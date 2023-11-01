BEGIN;
  
ALTER TABLE "points" DROP CONSTRAINT "FK_Points.game_id";
ALTER TABLE  "points" DROP COLUMN game_id;
ALTER TABLE "points" ADD COLUMN score_id uuid;
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
ALTER TABLE "points" DROP COLUMN points_order;

COMMIT;
