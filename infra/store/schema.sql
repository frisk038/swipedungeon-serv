DROP TABLE userlocation;
DROP TABLE users;

CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text UNIQUE NOT NULL,
    player_id text NOT NULL,
    power_type text,
    chara_id int
);

CREATE TABLE userlocation (
    user_id uuid REFERENCES users(user_id),
    coord geography(point) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON userlocation USING gist(coord);

CREATE TABLE userscore (
    user_id uuid REFERENCES users(user_id),
    user_level int,
    score int,
    PRIMARY KEY(user_id)
);

ALTER TABLE userlocation ADD city text DEFAULT ''; 

CREATE TABLE seen (
    user_id uuid REFERENCES users(user_id),
    seen_user uuid REFERENCES users(user_id),
    seen_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, seen_user)
);