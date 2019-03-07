#!/bin/bash
echo "DROP DATABASE benchmark; 
CREATE USER IF NOT EXISTS taiptht; 
CREATE DATABASE benchmark;
GRANT ALL ON DATABASE benchmark TO taiptht;"| cockroach sql --insecure --host=localhost:8000