package chase

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/aclindsa/ofxgo"
	"github.com/proctorinc/banker/internal/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ChaseOFXAccount struct {
	BankId           string
	AccountId        string
	IsoCurrencyCode  string
	Type             db.AccountType
	CurrentBalance   float32
	AvailableBalance float32
	Name             string
}

type ChaseOFXTransaction struct {
	Id          string
	Type        string
	DatePosted  time.Time
	Amount      float32
	PayeeId     string
	Payee       string
	PayeeFull   string
	CheckNumber string
	Description string
}

type ChaseOFXResult struct {
	Account      ChaseOFXAccount
	Transactions []ChaseOFXTransaction
}

func ParseChaseOFX(reader io.Reader) (*ChaseOFXResult, error) {
	response, err := ofxgo.ParseResponse(reader)

	if err != nil {
		log.Println("Failed to parse OFX file")
		return nil, fmt.Errorf("Failed to parse OFX file: %w", err)
	}

	// Was there an OFX error while processing our request?
	if response.Signon.Status.Code != 0 {
		meaning, _ := response.Signon.Status.CodeMeaning()
		return nil, fmt.Errorf("Nonzero signon status (%d: %s) with message: %s", response.Signon.Status.Code, meaning, response.Signon.Status.Message)
	}

	if len(response.Bank) > 0 {
		return parseBankAccount(response)
	} else if len(response.CreditCard) > 0 {
		return parseCreditCard(response)
	}

	return nil, fmt.Errorf("Unsupported account type. Supported: Bank account, credit card")
}

func parseBankAccount(response *ofxgo.Response) (*ChaseOFXResult, error) {
	if stmt, ok := response.Bank[0].(*ofxgo.StatementResponse); ok {
		current, _ := stmt.BalAmt.Float32()
		available, _ := stmt.AvailBalAmt.Float32()
		accountId := stmt.BankAcctFrom.AcctID.String()
		accountType := db.AccountType(stmt.BankAcctFrom.AcctType.String())
		caser := cases.Title(language.AmericanEnglish)
		accountTypeFormatted := caser.String(strings.ToLower(string(accountType)))

		account := ChaseOFXAccount{
			BankId:           stmt.BankAcctFrom.BankID.String(),
			AccountId:        accountId,
			IsoCurrencyCode:  stmt.CurDef.String(),
			Type:             accountType,
			CurrentBalance:   current,
			AvailableBalance: available,
			Name:             fmt.Sprintf("%s Account %s", accountTypeFormatted, accountId[len(accountId)-4:]),
		}
		var transactions []ChaseOFXTransaction

		for _, tx := range stmt.BankTranList.Transactions {
			amount, _ := tx.TrnAmt.Float32()
			name := tx.Name.String()

			if tx.Payee != nil {
				name = tx.Payee.Name.String()
			}

			transaction := ChaseOFXTransaction{
				Id:          tx.FiTID.String(),
				Type:        tx.TrnType.String(),
				DatePosted:  tx.DtPosted.Time,
				Amount:      amount,
				PayeeId:     tx.PayeeID.String(),
				Payee:       name,
				PayeeFull:   tx.ExtdName.String(),
				CheckNumber: tx.CheckNum.String(),
				Description: getDescription(name, tx.Memo.String()),
			}

			transactions = append(transactions, transaction)
		}

		result := &ChaseOFXResult{
			Account:      account,
			Transactions: transactions,
		}

		return result, nil
	}

	return nil, fmt.Errorf("ofx: failed to parse bank account data")
}

func parseCreditCard(response *ofxgo.Response) (*ChaseOFXResult, error) {
	if stmt, ok := response.CreditCard[0].(*ofxgo.CCStatementResponse); ok {
		current, _ := stmt.BalAmt.Float32()
		available, _ := stmt.AvailBalAmt.Float32()
		accountId := stmt.CCAcctFrom.AcctID.String()
		accountType := db.AccountTypeCREDIT
		caser := cases.Title(language.AmericanEnglish)
		accountTypeFormatted := caser.String(strings.ToLower(string(accountType)))

		account := ChaseOFXAccount{
			AccountId:        accountId,
			IsoCurrencyCode:  stmt.CurDef.String(),
			Type:             accountType,
			CurrentBalance:   current,
			AvailableBalance: available,
			Name:             fmt.Sprintf("%s Card %s", accountTypeFormatted, accountId[len(accountId)-4:]),
		}
		var transactions []ChaseOFXTransaction

		for _, tx := range stmt.BankTranList.Transactions {
			amount, _ := tx.TrnAmt.Float32()
			name := tx.Name.String()

			if tx.Payee != nil {
				name = tx.Payee.Name.String()
			}

			transaction := ChaseOFXTransaction{
				Id:          tx.FiTID.String(),
				Type:        tx.TrnType.String(),
				DatePosted:  tx.DtPosted.Time,
				Amount:      amount,
				PayeeId:     tx.PayeeID.String(),
				Payee:       name,
				PayeeFull:   tx.ExtdName.String(),
				CheckNumber: tx.CheckNum.String(),
				Description: getDescription(name, tx.Memo.String()),
			}

			transactions = append(transactions, transaction)
		}

		result := &ChaseOFXResult{
			Account:      account,
			Transactions: transactions,
		}

		return result, nil
	}

	return nil, fmt.Errorf("ofx: failed to parse credit card data")
}

func getDescription(name string, memo string) string {
	if name == "" {
		return memo
	} else if memo == "" {
		return name
	}

	return strings.Trim(name, " ") + " " + strings.Trim(memo, " ")
}
