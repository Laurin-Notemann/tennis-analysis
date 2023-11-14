BEGIN;
  DROP TRIGGER IF EXISTS trigger_delete_user_cascade ON users;
  DROP FUNCTION IF EXISTS delete_user_cascade();
COMMIT;
