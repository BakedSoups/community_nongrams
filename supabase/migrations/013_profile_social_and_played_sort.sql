alter table public.profiles
  add column if not exists social text not null default '' check (
    char_length(social) <= 80
    and social !~* 'https?://'
    and social !~* 'www\.'
    and social !~ '[/\\]'
  );

create or replace function public.save_creator_profile(
  p_avatar_puzzle jsonb,
  p_bio text default '',
  p_display_name text default 'Creator',
  p_social text default ''
) returns void language plpgsql security definer set search_path = public as $$
begin
  if auth.uid() is null then raise exception 'Sign in before saving a profile'; end if;
  if trim(coalesce(p_display_name, '')) = '' then raise exception 'Name is required'; end if;
  if p_social ~* 'https?://' or p_social ~* 'www\.' or p_social ~ '[/\\]' then
    raise exception 'Social must be a handle, not a link';
  end if;
  update public.profiles
  set avatar_puzzle = p_avatar_puzzle,
      bio = left(coalesce(p_bio, ''), 120),
      display_name = left(trim(p_display_name), 40),
      social = left(trim(coalesce(p_social, '')), 80)
  where id = auth.uid();
end;
$$;

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
        'kind', 'art', 'id', level.id, 'ownerId', level.owner_id,
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
          'id', level_version.id, 'levelId', level.id, 'version', level_version.version,
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

create or replace function public.browse_creators()
returns jsonb language sql stable security definer set search_path = public as $$
  select coalesce(jsonb_agg(creator order by creator->>'displayName'), '[]'::jsonb)
  from (
    select jsonb_build_object(
      'id', profile.id, 'displayName', profile.display_name, 'bio', profile.bio, 'social', profile.social, 'avatarPuzzle', profile.avatar_puzzle,
      'featured', coalesce((select jsonb_agg(promoted_item.featured) from (
        select jsonb_build_object(
          'kind', 'art', 'id', level.id, 'ownerId', level.owner_id, 'creatorName', profile.display_name,
          'title', version.title, 'likes', (select count(*) from public.likes where level_id = level.id),
          'puzzle', version.puzzle, 'promoted', true, 'publishedAt', version.published_at
        ) featured from public.profile_promotions promotion
        join public.levels level on level.id = promotion.level_id
        join public.level_versions version on version.level_id = level.id and version.version = level.current_version
        where promotion.owner_id = profile.id
        union all
        select jsonb_build_object(
          'kind', 'pack', 'id', pack.id, 'ownerId', pack.owner_id, 'creatorName', profile.display_name,
          'title', pack_version.title, 'likes', (select count(*) from public.pack_likes where pack_id = pack.id),
          'levels', coalesce((select jsonb_agg(jsonb_build_object('id', lv.id, 'levelId', l.id, 'version', lv.version, 'title', lv.title, 'puzzle', lv.puzzle, 'publishedAt', lv.published_at) order by pi.position)
            from public.pack_items pi join public.level_versions lv on lv.id = pi.level_version_id join public.levels l on l.id = lv.level_id where pi.pack_version_id = pack_version.id), '[]'::jsonb),
          'promoted', true, 'publishedAt', pack_version.published_at
        ) featured from public.profile_promotions promotion
        join public.packs pack on pack.id = promotion.pack_id
        join public.pack_versions pack_version on pack_version.pack_id = pack.id and pack_version.version = pack.current_version
        where promotion.owner_id = profile.id
      ) promoted_item), '[]'::jsonb),
      'levels', coalesce((select jsonb_agg(jsonb_build_object(
        'id', version.id, 'levelId', level.id, 'version', version.version, 'title', version.title,
        'description', version.description, 'tags', version.tags, 'puzzle', version.puzzle, 'publishedAt', version.published_at
      ) order by version.published_at desc)
      from public.levels level join public.level_versions version on version.level_id = level.id and version.version = level.current_version
      where level.owner_id = profile.id and level.status = 'published' and level.visibility = 'public'), '[]'::jsonb)
    ) creator from public.profiles profile order by profile.created_at desc limit 100
  ) creators;
$$;

grant execute on function public.save_creator_profile(jsonb, text, text, text) to authenticated;
grant execute on function public.browse_gallery(text, text) to anon, authenticated;
grant execute on function public.browse_creators() to anon, authenticated;
