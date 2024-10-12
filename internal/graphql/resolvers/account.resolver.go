package resolvers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/jdkato/prose/v2"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/chase"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/utils"
)

func (r *accountResolver) ID(ctx context.Context, account *db.Account) (uuid.UUID, error) {
	return account.ID, nil
}

func (r *accountResolver) SourceId(ctx context.Context, account *db.Account) (string, error) {
	masked := utils.MaskData(account.Sourceid)
	return masked, nil
}

func (r *accountResolver) UploadSource(ctx context.Context, account *db.Account) (string, error) {
	return string(account.Uploadsource), nil
}

func (r *accountResolver) Type(ctx context.Context, account *db.Account) (string, error) {
	return string(account.Type), nil
}

func (r *accountResolver) Name(ctx context.Context, account *db.Account) (string, error) {
	return account.Name, nil
}

func (r *accountResolver) RoutingNumber(ctx context.Context, account *db.Account) (*string, error) {
	if len(account.Routingnumber.String) > 0 {
		masked := utils.MaskData(account.Routingnumber.String)
		return &masked, nil
	}

	return nil, nil
}

func (r *accountResolver) Transactions(ctx context.Context, account *db.Account) ([]db.Transaction, error) {
	return r.DataLoaders.Retrieve(ctx).TransactionsByAccountId.Load(account.ID.String())
}

func (r *accountResolver) Stats(ctx context.Context, account *db.Account, input *gen.StatsInput) (*gen.StatsResponse, error) {
	user := auth.GetCurrentUser(ctx)
	incomeTotal, err := r.Repository.GetAccountIncome(ctx, db.GetAccountIncomeParams{
		Ownerid:   user.ID,
		Accountid: account.ID,
	})

	if err != nil {
		return nil, err
	}

	amount := incomeTotal.(int32)

	income := gen.IncomeStats{
		Total:        utils.FormatCurrencyFloat64(amount),
		Transactions: []db.Transaction{},
	}

	spendingTotal, err := r.Repository.GetAccountSpending(ctx, db.GetAccountSpendingParams{
		Ownerid:   user.ID,
		Accountid: account.ID,
	})

	if err != nil {
		return nil, err
	}

	spending := gen.SpendingStats{
		Total:        utils.FormatCurrencyFloat64(spendingTotal.(int32)),
		Transactions: []db.Transaction{},
	}

	total := incomeTotal.(int32) + spendingTotal.(int32)

	net := gen.NetStats{
		Total: utils.FormatCurrencyFloat64(total),
	}

	response := &gen.StatsResponse{
		Spending: &spending,
		Income:   &income,
		Net:      &net,
	}

	return response, nil
}

// Queries

func (r *queryResolver) Account(ctx context.Context, accountId uuid.UUID) (*db.Account, error) {
	user := auth.GetCurrentUser(ctx)
	account, err := r.Repository.GetAccount(ctx, db.GetAccountParams{
		ID:      accountId,
		Ownerid: user.ID,
	})

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *queryResolver) Accounts(ctx context.Context) ([]db.Account, error) {
	user := auth.GetCurrentUser(ctx)

	return r.Repository.ListAccounts(ctx, user.ID)
}

// Mutations

