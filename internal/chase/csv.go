package chase

import (
	"encoding/csv"
	"io"
	"time"

	"github.com/gocarina/gocsv"
)

type Date struct {
	time.Time
}

type ChaseCSVTransaction struct {
	Details        string  `csv:"Details"`
	PostingDate    Date    `csv:"Posting Date"`
	Description    string  `csv:"Description"`
	Amount         float32 `csv:"Amount"`
	Type           string  `csv:"Type"`
	Balance        float32 `csv:"Balance"`
	CheckOrSlipNum string  `csv:"Check or Slip #"`
}

func (date *Date) MarshalCSV() (string, error) {
	return date.Format("01/02/2006"), nil
}

func (date *Date) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("01/02/2006", csv)
	return err
}

func ParseChaseCSV(reader io.Reader) ([]ChaseCSVTransaction, error) {
	transactions := []ChaseCSVTransaction{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.FieldsPerRecord = -1
		return r
	})

	if err := gocsv.Unmarshal(reader, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
