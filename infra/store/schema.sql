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

UPDATE userscore
	SET (score, user_level) = 
		(SELECT 
			CASE WHEN $1 > score THEN $1 ELSE score END,
			CASE WHEN $1 > score THEN $2 ELSE user_level END
		 	FROM userscore
			WHERE user_id = '5d28f627-b30c-4171-939f-cc577ced454')
		WHERE user_id = '5d28f627-b30c-4171-939f-cc577ced454';