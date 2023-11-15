BEGIN;
  ALTER TABLE "stats" DROP CONSTRAINT "FK_Stats.team_id";
  ALTER TABLE "stats" ADD CONSTRAINT "FK_Stats.team_id" FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL;
  CREATE OR REPLACE FUNCTION delete_stats_row()
    RETURNS TRIGGER AS $$
    BEGIN
      IF OLD.team_id IS NULL AND OLD.game_id IS NULL THEN
          DELETE FROM stats WHERE id = OLD.id;
      END IF;
      RETURN OLD;
    END;
  $$ LANGUAGE plpgsql;

  CREATE TRIGGER delete_stats_after_team_delete
  AFTER DELETE ON teams
  FOR EACH ROW
  EXECUTE FUNCTION delete_stats_row();

  CREATE TRIGGER delete_stats_after_game_delete
  AFTER DELETE ON games
  FOR EACH ROW
  EXECUTE FUNCTION delete_stats_row();
COMMIT;
