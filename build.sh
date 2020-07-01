#!/bin/bash
clear
tag=`./version.sh`
docker build --rm -t figassis/goinagbe:$tag . && docker push figassis/goinagbe:$tag
