create table scenario_template (
    id         serial primary key,
    project    varchar(128),
    name       varchar(128),
    desciption varchar(256),
    template   text
);

create table scenario (
    scenario_templ_id  int,
    params             jsonb
    -- constraint fk_scenario_template
    --   foreign key(scenario_templ_id)
	--   references scenario_template(id)
	--   on delete cascade
);