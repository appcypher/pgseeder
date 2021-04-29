<h1 align="center">PGSEEDER</h1>

<p align="center">
:warning:  This project is experimental and in active development  :warning:
</p>

A Seeder for Postgres Databases.

### USAGE

Clone pgseeder repo and cd into the created folder

```sh
git clone https://github.com/gigamono/pgseeder
cd pgseeder
```

Build the binary.

```sh
go build cmd/pgseeder.go
```

Add the binary to system path and run `pgseeder` command in a directory with seeds or specify with seed directory with `--dir` flag

```sh
pgseeder -c "postgres://appcypher@localhost:5432/resourcedb?sslmode=disable" --add users
pgseeder -c "postgres://appcypher@localhost:5432/resourcedb?sslmode=disable" --remove users
pgseeder -c "postgres://appcypher@localhost:5432/resourcedb?sslmode=disable" --add-all -d internal/db/seeds/resource
pgseeder -c "postgres://appcypher@localhost:5432/resourcedb?sslmode=disable" --remove-all
```

### LIMITATIONS

Only supports .sql files right now.
