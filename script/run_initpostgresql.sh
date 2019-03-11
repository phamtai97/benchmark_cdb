#!/bin/bash
echo "DROP DATABASE benchmark; 
CREATE DATABASE benchmark;" |  sudo -u postgres psql