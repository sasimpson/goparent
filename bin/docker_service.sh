#!/bin/sh

docker run --name service1 --rm -d -p 8000:8000 --network goparent goparent:latest