-- name: SearchDefinedTags :many
select *
from defined_tags
where lowercased_name like $1
limit $2;

-- name: SearchDefinedTagsWithType :many
select *
from defined_tags
where lowercased_name like $1 and tag_type = $3
limit $2;

-- name: InsertDefinedTag :exec
insert into defined_tags
(id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: InsertDefinedTagEnMasse :copyfrom
insert into defined_tags
(id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of)
values ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: DefinedTagsAreInitialized :one
select exists(select 1 from defined_tags) as initialized;

-- name: GetTag :one
select defined_tags.*, syn.name as synonym_name
from defined_tags
left join defined_tags as syn on defined_tags.synonym_of = syn.id 
where defined_tags.id = $1;

-- name: GetTagParent :one
select t0.* 
from defined_tags t0
where t0.id = (select coalesce(t1.synonym_of, t1.id) from defined_tags t1 where t1.id = $1);

-- name: GetTagsByName :many
select * from defined_tags where name = ANY(sqlc.arg('names')::text[]);

-- name: GetTagsByIds :many
select * from defined_tags where id = ANY(sqlc.arg('ids')::int8[]);

-- name: ImportTags :exec
insert into defined_tags
(id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of, is_default)
select
    id, name, description, is_spoiler, is_adult, created_at, tag_type, synonym_of, true
from jsonb_to_recordset($1::jsonb)
as json_set (
    id int8, 
    name text, 
    description text, 
    is_spoiler boolean, 
    is_adult boolean, 
    created_at timestamptz, 
    tag_type tag_type, 
    synonym_of int8)
where not exists (select 1 from defined_tags where name = json_set.name);

-- name: RemoveUnusedDefaultTags :exec
delete from defined_tags d
where 
    d.name <> ANY(sqlc.arg(names)::text[]) and
    d.is_default = true and 
    not exists (select 1 from defined_tags where synonym_of = d.id) and
    not exists (select 1 from books where tag_ids @> array[d.id]);

