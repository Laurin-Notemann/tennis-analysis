BEGIN;
  ALTER TABLE "games" ADD COLUMN "winner" uuid;
  ALTER TABLE "games" ADD CONSTRAINT "FK_Games.winner" FOREIGN KEY ("winner") REFERENCES teams(id);
  
  ALTER TABLE "games" ADD COLUMN "match_id" uuid;
  ALTER TABLE "games" ADD CONSTRAINT "FK_Games.match_id" FOREIGN KEY ("match_id") REFERENCES matches(id);

  ALTER TABLE "sets" ADD COLUMN "winner" uuid;
  ALTER TABLE "sets" ADD CONSTRAINT "FK_Sets.winner" FOREIGN KEY ("winner") REFERENCES teams(id);

  ALTER TABLE "matches" ADD COLUMN "winner" uuid NOT NULL;
  ALTER TABLE "matches" ADD CONSTRAINT "FK_Matches.winner" FOREIGN KEY ("winner") REFERENCES teams(id);
COMMIT;

