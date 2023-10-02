#!/bin/bash

# Define database connection parameters
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_NAME="postgres"

# Path to your SQL files
SQL_FILE1="data.sql"
SQL_FILE2="functions.sql"

# Run psql with the provided parameters
export PGPASSWORD="postgres"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -v common_path="'$(pwd)/data/'" -f $SQL_FILE1 -f $SQL_FILE2
