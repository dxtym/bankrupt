#!/bin/sh

#exit on error
set -e

#migrate up
echo "run db migration"
/app/migrate -path /app/migration -database "$SOURCE" -verbose up

#start app
echo "start app"
exec "$@"
