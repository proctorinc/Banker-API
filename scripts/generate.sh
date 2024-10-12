printf "generating sqlc models..."
if sqlc generate
then
    echo " [done]"
else
    echo " [failed]"
    exit 1
fi

printf "removing old dataloaders..."
if (rm $PWD/internal/dataloaders/*_gen.go)
then
    echo " [done]"
fi

printf "generating new dataloaders..."
if go generate ./internal/dataloaders/...
then
    echo " [done]"
else
    echo " [failed]"
    exit 1
fi

printf "removing old graphql models..."
if (rm -r $PWD/internal/graphql/generated)
then
    echo " [done]"
fi

printf "generating new gql models..."
if go run github.com/99designs/gqlgen generate
then
    echo " [done]"
else
    echo " [failed]"
    exit 1
fi
