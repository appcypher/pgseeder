<h1 align="center">PGSEEDER</h1>

A seeder for adding and removing seeds from a Postgres database.

## <sup><sup>🛠</sup></sup> INSTALLATION

- Clone pgseeder repo and cd into the created folder

  ```sh
  git clone https://github.com/gigamono/pgseeder
  ```

  ```sh
  cd pgseeder
  ```

- Build the binary.

  ```sh
  go build cmd/pgseeder.go
  ```

- Add the binary to system path and run `pgseeder -h` command

## <sup><sup>📜</sup></sup>  USAGE

### CREATING SEED FILES

- Pgseeder expects `.sql` filenames to be ordered based on their dependency.

  The ordering format has to be a prefix and an integer.

  For example:

  ```
  01_users.sql
  ```

  ```
  20210702_users.sql
  ```

- You can use Pgseeder to create a timestamped `.sql` file.

  ```sh
  pgseeder --create "users"
  ```

  This generates a timestamped `.sql` file similar to this: `20210620150637_users.sql`

> <sup><sup>ℹ️</sup></sup> Ordering allows Pgseeder to avoid foreign key constraint errors when adding or removing all seeds.

### ADDING SEEDS TO THE DATABASE

- To add seeds from a `20210620150637_users.sql` file, you run following command:

  ```sh
  pgseeder -c "postgres://postgres@localhost:5432/resourcedb?sslmode=disable" --add users
  ```

  The `-c` is needed to connect to the database. It takes postgres connection string.

- You can also run all the `.sql` files in the current directory with the `--add-all` flag

  ```sh
  pgseeder -c "postgres://postgres@localhost:5432/resourcedb?sslmode=disable" --add-all
  ```

- You can specify a directory where the `.sql` files live with `-d` flag

  ```sh
  pgseeder -c "postgres://postgres@localhost:5432/resourcedb?sslmode=disable" --add-all -d "./seeds"
  ```

> <sup><sup>ℹ️</sup></sup> To support removing seeds later, Pgseeder expects the `.sql` filename to correspond to seeded table's name.


> <sup><sup>⚠</sup></sup> If multiple tables are seeded in a `.sql` file, Pgseeder will only recognise the one that matches the filename.


### REMOVING SEEDS FROM THE DATABASE

- Here is how to remove seeds from `users` table.

  ```sh
  pgseeder -c "postgres://postgres@localhost:5432/resourcedb?sslmode=disable" --remove users
  ```

- To remove all seeds.

  ```sh
  pgseeder -c "postgres://postgres@localhost:5432/resourcedb?sslmode=disable" --remove-all
  ```

### SEED KEYS

- By default Pgseeder stores the `id` (if there is one) of the seed which it uses later to remove it.

- Sometimes the primary key is not `id`. You can specify the primary key with `-k` flag.

  ```sh
  pgseeder -c "postgres://postgres@localhost:5432/resourcedb?sslmode=disable" --add users -k "user_id"
  ```

- And to uniquely identify seeds using composite keys, you separate them by colons, commas or spaces.

  ```sh
  pgseeder -c "postgres://postgres@localhost:5432/resourcedb?sslmode=disable" --add x_memberships -k "user_id, workspace_id"
  ```


## <sup><sup>💀</sup></sup> CURRENT LIMITATIONS

- Only supports `.sql` files.

- Does not support seeding multiple tables in a single `.sql` file.
