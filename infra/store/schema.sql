DROP TABLE userlocation;
DROP TABLE users;

CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    player_id text UNIQUE NOT NULL,
    power_type text,
    chara_id int
);

CREATE TABLE userlocation (
    user_id uuid REFERENCES users(user_id),
    coord geography(point) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON userlocation USING gist(coord);
