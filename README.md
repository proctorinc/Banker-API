# Banker API
Banker is a GraphQL banking API for fetching bank account and transaction data and statistics. Uses OFX transaction files to upload data from common banks like Chase.

## Libraries
- gin-gonic > HTTP web framework (https://github.com/gin-gonic/gin)
- gqlgen > GraphQL code generation (https://github.com/99designs/gqlgen)
- dataloaden > GraphQL dataloader code generation (https://github.com/vektah/dataloaden)
- sqlc > SQL model/query code generation (https://github.com/sqlc-dev/sqlc)


## Setup
#### Prereqs
- Brew (mac only) (https://mac.install.guide/homebrew/3)
- Go (https://go.dev/doc/install)
- Air (https://github.com/air-verse/air)

## How to run (development server)
To start the server, run Air (May depend upon where you have air installed)
```sh
~/.air
```

## Latest Updates
- Added dataloaders to efficiently query and cache data for nested subqueries in large queries
- Query cursor pagination via GraphQL edges and nodes
- File upload that supports OFX transaction files to retrieve transaction and account data
- Merchant matching by using a Natural Language Model (NLM) by parsing merchant data from transaction descriptions

## Example Queries
### Accounts data query
```graphql
query accounts {
  accounts(page:{
    first: 2
    after: "Y3Vyc29yOm9mZnNldDox"
  }) {
    edges {
      node {
        id
        uploadSource
        type
        name
        routingNumber
        transactions(page: {
          first: 4
          after: "Y3Vyc29yOm9mZnNldDox"
        }) {
          edges {
            node {
              description
              amount
              date
            }
            cursor
          }
        }
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      totalCount
    }
  }
}
```

### Transaction data query
```graphql
query transactions {
  transactions(page: {
    after: "Y3Vyc29yOm9mZnNldDo2"
    first: 10
  }) {
    edges {
      node {
        id
        description
        amount
        payee
        date
        updated
        merchant {
          name
          sourceId
        }
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      totalCount
    }
  }
}
```

### Merchants data query
```graphql
query merchants {
  merchants(page:{
    first: 10
  }) {
    edges {
      node {
        id
        name
        sourceId
        ownerId
        transactions(page:{
          first: 2
        }) {
          edges {
            node {
              description
              payee
              amount
              date
            }
          }
        }
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      totalCount
    }
  }
}
```

## Development
To generate sqlc models, graphql models, and graphql dataloaders, run:
```sh
./scripts/generate.sh
```
This command runs ```sqlc generate```, ```gqlgen generate```, ```go generate ./internal/dataloaders/...``` to compile sql and graphql schema files into go
