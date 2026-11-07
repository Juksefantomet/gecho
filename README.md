# Gecho

> *"Time is the most valuable commodity."* — Gordon Gekko

**Gecho** is a Rails-inspired scaffolding and migrations CLI toolset for Go projects using:

- [Echo](https://echo.labstack.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)


It helps you bootstrap new apps with more haste by generating templates for models, routes, queries, migrations, and project structure in seconds.

Disclaimer: This is probably highly opinionated and might not match your expectations on how to structure your project, but it might get you started if you don't have one already.

---

## Installation

Install the CLI globally with:

```bash
go install github.com/Juksefantomet/gecho@latest
```

Install tag specific version:

```bash
go install github.com/Juksefantomet/gecho@v0.1.5
```

---

## Commands

### `gecho init`

Creates a Gecho-compatible folder structure.

If no `go.mod` exists, you will be prompted to enter a module name.

It also generates:
- `.env` file (PostgreSQL connection)
- `main.go` with Echo + Swagger boilerplate
- `app/routes/helloWorld.go`
- `app/services/database/database.go`
- `app/services/database/migrations.go` with `database.AutoMigrate()`, which applies pending `db/migrations` on app start. The call is commented out in `main.go`; uncomment it to enable.

```bash
gecho init
```

---

### `gecho scaffold <name>`

Follows Rails naming: `<name>` may be singular or plural (`user`, `users`, `BlogPost`). The model is singular, the table and route are plural.

Generates (for `gecho scaffold user`):
- Model: `app/models/user.go` (`type User`)
- Route: `app/routes/users.go` (`GetUsers`, `GET /users`)
- Query: `app/services/database/userQueries.go` (`SELECT * FROM users`)
- Migrations: `db/migrations/*.up.sql` + `.down.sql` (`CREATE TABLE users`)

**Note**: You must wire the route manually in `main.go`. There is a commented example present in the generated main.go

```bash
gecho scaffold user
```

---

### `gecho create-migration <name>`

Creates a pair of empty `.up.sql` and `.down.sql` migration files.

```bash
gecho create-migration add_index_to_users
```

---

### `gecho migrate [down]`

Applies or rolls back migrations defined in `db/migrations`.

```bash
gecho migrate         # applies all .up.sql files
gecho migrate down    # rolls back the last applied migration in migrations table in the database
```

---

### `gecho automigrate [generate]`

Runs pending migrations when the app starts, using the same `migrations` table as `gecho migrate`.

```bash
gecho automigrate            # shows setup instructions
gecho automigrate generate   # creates app/services/database/migrations.go if missing
```

Then add this to `main.go`, right after the existing `database.InitDB()`:

```go
if err := database.AutoMigrate(); err != nil {
	log.Fatalf("Auto-migrate failed: %v", err)
}
```

---

## After `gecho init`

Run this to get everything working: (first time users of swagger must install the binary)

```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.4
swag init
go mod tidy
go run main.go
```

`swag init` must run before `go mod tidy`: the generated `main.go` imports the `docs` package that `swag init` creates.

NOTE:
Installing swag@latest can cause issue when running go run main.go - lock in on version


Then visit:

```
http://localhost:3000/swagger/index.html
```

---

## Releasing

Releases are based on the version in the `VERSION` file (e.g., `0.1.0`).

To tag and push a release:

```bash
./release.sh
```

Then draft a GitHub release at:

```
https://github.com/Juksefantomet/gecho/releases/new?tag=v0.1.0
```

---

## Quote

> *"Time is the most valuable commodity."* — Gordon Gekko

Gecho was built with speed, repetition, and scale in mind.

---


## License

Licensed under the [MIT License](LICENSE).
