BEGIN;
  DROP TABLE IF EXISTS "stats";

  DROP TRIGGER IF EXISTS delete_stats_after_team_delete ON teams;
  DROP TRIGGER IF EXISTS delete_stats_after_game_delete ON games;
  DROP FUNCTION IF EXISTS delete_stats_row();

COMMIT;
