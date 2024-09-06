// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"github.com/google/uuid"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Mutation struct {
}

type Query struct {
}

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TransactionInput struct {
	Amount  float64   `json:"amount"`
	OwnerID uuid.UUID `json:"ownerId"`
}
