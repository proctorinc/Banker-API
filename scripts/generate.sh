DB_NAME=chase-data

echo "generating sqlc files..."
sqlc generate
echo "generating gql files..."
go run github.com/99designs/gqlgen generate
echo "pushing changes to db..."
psql -d $DB_NAME -a -f internal/db/schema.sql
echo "done"
