# Tennis-Analysis (working title)

The webapp to keep track of regional tennis matches (https://tennis.laurinnotemann.dev)

### Prerequisites

- Docker 
- Go (v1.20)

### Dev Setup
1. Create .env and copy .env.example in there in specify own values
```bash
cp .env.example .env
```

2. Run tests (docker has to be running because it will temporarily build a test db)
```bash
make test
```

3. Run in development (this will create dev db in docker and apply all migrations to it)
```bash
make start-dev-env
```

4. Open Url
http://localhost:3000/

### Technoligies
- Golang
- Echo (Go backend framework)
- Postgres
- sqlc (Codegen from SQL migrations and queries)
- HTML, CSS, JavaScript (no framework)

### Missing features 
- [ ] If team already exsits (returns 500 error rn) - Backend and Frontend
- [ ] warning when you delete a player if has teams, stats or matches or anything else - Backend and Frontend
- [ ] warning for team as well - Backend and Frontend
- [ ] reset password for users - frontend and mail service
- [ ] delete user - frontend
- [ ] edit team <- backlog for now - frontend
- [ ] fetch player name in /edit-player html page (id from url) instead of putting in into the local storage lol - frontend
- [ ] filter matches by teams/players
- [ ] loading spinner or something instead od loading every thing indepenetly - frontend

## Contribution
Everything built by me (Laurin Notemann)
