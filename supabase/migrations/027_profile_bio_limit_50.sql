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
      bio = left(coalesce(p_bio, ''), 50),
      display_name = left(trim(p_display_name), 40),
      social = left(trim(coalesce(p_social, '')), 160),
      favorite_palette = left(trim(coalesce(p_palette, '')), 40),
      favorite_color = left(trim(coalesce(p_favorite_color, '')), 9)
  where id = auth.uid();
end;
$$;

grant execute on function public.save_creator_profile(jsonb, text, text, text, text, text) to authenticated;
