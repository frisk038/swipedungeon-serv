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

INSERT INTO userscore (user_id, score, user_level)
VALUES ('d84bbea4-4cb3-4539-b5f4-d04c0067c61e', 9, 1)
ON CONFLICT (user_id) 
DO UPDATE
SET score = GREATEST(EXCLUDED.score, userscore.score),
user_level = CASE WHEN EXCLUDED.score > userscore.score THEN EXCLUDED.user_level ELSE userscore.user_level END;