insert into site_config (key, value)
select 'main', jsonb_object_agg(key, value->'Data') ||
       jsonb_build_object('Version', coalesce(max((value->>'Version')::int)
           filter (where key = 'FontConfiguration'), 0))
from site_config
where key in ('FontConfiguration', 'CaptchaSettings',
              'PasswordRequirements', 'ContentRestrictions')
having count(*) > 0
on conflict (key) do update set value = excluded.value;

delete from site_config
where key in ('FontConfiguration', 'CaptchaSettings',
              'PasswordRequirements', 'ContentRestrictions');
