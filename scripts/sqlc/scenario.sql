-- name: AddSeriesItem :exec
insert into asset_series (
    symbol,
    date,
    price
) VALUES ($1, $2, $3)
on conflict on constraint unique_elem
do update
set price = EXCLUDED.price;


-- name: DeleteSeries :exec
delete from asset_series
where symbol = $1;

-- name: DeleteSeriesBefore :exec
delete from asset_series
where date < $1 and symbol = $2;

-- name: GetSeries :many
select date, price from asset_series
where symbol = $1
order by date;

-- name: UpdateSeriesItem :exec
update asset_series
set state_id = $3
where symbol = $1 and date = $2;

-- name: GetLastSeriesItem :one
select date, price, st.state_id, st.rank from asset_series t1
left join asset_state st
on t1.state_id = st.state_id
where t1.symbol = $1 and not exists (
  select *
  from asset_series t2
  where t2.symbol = $1 and t2.date > t1.date
);

-- name: GetSeriesWithStates :many
select date, price, st.state_id, st.rank from asset_series se
join asset_state st
on se.state_id = st.state_id
where symbol = $1
order by date;

-- name: GetSeriesAfter :many
select date, price from asset_series
where symbol = $1 and date >= $2
order by date;

-- name: GetSeriesWithStatesAfter :many
select date, price, st.state_id, st.rank from asset_series se
join asset_state st
on se.state_id = st.state_id
where symbol = $1 and date >= $2
order by date;