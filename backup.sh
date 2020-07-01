#!/bin/bash

mysqldump --column-statistics=0 -u test -ptest goinagbe | gzip > goinagbe.sql.gz
# scp goinagbe.sql.gz rds:/home/app/goinagbe.sql.gz