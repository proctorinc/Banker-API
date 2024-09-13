package resolvers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/chase"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
)

type transactionResolver struct{ *Resolver }

func (r *transactionResolver) ID(ctx context.Context, transaction *db.Transaction) (uuid.UUID, error) {
	return transaction.ID, nil
}

func (r *transactionResolver) SourceId(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Sourceid, nil
}

func (r *transactionResolver) UploadSource(ctx context.Context, transaction *db.Transaction) (string, error) {
	return string(transaction.Uploadsource), nil
}

func (r *transactionResolver) Amount(ctx context.Context, transaction *db.Transaction) (float64, error) {
	return float64(transaction.Amount) / 100, nil
}

func (r *transactionResolver) PayeeID(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Payeeid.String) > 0 {
		return &transaction.Payeeid.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) Payee(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Payee.String) > 0 {
		return &transaction.Payee.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) PayeeFull(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Payeefull.String) > 0 {
		return &transaction.Payeefull.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) IsoCurrencyCode(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Isocurrencycode, nil
}

func (r *transactionResolver) Date(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Date.Format(time.RFC3339), nil
}

func (r *transactionResolver) Description(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Description, nil
}

func (r *transactionResolver) Type(ctx context.Context, transaction *db.Transaction) (string, error) {
	return string(transaction.Type), nil
}

func (r *transactionResolver) CheckNumber(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Checknumber.String) > 0 {
		return &transaction.Checknumber.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) Updated(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Updated.Format(time.RFC3339), nil
}

// Queries

func (r *queryResolver) Transaction(ctx context.Context, transactionId uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.GetTransaction(ctx, transactionId)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *queryResolver) Transactions(ctx context.Context) ([]db.Transaction, error) {
	user := auth.GetCurrentUser(ctx)

	return r.Repository.ListTransactions(ctx, user.ID)
}

// Mutations

func (r *mutationResolver) CreateTransaction(ctx context.Context, data gen.TransactionInput) (*db.Transaction, error) {
	user := auth.GetCurrentUser(ctx)

	transaction, err := r.Repository.CreateTransaction(ctx, db.CreateTransactionParams{
		Ownerid: user.ID,
		Amount:  int32(data.Amount * 100),
	})

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *mutationResolver) DeleteTransaction(ctx context.Context, id uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.DeleteTransaction(ctx, id)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *mutationResolver) ChaseCSVTransactionsUpload(ctx context.Context, reader graphql.Upload) (bool, error) {
	if !bytes.HasSuffix([]byte(reader.Filename), []byte(".csv")) {
		log.Printf("Invalid extension: %s", reader.Filename)
		return false, fmt.Errorf("Invalid file extension. .CSV required")
	}

	user := auth.GetCurrentUser(ctx)
	transactions, err := chase.ParseChaseCSV(reader.File)

	if err != nil {
		return false, err
	}

	for _, transaction := range transactions {
		r.Repository.CreateTransaction(ctx, db.CreateTransactionParams{
			Ownerid: user.ID,
			Amount:  int32(transaction.Amount * 100),
		})
	}

	return true, nil
}

func (r *mutationResolver) ChaseOFXTransactionsUpload(ctx context.Context, reader graphql.Upload) (bool, error) {
	log.Printf("File extension: %s", strings.ToLower(reader.Filename))
	if !bytes.HasSuffix([]byte(reader.Filename), []byte(".ofx")) &&
		!bytes.HasSuffix([]byte(reader.Filename), []byte(".QFX")) &&
		!bytes.HasSuffix([]byte(reader.Filename), []byte(".")) {
		log.Printf("Invalid extension: %s", reader.Filename)
		return false, fmt.Errorf("Invalid file extension. .OFX/.QBX/.QBO required")
	}

	user := auth.GetCurrentUser(ctx)
	ofxResult, err := chase.ParseChaseOFX(reader.File)
	ofxAccount := ofxResult.Account

	if err != nil {
		return false, err
	}

	account, err := r.Repository.UpsertAccount(ctx, db.UpsertAccountParams{
		Sourceid:      ofxAccount.AccountId,
		Uploadsource:  db.UploadSourceCHASEOFXUPLOAD,
		Name:          ofxAccount.Name,
		Type:          db.AccountType(ofxAccount.Type),
		Routingnumber: sql.NullString{String: ofxAccount.BankId, Valid: len(ofxAccount.BankId) > 0},
		Ownerid:       user.ID,
	})

	if err != nil {
		return false, err
	}

	log.Printf("Resolver transactions: %d", len(ofxResult.Transactions))

	for _, tx := range ofxResult.Transactions {
		r.Repository.UpsertTransaction(ctx, db.UpsertTransactionParams{
			Ownerid:         user.ID,
			Amount:          int32(tx.Amount * 100),
			Payeeid:         sql.NullString{String: tx.PayeeId, Valid: len(tx.PayeeId) > 0},
			Payee:           sql.NullString{String: tx.Payee, Valid: len(tx.Payee) > 0},
			Payeefull:       sql.NullString{String: tx.PayeeFull, Valid: len(tx.PayeeFull) > 0},
			Sourceid:        tx.Id,
			Uploadsource:    db.UploadSourceCHASEOFXUPLOAD,
			Isocurrencycode: ofxResult.Account.IsoCurrencyCode,
			Date:            tx.DatePosted,
			Description:     tx.Description,
			Type:            db.TransactionType(tx.Type),
			Updated:         time.Now(),
			Checknumber:     sql.NullString{String: tx.CheckNumber, Valid: len(tx.CheckNumber) > 0},
			Accountid:       account.ID,
		})
	}

	return true, nil
}
