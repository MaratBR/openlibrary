drop table if exists comment_logs;
drop index if exists ix_user_logs_user_time;

delete from user_logs
where action_type not in ('sec_password_reset', 'sec_2fa_umbrella', 'ban', 'unban', 'mute', 'unmute');

alter table user_logs
    drop column if exists "time",
    drop column if exists reason,
    alter column action_type type user_action_type using action_type::user_action_type;
