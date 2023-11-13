BEGIN;
  CREATE TABLE "stats" (
    id uuid NOT NULL DEFAULT gen_random_uuid(),

    aces INT DEFAULT 0,
    double_faults INT DEFAULT 0,
    net_points INT DEFAULT 0,
    deuce INT DEFAULT 0,
    points_won_deuce INT DEFAULT 0,

    game_id uuid,
    team_id uuid,
    
    CONSTRAINT "FK_Stats.game_id" FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE SET NULL,
    CONSTRAINT "FK_Stats.team_id" FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL
  );
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
