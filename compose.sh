#!/bin/bash
# docker-compose build --no-cache && docker-compose up -d
# docker-compose down
docker-compose build && docker-compose up -d
docker logs -f goinagbe