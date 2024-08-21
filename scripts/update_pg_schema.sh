DB_NAME=chase-data

echo "updating pg db schema..."
if psql -d $DB_NAME -a -f internal/db/schema.sql
then
    echo "done"
else
    exit 1
fi
