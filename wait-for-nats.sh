#!/bin/sh
# wait-for-nats.sh

set -e

host="$1"
shift
cmd="$@"

until nc -z $host 4222; do
  >&2 echo "Nats is unavailable - sleeping"
  sleep 1
done

>&2 echo "Nats is up - executing command"
exec $cmd
