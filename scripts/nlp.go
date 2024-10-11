package main

import (
	"fmt"
	"regexp"
)

var descriptions = []string{
	"CAPITAL ONE      MOBILE PMT 3Y9E AIPUVP8EKFQ WEB ID: 9279744380",
	"Zelle payment from JEREMY T LOCK E 22228690365",
	"NEWREZ-SHELLPOIN ACH PMT PPD ID: 6371542226",
	"WINCO #014 3025 SW CED BEAVERTON OR  826828  09/30",
	"WINCO #014 3025 SW CED BEAVERTON OR  756262  08/08",
	"OREGON HEALTH SC SALARY PPD ID: 1931176109",
	"A JESUS CHURCH F A JESUS CH WEB ID: 1201681064",
	"C27840 SECURE CO DIR DEP PPD ID: 4462283648",
	"Zelle payment to Moother 2158001 4472",
	"Zelle payment to Moother 2126372 2103",
	"NEWREZ-SHELLPOIN ACH PMT PPD ID: 6371542226",
	"Payment to Chase card ending in 6268",
	"SCHWAB BROKERAGE MONEYLINK  5586 22442572234 WEB ID: 9005586224",
	"IRS              USATAXPYMT PPD ID: 3387702000",
}

func main() {
	fmt.Println("Running NLP test")
	for _, description := range descriptions {
		id, err := parseMerchantId(description)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(id)
		}
		// fmt.Println(description)
		// // Create a new document with the default configuration:
		// doc, err := prose.NewDocument(description)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// for _, tok := range doc.Entities() {
		// 	fmt.Println(tok.Text, tok.Label)
		// }

		// // Iterate over the doc's tokens:
		// for _, tok := range doc.Tokens() {
		// 	fmt.Println(tok.Text, tok.Tag)
		// }
		// fmt.Println()
	}
}

func parseMerchantId(description string) (string, error) {
	expression, err := regexp.Compile(`\b\d{10}\b`)

	if err != nil {
		return "", err
	}

	matches := expression.FindAllString(description, -1)

	if len(matches) > 0 {
		// Get the last match
		return matches[len(matches)-1], nil
	} else {
		return "", fmt.Errorf("no merchantId found")
	}
}
