drop table if exists report_number_counters;
drop table if exists reports;
drop table if exists moderation_logs;
alter table books
    drop column if exists is_shadow_banned,
    drop column if exists is_perm_removed;
