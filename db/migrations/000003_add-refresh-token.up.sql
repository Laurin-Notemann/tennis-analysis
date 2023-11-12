BEGIN;
  CREATE TABLE "refresh_tokens" (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid UNIQUE NOT NULL,
    token text UNIQUE NOT NULL,
    expiry_date timestamptz NOT NULL DEFAULT (NOW()+ INTERVAL '30 days'),

    created_at timestamptz NOT NULL DEFAULT Now(),
    updated_at timestamptz NOT NULL DEFAULT Now(),
    PRIMARY KEY (id),
    CONSTRAINT "FK_Refresh_token.user_id" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
  );

  ALTER TABLE "users" ADD COLUMN refresh_token_id uuid;
  ALTER TABLE "users" ADD CONSTRAINT "FK_Users.refresh_token_id" FOREIGN KEY (refresh_token_id) REFERENCES refresh_tokens(id) ON DELETE SET NULL;
COMMIT;
