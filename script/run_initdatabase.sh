#!/bin/bash
echo "DROP DATABASE benchmark; 
CREATE USER taiptht; 
CREATE DATABASE benchmark;
GRANT ALL PRIVILEGES ON DATABASE benchmark TO taiptht;"|  cockroach sql --insecure --host=localhost:8000