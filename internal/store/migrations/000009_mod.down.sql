drop table if exists reports;
drop table if exists moderation_logs;
alter table books drop column if exists is_perm_removed;
