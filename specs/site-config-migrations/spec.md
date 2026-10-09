# Versioned settings rows

Each settings type (fonts, captcha, password requirements, content restrictions)
is stored in its own keyed database row with independent Version and Data.
SettingsDefinition owns its defaults and migration history. Loading runs only
pending migrations for each known definition, in order. The initial font
migration initializes an empty FontConfiguration.

Missing settings use their definition's defaults and migrations. When there are
no stored rows, static defaults remain active. Current settings retain their
values. Unknown rows remain untouched. Upgrades are persisted atomically before
activation; invalid data, unsupported versions, or migration/write failures leave
the active snapshot intact and return an error.

A reversible database migration splits the legacy main JSON into settings rows,
preserving existing values and assigning the legacy version to font settings.
Other settings begin at version zero. There is no global settings version.
