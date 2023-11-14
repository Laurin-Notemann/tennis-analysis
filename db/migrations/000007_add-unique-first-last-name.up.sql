BEGIN;
  ALTER TABLE "players" ADD CONSTRAINT "unique_firstlast_name" UNIQUE (first_name, last_name);
COMMIT;
