create table scenario_template (
    id         varchar(128) primary key,
    project    varchar(128),
    desciption varchar(256),
    template   text
);

create table scenario (
    id          serial primary key,
    template_id varchar(128),
    params      jsonb,
    constraint fk_scenario_template
      foreign key(template_id)
	  references scenario_template(id)
	  on delete set null
);