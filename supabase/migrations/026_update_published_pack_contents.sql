drop function if exists public.update_published_content(text, uuid, text, text);

create or replace function public.update_published_content(
  p_kind text,
  p_content_id uuid,
  p_title text,
  p_description text default '',
  p_levels jsonb default null
) returns void language plpgsql security definer set search_path = public as $$
declare
  next_version integer;
  previous_level_version public.level_versions%rowtype;
  previous_pack_version public.pack_versions%rowtype;
  new_pack_version_id uuid;
  level_payload jsonb;
  target_level public.levels;
  target_version public.level_versions;
  level_local_id text;
  level_title text;
  level_description text;
  level_tags text[];
  level_puzzle jsonb;
  width integer;
  height integer;
  row_text text;
  pixel_row jsonb;
  pixel_value text;
  filled_cells integer;
  position_index integer := 0;
begin
  if auth.uid() is null then raise exception 'Sign in to manage published work'; end if;

  p_title := left(trim(coalesce(p_title, '')), 80);
  p_description := left(coalesce(p_description, ''), 500);
  if p_title = '' then raise exception 'Name is required'; end if;

  if p_kind = 'art' then
    select version.* into previous_level_version
    from public.levels level
    join public.level_versions version on version.level_id = level.id and version.version = level.current_version
    where level.id = p_content_id and level.owner_id = auth.uid() and level.status = 'published';

    if previous_level_version.id is null then raise exception 'Published art not found'; end if;

    next_version := previous_level_version.version + 1;
    insert into public.level_versions(level_id, version, title, description, tags, puzzle)
    values (p_content_id, next_version, p_title, p_description, previous_level_version.tags, previous_level_version.puzzle);

    update public.levels
    set title = p_title, description = p_description, current_version = next_version, updated_at = now()
    where id = p_content_id and owner_id = auth.uid();
  elsif p_kind = 'pack' then
    select version.* into previous_pack_version
    from public.packs pack
    join public.pack_versions version on version.pack_id = pack.id and version.version = pack.current_version
    where pack.id = p_content_id and pack.owner_id = auth.uid() and pack.status = 'published';

    if previous_pack_version.id is null then raise exception 'Published pack not found'; end if;

    if p_levels is not null then
      if jsonb_typeof(p_levels) <> 'array' or jsonb_array_length(p_levels) not between 1 and 20 then
        raise exception 'Packs must contain 1 to 20 levels';
      end if;
      if (select count(distinct value->>'id') from jsonb_array_elements(p_levels) value) <> jsonb_array_length(p_levels) then
        raise exception 'A level can appear only once';
      end if;
    end if;

    next_version := previous_pack_version.version + 1;
    insert into public.pack_versions(pack_id, version, title, description)
    values (p_content_id, next_version, p_title, p_description)
    returning id into new_pack_version_id;

    if p_levels is null then
      insert into public.pack_items(pack_version_id, level_version_id, position)
      select new_pack_version_id, level_version_id, position
      from public.pack_items
      where pack_version_id = previous_pack_version.id
      order by position;
    else
      for level_payload in select value from jsonb_array_elements(p_levels) value loop
        level_local_id := level_payload->>'id';
        level_title := trim(coalesce(level_payload->>'title', ''));
        level_description := coalesce(level_payload->>'description', '');
        level_tags := coalesce(array(select jsonb_array_elements_text(coalesce(level_payload->'tags', '[]'::jsonb))), '{}');
        level_puzzle := level_payload->'puzzle';
        width := (level_puzzle->>'width')::integer;
        height := (level_puzzle->>'height')::integer;
        filled_cells := 0;

        if coalesce(level_local_id, '') = '' then raise exception 'Pack level is missing an id'; end if;
        if level_title = '' or char_length(level_title) > 80 then raise exception 'Pack level title must be 1 to 80 characters'; end if;
        if width not in (8, 10, 15, 20, 32) or height not in (8, 10, 15, 20, 32) then raise exception 'Unsupported puzzle dimensions'; end if;
        if jsonb_typeof(level_puzzle->'solution') <> 'array' or jsonb_array_length(level_puzzle->'solution') <> height then raise exception 'Invalid solution rows'; end if;
        for row_text in select jsonb_array_elements_text(level_puzzle->'solution') loop
          if char_length(row_text) <> width or row_text !~ '^[01]+$' then raise exception 'Invalid solution row'; end if;
          filled_cells := filled_cells + char_length(replace(row_text, '0', ''));
        end loop;
        if filled_cells = 0 then raise exception 'Puzzle must contain at least one filled cell'; end if;
        if jsonb_typeof(level_puzzle->'skeletonPixels') <> 'array' or jsonb_array_length(level_puzzle->'skeletonPixels') <> height then raise exception 'Invalid Before layer'; end if;
        for pixel_row in select jsonb_array_elements(level_puzzle->'skeletonPixels') loop
          if jsonb_typeof(pixel_row) <> 'array' or jsonb_array_length(pixel_row) <> width then raise exception 'Invalid Before row'; end if;
          for pixel_value in select jsonb_array_elements_text(pixel_row) loop
            if lower(pixel_value) not in ('', 'transparent', '#000000ff') then raise exception 'Before art must be black or transparent'; end if;
          end loop;
        end loop;
        if jsonb_typeof(level_puzzle->'revealPixels') <> 'array' or jsonb_array_length(level_puzzle->'revealPixels') <> height then raise exception 'Invalid After layer'; end if;

        insert into public.levels (owner_id, local_id, title, description, tags, visibility, status)
        values (auth.uid(), level_local_id, level_title, level_description, level_tags, 'pack_only'::public.content_visibility, 'published'::public.content_status)
        on conflict (owner_id, local_id) do update set
          title = excluded.title,
          description = excluded.description,
          tags = excluded.tags,
          status = 'published'::public.content_status,
          visibility = case
            when public.levels.visibility = 'public'::public.content_visibility then 'public'::public.content_visibility
            else 'pack_only'::public.content_visibility
          end,
          current_version = public.levels.current_version + 1,
          updated_at = now()
        returning * into target_level;

        insert into public.level_versions(level_id, version, title, description, tags, puzzle)
        values (target_level.id, target_level.current_version, target_level.title, target_level.description, target_level.tags, level_puzzle)
        returning * into target_version;

        insert into public.pack_items(pack_version_id, level_version_id, position)
        values (new_pack_version_id, target_version.id, position_index);
        position_index := position_index + 1;
      end loop;
    end if;

    update public.packs
    set title = p_title, description = p_description, current_version = next_version, updated_at = now()
    where id = p_content_id and owner_id = auth.uid();
  else
    raise exception 'Invalid published content kind';
  end if;
end;
$$;

grant execute on function public.update_published_content(text, uuid, text, text, jsonb) to authenticated;
