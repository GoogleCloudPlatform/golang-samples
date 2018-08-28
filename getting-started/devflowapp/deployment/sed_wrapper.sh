#!/bin/bash

# Substitue values for environment variables in file $1

echo "Substituting values for password and project ($PROJECT)"
sed -e "s/_DB_PASSWORD/$_DB_PASSWORD/g" -i $1
sed -e "s/PROJECT_ID/$PROJECT/g" -i $1

