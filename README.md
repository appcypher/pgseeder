<h1 align="center">PGSEEDER</h1>

A seeder for adding and removing seeds from a Postgres database.

### INSTALLATION

Clone pgseeder repo and cd into the created folder

```sh
git clone https://github.com/gigamono/pgseeder
```

```sh
cd pgseeder
```

Build the binary.

```sh
go build cmd/pgseeder.go
```

Add the binary to system path and run `pgseeder -h` command

### USAGE

##### ADDING SEEDS

For example, you can add seeds from a `user.sql` file with the following command:

```sh
pgseeder -c "postgres://postgres@localhost:5432/resourcedb" --add users
```

The `-c` is needed to connect to the database. It takes postgres connection string.

You can also run all the `.sql` files in the current directory with the `-add-all` flag

```sh
pgseeder -c "postgres://postgres@localhost:5432/resourcedb" --add-all
```

You can specify a directory where the `.sql` files live with `-d` flag

```sh
pgseeder -c "postgres://postgres@localhost:5432/resourcedb" --add-all -d "./seeds"
```

:information_source: To allow removing seeds later, Pgseeder expects the `.sql` filename to correspond to seeded table's name.

:warning: If multiple tables are seeded in a `.sql` file, Pgseeder will only recognise the one that matches the filename.


##### REMOVING SEEDS

Here is how to remove seeds from some `users` table.

```sh
pgseeder -c "postgres://postgres@localhost:5432/resourcedb" --remove users
```

To remove all seeds.

```sh
pgseeder -c "postgres://postgres@localhost:5432/resourcedb" --remove-all
```


##### SEED KEYS

By default Pgseeder stores the `id` (if there is one) of the seed which it uses later to remove them.

Sometimes the primary key is not `id`. You can specify the primary key with `k` flag.

```sh
pgseeder -c "postgres://postgres@localhost:5432/resourcedb" --add users -k "user_id"
```

And to uniquely identify seeds using composite keys, you separate them by colons, commas or spaces.

```sh
pgseeder -c "postgres://postgres@localhost:5432/resourcedb" --add x_memberships -k "user_id, workspace_id"
```


### LIMITATIONS

- Only supports `.sql` files.
- Does not support seeding multiple tables in a single `.sql` files