func (r *mutationResolver) ChaseOFXUpload(ctx context.Context, reader graphql.Upload) (*gen.UploadResponse, error) {
	response := &gen.UploadResponse{
		Success: false,
		Accounts: &gen.UploadStats{
			Updated: 0,
			Failed:  0,
		},
		Transactions: &gen.UploadStats{
			Updated: 0,
			Failed:  0,
		},
	}

	if !bytes.HasSuffix([]byte(reader.Filename), []byte(".ofx")) &&
		!bytes.HasSuffix([]byte(reader.Filename), []byte(".QFX")) &&
		!bytes.HasSuffix([]byte(reader.Filename), []byte(".")) {
		log.Printf("Invalid extension: %s", reader.Filename)
		return response, fmt.Errorf("Invalid file extension. .OFX/.QBX/.QBO required")
	}

	user := auth.GetCurrentUser(ctx)

	// Chase QFX contains extra line at beginning of the file
	// This breaks the OFX reader
	// ofxFile, err := skipFirstLine(reader.File)

	// if err != nil {
	// 	return response, err
	// }

	ofxResult, err := chase.ParseChaseOFX(reader.File)

	if err != nil {
		return response, err
	}

	account, err := r.Repository.UpsertAccount(ctx, db.UpsertAccountParams{
		Sourceid:      ofxResult.Account.AccountId,
		Uploadsource:  db.UploadSourceCHASEOFXUPLOAD,
		Name:          ofxResult.Account.Name,
		Type:          db.AccountType(ofxResult.Account.Type),
		Routingnumber: sql.NullString{String: ofxResult.Account.BankId, Valid: len(ofxResult.Account.BankId) > 0},
		Ownerid:       user.ID,
	})

	if err != nil {
		response.Accounts.Updated++
		return response, err
	}

	// Increment successful account upload
	response.Accounts.Updated++

	for _, tx := range ofxResult.Transactions {
		merchant := new(db.Merchant)
		merchantName := parseMerchantName(tx.Description)
		merchantId, err := parseMerchantId(tx.Description)

		if err != nil {
			merchantId = strings.ToUpper(tx.Description)
		}

		res, err := r.Repository.GetMerchantBySourceId(ctx, sql.NullString{String: merchantId, Valid: merchantId != ""})

		if err != nil {
			res, err = r.Repository.GetMerchantByName(ctx, merchantName)

			if err != nil {
				res, err := r.Repository.LinkMerchant(ctx, db.LinkMerchantParams{
					MerchantName: merchantName, //tx.Description,
					KeyMatch:     merchantId,
					UploadSource: db.UploadSourceCHASEOFXUPLOAD,
					SourceId:     sql.NullString{String: merchantId, Valid: merchantId != ""},
					UserId:       user.ID,
				})

				if err != nil {
					log.Println(err)
				}

				merchant = res
			} else {
				merchant = &res
			}
		} else {
			merchant = &res
		}

		if err != nil {
			fmt.Println("Merchant unable to be found!")
		}

		_, err = r.Repository.UpsertTransaction(ctx, db.UpsertTransactionParams{
			Ownerid:         user.ID,
			Amount:          utils.FormatCurrencyInt(tx.Amount),
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
			Merchantid:      merchant.ID,
		})

		// Increment successful transaction upload
		if err != nil {
			response.Transactions.Failed++
			log.Println(tx)
			log.Println(err)
		} else {
			response.Transactions.Updated++
		}
	}

	response.Success = true

	return response, nil
}

func skipFirstLine(reader io.ReadSeeker) (io.ReadSeeker, error) {
	buf := make([]byte, 2)
	_, err := reader.Read(buf)

	if err != nil {
		return nil, err
	}

	return reader, nil
}

func parseMerchantName(description string) string {
	doc, err := prose.NewDocument(description)
	if err != nil {
		return description
	}

	if len(doc.Entities()) > 0 {
		fmt.Printf("Found merchant name?(entity) %s\n", doc.Entities()[0].Text)
		return doc.Entities()[0].Text
	}

	for _, token := range doc.Tokens() {
		if len(token.Tag) >= 2 && (token.Tag[0:2] == "NN" ||
			token.Tag[0:2] == "VB" ||
			token.Tag[0:2] == "MD" ||
			token.Tag[0:2] == "RB") {
			fmt.Printf("Found merchant name?(token) %s\n", token.Text)
			return token.Text
		}
	}

	return description
}

func parseMerchantId(description string) (string, error) {
	expression := regexp.MustCompile(`\b\d{10}\b`)
	matches := expression.FindAllString(description, -1)

	if len(matches) > 0 {
		// Get the last match
		fmt.Printf("Found merchant id! %s\n", matches[len(matches)-1])
		return matches[len(matches)-1], nil
	} else {
		return "", fmt.Errorf("no merchantId found")
	}
}
