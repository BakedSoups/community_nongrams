create or replace function public.unpublish_community_local_art(p_local_id text)
returns void language plpgsql security definer set search_path = public as $$
declare
  target_id uuid;
begin
  if auth.uid() is null then raise exception 'Sign in to manage published work'; end if;

  update public.levels set status = 'hidden', updated_at = now()
  where local_id = p_local_id and owner_id = auth.uid() and status = 'published'
  returning id into target_id;

  if target_id is null then raise exception 'Published item not found'; end if;

  delete from public.profile_promotions
  where owner_id = auth.uid() and level_id = target_id;
end;
$$;

grant execute on function public.unpublish_community_local_art(text) to authenticated;
