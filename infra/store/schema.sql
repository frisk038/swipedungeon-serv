CREATE TABLE users (
    user_id int GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    player_id text NOT NULL,
    power_type text NOT NULL
)