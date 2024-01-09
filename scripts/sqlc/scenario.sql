-- name: CreateScenario :exec
insert into scenario (
    template_id,
    params
) VALUES ($1, $2);

-- name: DeleteScenario :exec
delete from scenario
where id = $1;

-- name: UpdateScenarioParams :exec
update scenario
set params = $1
where id = $2;

-- name: GetScenario :one
select id, template_id, params from scenario
where id = $1;

-- name: GetScenariosByTemplate :many
select id, template_id, params from asset_series
where template_id = $1
order by id;