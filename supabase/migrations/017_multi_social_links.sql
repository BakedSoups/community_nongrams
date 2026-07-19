alter table public.profiles
  drop constraint if exists profiles_social_check;

alter table public.profiles
  add constraint profiles_social_check check (
    char_length(social) <= 160
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
    raise exception 'Social must use supported profile entries';
  end if;
  update public.profiles
  set avatar_puzzle = p_avatar_puzzle,
      bio = left(coalesce(p_bio, ''), 120),
      display_name = left(trim(p_display_name), 40),
      social = left(trim(coalesce(p_social, '')), 160)
  where id = auth.uid();
end;
$$;
