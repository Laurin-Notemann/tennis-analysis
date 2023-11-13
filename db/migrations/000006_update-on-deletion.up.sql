BEGIN;
  ALTER TABLE "stats" DROP CONSTRAINT "FK_Stats.team_id";
  ALTER TABLE "stats" ADD CONSTRAINT "FK_Stats.team_id" FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE;

  DROP TRIGGER IF EXISTS delete_stats_after_team_delete ON teams;
  DROP TRIGGER IF EXISTS delete_stats_after_game_delete ON games;
  DROP FUNCTION IF EXISTS delete_stats_row();

COMMIT;
