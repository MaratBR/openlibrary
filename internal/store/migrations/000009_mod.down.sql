drop table if exists user_logs;
drop type if exists user_action_type;
drop table if exists book_logs;
drop type if exists book_action_type;
alter table books drop column if exists is_perm_removed;