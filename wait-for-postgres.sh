#!/bin/sh
# wait-for-postgres.sh

set -e

pass="$1"; shift
host="$1"; shift
port="$1"; shift
user="$1"; shift
cmd="$@"

until PGPASSWORD="$pass" psql -h "$host" -p "$port" -U "$user" -c '\q'; do
    >&2 echo "Postgres is unavailable - sleeping"

    sleep 2
done

>&2 echo "Postgres is up - executing command"

exec $cmd
