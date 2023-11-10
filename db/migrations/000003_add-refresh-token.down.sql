BEGIN;
  ALTER TABLE "users" DROP CONSTRAINT "FK_Users.refresh_token_id";
  ALTER TABLE "users" DROP COLUMN refresh_token_id;
  DROP TABLE "refresh_tokens";

COMMIT;
