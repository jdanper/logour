#!/bin/bash
wrk -c10 -d20 -t6 -s wrk.lua http://localhost:8080/v0/event