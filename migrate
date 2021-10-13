#!/bin/bash

echo "Starting migration"
pg_dump -d $SOURCE_DATABASE_URI -Fc > /data/dump.sql

pg_restore --clean -d $TARGET_DATABASE_URI < /data/dump.sql
pg_restore --create -d $TARGET_DATABASE_URI < /data/dump.sql

echo "All done!"

