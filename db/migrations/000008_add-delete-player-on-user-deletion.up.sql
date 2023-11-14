BEGIN;
  -- Create the trigger function
  CREATE OR REPLACE FUNCTION delete_user_cascade()
  RETURNS TRIGGER AS $$
  BEGIN
      -- Delete players associated with the user's teams
      DELETE FROM players
      WHERE id IN (
          SELECT player_one FROM teams WHERE user_id = OLD.id
      ) OR id IN (
          SELECT player_two FROM teams WHERE user_id = OLD.id
      );

      -- Delete teams associated with the user
      DELETE FROM teams
      WHERE user_id = OLD.id;

      -- Return the old value to indicate successful deletion
      RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;

  -- Create the trigger on the users table
  CREATE TRIGGER trigger_delete_user_cascade
  BEFORE DELETE ON users
  FOR EACH ROW EXECUTE FUNCTION delete_user_cascade();
COMMIT; 
