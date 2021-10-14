# postgres-migrator

## How it works

Performs a `pg_dump` against a source database and `pg_restores` it into the specified target.


1. Clone this project and launch a new app.
```
fly launch --name postgres-migration-tool
```

2. Specify the `SOURCE_DATABASE_URI` and `TARGET_DATABASE_URI` as secrets.
```shell
# The source postgres database.
fly secrets set SOURCE_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>/<database>

# The target postgres database.
fly secrets set TARGET_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>

```

3. Update the Dockerfile's PG_VERSION to match the Postgres version used by the target. 

4. Deploy the app
```
fly deploy .
```

5. Initiate the migration process
```
fly ssh console -C "migrate"
```

