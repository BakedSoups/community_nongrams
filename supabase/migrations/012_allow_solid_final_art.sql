create or replace function public.publish_level(
  p_local_id text,
  p_title text,
  p_description text,
  p_tags text[],
  p_puzzle jsonb,
  p_submit_official boolean default false,
  p_rights_confirmed boolean default false
) returns jsonb language plpgsql security definer set search_path = public as $$
declare
  uid uuid := auth.uid();
  target_level public.levels;
  target_version public.level_versions;
  next_version integer;
  width integer := (p_puzzle->>'width')::integer;
  height integer := (p_puzzle->>'height')::integer;
  row_text text;
  pixel_row jsonb;
  pixel_value text;
  filled_cells integer := 0;
begin
  if uid is null then raise exception 'Sign in before publishing'; end if;
  if trim(p_title) = '' or char_length(p_title) > 80 then raise exception 'Title must be 1 to 80 characters'; end if;
  if width not in (8, 10, 15, 20) or height not in (8, 10, 15, 20) then raise exception 'Unsupported puzzle dimensions'; end if;
  if jsonb_typeof(p_puzzle->'solution') <> 'array' or jsonb_array_length(p_puzzle->'solution') <> height then raise exception 'Invalid solution rows'; end if;
  for row_text in select jsonb_array_elements_text(p_puzzle->'solution') loop
    if char_length(row_text) <> width or row_text !~ '^[01]+$' then raise exception 'Invalid solution row'; end if;
    filled_cells := filled_cells + char_length(replace(row_text, '0', ''));
  end loop;
  if filled_cells = 0 then raise exception 'Puzzle must contain at least one filled cell'; end if;
  if jsonb_typeof(p_puzzle->'skeletonPixels') <> 'array' or jsonb_array_length(p_puzzle->'skeletonPixels') <> height then raise exception 'Invalid Before layer'; end if;
  for pixel_row in select jsonb_array_elements(p_puzzle->'skeletonPixels') loop
    if jsonb_typeof(pixel_row) <> 'array' or jsonb_array_length(pixel_row) <> width then raise exception 'Invalid Before row'; end if;
    for pixel_value in select jsonb_array_elements_text(pixel_row) loop
      if lower(pixel_value) not in ('', 'transparent', '#000000ff') then raise exception 'Before art must be black or transparent'; end if;
    end loop;
  end loop;
  if jsonb_typeof(p_puzzle->'revealPixels') <> 'array' or jsonb_array_length(p_puzzle->'revealPixels') <> height then raise exception 'Invalid After layer'; end if;
  for pixel_row in select jsonb_array_elements(p_puzzle->'revealPixels') loop
    if jsonb_typeof(pixel_row) <> 'array' or jsonb_array_length(pixel_row) <> width then raise exception 'Invalid After row'; end if;
  end loop;
  if p_submit_official and not p_rights_confirmed then raise exception 'Rights confirmation is required'; end if;

  insert into public.levels (owner_id, local_id, title, description, tags)
  values (uid, p_local_id, trim(p_title), coalesce(p_description, ''), coalesce(p_tags, '{}'))
  on conflict (owner_id, local_id) do update set
    title = excluded.title, description = excluded.description, tags = excluded.tags,
    status = 'published', visibility = 'public',
    current_version = public.levels.current_version + 1, updated_at = now()
  returning * into target_level;

  next_version := target_level.current_version;
  insert into public.level_versions (level_id, version, title, description, tags, puzzle)
  values (target_level.id, next_version, target_level.title, target_level.description, target_level.tags, p_puzzle)
  returning * into target_version;

  if p_submit_official then
    insert into public.official_submissions (owner_id, level_id, level_version_id, rights_confirmed)
    values (uid, target_level.id, target_version.id, true);
  end if;
  return jsonb_build_object('levelId', target_level.id, 'levelVersionId', target_version.id, 'version', next_version);
end;
$$;

grant execute on function public.publish_level(text, text, text, text[], jsonb, boolean, boolean) to authenticated;
