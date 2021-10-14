# postgres-migrator

## How it works

Performs a `pg_dump` against a source database and `pg_restores` it into the specified target.


## How to perform the migration

1. Clone this project and update the Dockerfile's `PG_VERSION` to match the Postgres version used by the target.

2. Specify the `SOURCE_DATABASE_URI` and `TARGET_DATABASE_URI` as secrets.
```shell
# The source postgres database.
fly secrets set SOURCE_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>/<database>

# The target postgres database.
fly secrets set TARGET_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>

```

3. Launch the app
```
fly launch --name postgres-migration-tool
```

4. Initiate the migration process
```
fly ssh console -C "migrate"
```

