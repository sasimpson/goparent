#!/bin/sh

docker run --name goparent_rethinkdb --network goparent --rm -d --network goparent -p 28015:28015 -p 29015:29015 -p 8080:8080 -v $PWD/rethinkdb_data:/data/rethinkdb_data rethinkdb:latest