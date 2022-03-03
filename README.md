# postgres-migrator

## How it works

Performs a `pg_dump` against a source database and restores it into the specified target.

## Performing the migration

1. Clone this project and update the Dockerfile's `PG_VERSION` to match the version of your target Postgres cluster.

2. Specify the `SOURCE_DATABASE_URI` and `TARGET_DATABASE_URI` as secrets.
```shell
# The source postgres database.
fly secrets set SOURCE_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>/<database>

# The target postgres database ( Leave off the database name unless the database already exists ).
fly secrets set TARGET_DATABASE_URI=postgres://<user>:<password>@<hostname>:<port>
```

3. Launch the app
```
fly launch --name postgres-migration-tool
```

4. Read the recommendations section below. :)

5. Initiate the migration process
```
fly ssh console -C "migrate"
```

## Recommendations

* Make sure target app has plenty of memory.  
* Scale your destination app down to 1 before initiating the migration.  
* If the migration looks like it has frozen, connect to your target postgres app and run  `\l+` a few times to make sure the data size is still growing.
* If migration connection is broken, your target app is likely under scaled.  Monitor your metrics dashboard and your Fly app's vm checks for guidance.  
