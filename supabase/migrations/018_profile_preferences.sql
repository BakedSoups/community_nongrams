alter table public.profiles
  add column if not exists favorite_palette text not null default '' check (char_length(favorite_palette) <= 40),
  add column if not exists favorite_color text not null default '' check (favorite_color = '' or favorite_color ~ '^#[0-9A-Fa-f]{6}$');

create or replace function public.save_creator_profile(
  p_avatar_puzzle jsonb,
  p_bio text default '',
  p_display_name text default 'Creator',
  p_social text default '',
  p_palette text default '',
  p_favorite_color text default ''
) returns void language plpgsql security definer set search_path = public as $$
begin
  if auth.uid() is null then raise exception 'Sign in before saving a profile'; end if;
  if trim(coalesce(p_display_name, '')) = '' then raise exception 'Name is required'; end if;
  if p_social ~* 'https?://' or p_social ~* 'www\.' or p_social ~ '[/\\]' then
    raise exception 'Social must use supported profile entries';
  end if;
  if p_favorite_color <> '' and p_favorite_color !~ '^#[0-9A-Fa-f]{6}$' then
    raise exception 'Favorite color must be a hex color';
  end if;
  update public.profiles
  set avatar_puzzle = p_avatar_puzzle,
      bio = left(coalesce(p_bio, ''), 120),
      display_name = left(trim(p_display_name), 40),
      social = left(trim(coalesce(p_social, '')), 160),
      favorite_palette = left(trim(coalesce(p_palette, '')), 40),
      favorite_color = left(trim(coalesce(p_favorite_color, '')), 9)
  where id = auth.uid();
end;
$$;

create or replace function public.browse_creators()
returns jsonb language sql stable security definer set search_path = public as $$
  select coalesce(jsonb_agg(creator order by creator->>'displayName'), '[]'::jsonb)
  from (
    select jsonb_build_object(
      'id', profile.id, 'displayName', profile.display_name, 'bio', profile.bio, 'social', profile.social,
      'palette', profile.favorite_palette, 'favoriteColor', profile.favorite_color, 'avatarPuzzle', profile.avatar_puzzle,
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

grant execute on function public.save_creator_profile(jsonb, text, text, text, text, text) to authenticated;
grant execute on function public.browse_creators() to anon, authenticated;
