CREATE TABLE users(
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  username text NOT NULL,
  email text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id)
);


CREATE TABLE users(
  id uuid NOT NULL DEFAULT gen_random_uuid(),

  created_at timestamptz NOT NULL DEFAULT Now(),
  updated_at timestamptz NOT NULL DEFAULT Now(),
  PRIMARY KEY (id)
);
