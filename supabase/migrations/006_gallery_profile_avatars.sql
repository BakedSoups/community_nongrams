do $$
declare
  definition text;
begin
  select pg_get_functiondef('public.browse_gallery(text,text)'::regprocedure)
  into definition;

  if position('avatarPuzzle' in definition) = 0 then
    definition := replace(
      definition,
      '''creatorName'', profile.display_name, ''title''',
      '''creatorName'', profile.display_name, ''avatarPuzzle'', profile.avatar_puzzle, ''title'''
    );
    execute definition;
  end if;
end;
$$;
