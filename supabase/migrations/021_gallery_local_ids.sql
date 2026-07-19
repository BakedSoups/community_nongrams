create or replace function public.browse_gallery(p_kind text default 'art', p_sort text default 'new')
returns jsonb language plpgsql stable security definer set search_path = public as $$
declare result jsonb;
begin
  if p_kind not in ('art', 'pack') then raise exception 'Invalid gallery kind'; end if;
  if p_sort not in ('new', 'top', 'played') then raise exception 'Invalid gallery sort'; end if;

  if p_kind = 'art' then
    select coalesce(jsonb_agg(item order by
      case when p_sort = 'top' then (item->>'likes')::integer end desc,
      case when p_sort in ('top', 'played') then (item->>'plays')::integer end desc,
      (item->>'publishedAt')::timestamptz desc), '[]'::jsonb) into result
    from (
      select jsonb_build_object(
        'kind', 'art', 'id', level.id, 'localId', level.local_id, 'ownerId', level.owner_id,
        'creatorName', profile.display_name, 'avatarPuzzle', profile.avatar_puzzle, 'title', version.title,
        'description', version.description,
        'plays', (select count(*) from public.play_events where level_id = level.id),
        'likes', (select count(*) from public.likes where level_id = level.id),
        'liked', exists(select 1 from public.likes where level_id = level.id and user_id = auth.uid()),
        'owned', level.owner_id = auth.uid(),
        'promoted', exists(select 1 from public.profile_promotions where owner_id = level.owner_id and level_id = level.id),
        'previewPixels', level.preview_pixels,
        'puzzle', version.puzzle, 'publishedAt', version.published_at
      ) item
      from public.levels level
      join public.profiles profile on profile.id = level.owner_id
      join public.level_versions version on version.level_id = level.id and version.version = level.current_version
      where level.status = 'published' and level.visibility = 'public'
      limit 200
    ) gallery;
  else
    select coalesce(jsonb_agg(item order by
      case when p_sort = 'top' then (item->>'likes')::integer end desc,
      case when p_sort in ('top', 'played') then (item->>'plays')::integer end desc,
      (item->>'publishedAt')::timestamptz desc), '[]'::jsonb) into result
    from (
      select jsonb_build_object(
        'kind', 'pack', 'id', pack.id, 'ownerId', pack.owner_id,
        'creatorName', profile.display_name, 'avatarPuzzle', profile.avatar_puzzle, 'title', pack_version.title,
        'description', pack_version.description,
        'plays', coalesce((select count(*) from public.play_events play
          join public.level_versions played_version on played_version.level_id = play.level_id
          join public.pack_items pack_item on pack_item.level_version_id = played_version.id
          where pack_item.pack_version_id = pack_version.id), 0),
        'likes', (select count(*) from public.pack_likes where pack_id = pack.id),
        'liked', exists(select 1 from public.pack_likes where pack_id = pack.id and user_id = auth.uid()),
        'owned', pack.owner_id = auth.uid(),
        'promoted', exists(select 1 from public.profile_promotions where owner_id = pack.owner_id and pack_id = pack.id),
        'previewPixels', pack.preview_pixels,
        'levels', coalesce((select jsonb_agg(jsonb_build_object(
          'id', level_version.id, 'levelId', level.id, 'localId', level.local_id, 'version', level_version.version,
          'title', level_version.title, 'description', level_version.description,
          'tags', level_version.tags, 'puzzle', level_version.puzzle,
          'publishedAt', level_version.published_at
        ) order by pack_item.position)
        from public.pack_items pack_item
        join public.level_versions level_version on level_version.id = pack_item.level_version_id
        join public.levels level on level.id = level_version.level_id
        where pack_item.pack_version_id = pack_version.id), '[]'::jsonb),
        'publishedAt', pack_version.published_at
      ) item
      from public.packs pack
      join public.profiles profile on profile.id = pack.owner_id
      join public.pack_versions pack_version on pack_version.pack_id = pack.id and pack_version.version = pack.current_version
      where pack.status = 'published' and pack.visibility = 'public'
      limit 200
    ) gallery;
  end if;
  return result;
end;
$$;

grant execute on function public.browse_gallery(text, text) to anon, authenticated;
