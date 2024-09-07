echo "generating sqlc files..."
if sqlc generate
then
    echo "done"
else
    exit 1
fi

echo "removing graphql models file"
if (rm $PWD/internal/graphql/models.go && rm $PWD/internal/graphql/models.go) || true
then
    echo "done"
fi

echo "generating gql files..."
if go run github.com/99designs/gqlgen generate
then
    echo "done"
else
    exit 1
fi
