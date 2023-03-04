# postgres-migrator

**This project is deprecated in favor of: [Postgers Importer](https://github.com/fly-apps/postgres-importer)**



## How it works

Performs a `pg_dump` against a source database and restores it into the specified target.

## Performing the migration

1. Clone this project and update the Dockerfile's `PG_VERSION` to match the version of your target Postgres cluster.

2. Run `fly apps create --name postgres-migration-tool`

3. Specify the `SOURCE_DATABASE_URI` and `TARGET_DATABASE_URI` as secrets.
```shell
# The source postgres database.
fly secrets set SOURCE_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>/<database>

# The target postgres database ( Leave off the database name unless the database already exists ).
fly secrets set TARGET_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>
```

4. Deploy app
```
fly deploy
```

5. Read the recommendations section below. :)

6. Initiate the migration process
```
fly ssh console -C "migrate"
```

## Recommendations


* Make sure the target app has plenty of memory.
* Scale your destination app down to 1 before initiating the migration: `fly scale count 1 -a target-db-name`. It might warn you that this app name is different than what's in the current fly.toml, but that's ok
* If the migration looks like it has frozen, connect to your target postgres app and run  `\l+` a few times to make sure the data size is still growing.
* If migration connection is broken, your target app is likely under scaled.  Monitor your metrics dashboard and your Fly app's vm checks for guidance.
