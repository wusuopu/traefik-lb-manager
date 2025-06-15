#!/bin/sh


start_server () {
  ./app
}

init_db () {
  if [[ "$DATABASE_TYPE" == "sqlite" ]]; then
    ./goose db:create
  fi
  ./goose up
}

if [ "$1" = 'start_server' ]; then
  start_server
elif [ "$1" = 'init_db' ]; then
  init_db
else
  echo $@
  exec $@
fi
