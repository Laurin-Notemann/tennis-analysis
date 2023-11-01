BEGIN;

ALTER TABLE "points" DROP CONSTRAINT "FK_Points.score_id";
ALTER TABLE "points" DROP COLUMN score_id;
ALTER TABLE  "points" ADD COLUMN game_id uuid;
ALTER TABLE "points" ADD CONSTRAINT "FK_Points.game_id" FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
ALTER TABLE "points" ADD COLUMN points_order INT;
DROP TABLE "scores" CASCADE;

COMMIT;
