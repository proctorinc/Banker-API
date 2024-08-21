echo "generating sqlc files..."
if sqlc generate
then
    echo "done"
else
    exit 1
fi

echo "removing graphql models file"
if rm /Users/mattyp/Documents/Projects/go-bank-backend/internal/graphql/models.go
then
    echo "done"
else
    exit 1
fi

echo "generating gql files..."
if go run github.com/99designs/gqlgen generate
then
    echo "done"
else
    exit 1
fi
