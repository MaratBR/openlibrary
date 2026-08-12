-- Run this only after the database has successfully applied all 16 migrations
-- from the pre-consolidation migration set. The schema itself is unchanged;
-- only golang-migrate's recorded version needs to be remapped.
begin;

do $migration_state$
declare
    current_version bigint;
    current_dirty boolean;
begin
    select version, dirty
      into current_version, current_dirty
      from schema_migrations;

    if not found then
        raise exception 'schema_migrations has no migration state';
    end if;

    if current_version <> 16 or current_dirty then
        raise exception 'expected clean migration version 16, found version %, dirty %',
            current_version, current_dirty;
    end if;

    update schema_migrations
       set version = 3,
           dirty = false;
end
$migration_state$;

commit;
