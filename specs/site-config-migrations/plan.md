# Plan

Use SettingsDefinition and generic Definition[T] to describe each settings type.
Keep the abstraction, concrete definitions/migrations, and Store implementation
separate. Keep Get's typed aggregate snapshot for current callers, but persist
one {Version, Data} envelope per definition key. Register definitions with their
snapshot field bindings in settingsDefinitions.

Append migrations to each definition: migration i upgrades version i to i+1.
Never reorder published migrations. New types require a definition and snapshot
field binding; persistence queries do not need changes.

Split legacy main via migration 000004, with a down migration to reconstruct it.
Use SQLC queries to aggregate reads and atomically upsert known rows. Serialize
Load and Save. Remember successful saves to avoid unchanged writes.

Verify independent versions, migration order/errors, preserved current values,
missing settings, malformed data, unsupported versions, and failed writes.
