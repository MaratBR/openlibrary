insert into site_config (key, value)
select setting.key, jsonb_build_object(
    'Version', case when setting.key = 'FontConfiguration'
                    then coalesce(legacy.value->'Version', '0'::jsonb)
                    else '0'::jsonb end,
    'Data', setting.value)
from site_config legacy
cross join lateral jsonb_each(legacy.value) setting
where legacy.key = 'main'
  and setting.key in ('FontConfiguration', 'CaptchaSettings',
                      'PasswordRequirements', 'ContentRestrictions')
on conflict (key) do nothing;

-- An absent legacy font section needs the initial font migration.
insert into site_config (key, value)
select 'FontConfiguration', '{"Version":0,"Data":{}}'::jsonb
from site_config where key = 'main'
on conflict (key) do nothing;

delete from site_config where key = 'main';
