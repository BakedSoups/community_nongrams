create or replace function public.update_published_content(
  p_kind text,
  p_content_id uuid,
  p_title text,
  p_description text default ''
) returns void language plpgsql security definer set search_path = public as $$
declare
  next_version integer;
  previous_level_version public.level_versions%rowtype;
  previous_pack_version public.pack_versions%rowtype;
  new_pack_version_id uuid;
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

    next_version := previous_pack_version.version + 1;
    insert into public.pack_versions(pack_id, version, title, description)
    values (p_content_id, next_version, p_title, p_description)
    returning id into new_pack_version_id;

    insert into public.pack_items(pack_version_id, level_version_id, position)
    select new_pack_version_id, level_version_id, position
    from public.pack_items
    where pack_version_id = previous_pack_version.id
    order by position;

    update public.packs
    set title = p_title, description = p_description, current_version = next_version, updated_at = now()
    where id = p_content_id and owner_id = auth.uid();
  else
    raise exception 'Invalid published content kind';
  end if;
end;
$$;

grant execute on function public.update_published_content(text, uuid, text, text) to authenticated;
