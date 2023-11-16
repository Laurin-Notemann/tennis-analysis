CREATE TABLE "users" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  username text UNIQUE NOT NULL,
  email text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id),
  CONSTRAINT "users_username_unique" UNIQUE("username"),
	CONSTRAINT "users_email_unique" UNIQUE("email")
);

CREATE TABLE "players" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  first_name text NOT NULL,
  last_name text NOT NULL,

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id)
);


CREATE TABLE "teams" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name text NOT NULL,

  user_id uuid NOT NULL,
  player_one uuid NOT NULL,
  player_two uuid,

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id),
  CONSTRAINT "FK_Teams.user_id" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Teams.player_one" FOREIGN KEY (player_one) REFERENCES players(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Teams.player_two" FOREIGN KEY (player_two) REFERENCES players(id) ON DELETE SET NULL,
  CONSTRAINT "CK_Teams_DistinctPlayers" CHECK (player_one <> player_two OR player_two IS NULL),
  CONSTRAINT "unique_players" UNIQUE (player_one, player_two)
);


CREATE TABLE "matches" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  number_of_sets INT,
  
  user_id uuid NOT NULL,
  team_one uuid NOT NULL,
  team_two uuid NOT NULL,
  winner uuid NOT NULL,

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id),
  CONSTRAINT "FK_Matches.user_id" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Matches.team_one" FOREIGN KEY (team_one) REFERENCES teams(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Matches.team_two" FOREIGN KEY (team_two) REFERENCES teams(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Matches.winner" FOREIGN KEY (winner) REFERENCES teams(id) ON DELETE CASCADE,
  CONSTRAINT "CK_Matches_DistinctPlayers" CHECK (team_one <> team_two)
);

CREATE TABLE "sets" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),

  match_id uuid NOT NULL,
  winner uuid,

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id),
  CONSTRAINT "FK_Sets.match_id" FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Sets.winner" FOREIGN KEY (winner) REFERENCES teams(id) ON DELETE SET NULL
);

CREATE TABLE "games" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),

  winner uuid,
  server_id uuid NOT NULL,
  set_id uuid,
  match_id uuid,

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id),
  CONSTRAINT "FK_Games.winner" FOREIGN KEY (winner) REFERENCES teams(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Games.server_id" FOREIGN KEY (server_id) REFERENCES teams(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Games.set_id" FOREIGN KEY (set_id) REFERENCES sets(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Games.match_id" FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE
);

CREATE TABLE "points" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  value INT,

  team_id uuid NOT NULL,
  score_id uuid NOT NULL,

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id),
  CONSTRAINT "FK_Points.team_id" FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

CREATE TABLE "scores" (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  score_order INT,

  team_one_point uuid NOT NULL,
  team_two_point uuid NOT NULL,
  game_id uuid NOT NULL,

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id),
  CONSTRAINT "FK_Scores.team_one_point" FOREIGN KEY (team_one_point) REFERENCES points(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Scores.team_two_point" FOREIGN KEY (team_two_point) REFERENCES points(id) ON DELETE CASCADE,
  CONSTRAINT "FK_Scores.game_id" FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
  CONSTRAINT "CK_Scores_DistinctTeams" CHECK (team_one_point <> team_two_point)
);

ALTER TABLE "points" ADD CONSTRAINT "FK_Points.score_id" FOREIGN KEY (score_id) REFERENCES scores(id) ON DELETE CASCADE;
