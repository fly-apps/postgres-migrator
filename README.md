# postgres-migrator

## How it works

Performs a `pg_dump` against a specified source database and restores it into the specified target.


1. Clone this project and launch a new app.
```
fly launch --name postgres-migration-tool
```

2. Specify the SOURCE_DATABASE_URI and TARGET_DATABASE_URI as secrets.
```shell
# The source postgres database.
fly secrets set SOURCE_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>/<database>

# The target postgres database.
fly secrets set TARGET_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>/<database>

```

3. Provision a new volume to hold your migration data.
```
  fly volumes create migration_data --region iad --size 10
```

4. Update the Dockerfile's PG_VERSION to match the Postgres version used by the target. 

5. Deploy the app
```
fly deploy .
```

6. Initiate the migration process
```
fly ssh console -C "migrate"
```

