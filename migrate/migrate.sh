#!/bin/bash

SQL_FILE1="/scripts/data.sql"
SQL_FILE2="/scripts/functions.sql"

psql -v data_path="'$DATA_PATH'" -f "$DATA_PATH$SQL_FILE1" -f "$DATA_PATH$SQL_FILE2"
