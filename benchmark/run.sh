#!/bin/bash
wrk -c10 -d20 -t6 -s wrk.lua http://localhost:8080/v0/hit/TST8000
# roboto -b payload -x POST -n 50000 -t 1m -u http://localhost:8080/v0/hit/TST8000 -r 0