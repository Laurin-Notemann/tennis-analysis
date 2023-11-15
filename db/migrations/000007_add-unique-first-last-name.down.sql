BEGIN;
  ALTER TABLE "players" DROP CONSTRAINT "unique_firstlast_name";
COMMIT;
