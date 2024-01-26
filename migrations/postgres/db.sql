create table cities (
    id uuid primary key,
    name text check (char_length(name) > 3 AND char_length(name) <= 30),
    created_at timestamp default now()
);

create table customers (
    id uuid primary key,
    full_name text,
    phone text unique,
    email text unique,
    created_at timestamp default now()
);

create table drivers (
    id uuid primary key ,
    full_name text,
    phone text unique,
    from_city_id uuid references cities(id),
    to_city_id uuid references cities(id),
    created_at timestamp default now()
);

create table cars (
    id uuid primary key ,
    model varchar(30),
    brand varchar(30),
    number varchar(30) unique,
    status boolean default true,
    driver_id uuid references drivers(id),
    created_at timestamp default now()
);

create table trips (
    id uuid primary key,
    trip_number_id varchar(5) unique,
    from_city_id uuid references cities(id),
    to_city_id uuid references cities(id),
    driver_id uuid references drivers(id),
    price int default 0 check (price >= 0),
    created_at timestamp default now()
);

create table trip_customers (
    id uuid primary key,
    trip_id uuid references trips(id),
    customer_id uuid references customers(id),
    created_at timestamp default now()
);





CREATE SEQUENCE trip_num_id_sequence START 1;

CREATE OR REPLACE FUNCTION generate_trip_number_id() RETURNS TRIGGER AS $$
BEGIN
    NEW.trip_number_id := 'T-' || nextval('trip_number_id_sequence')::text;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trips_before_insert_trigger
BEFORE INSERT ON trips
FOR EACH ROW EXECUTE FUNCTION generate_trip_number_id();