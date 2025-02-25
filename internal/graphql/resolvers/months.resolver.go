package resolvers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/proctorinc/banker/internal/auth"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
)

type MonthsResolver struct {
	months []string
}

var monthNames = map[string]string{
	"1":  time.January.String(),
	"2":  time.February.String(),
	"3":  time.March.String(),
	"4":  time.April.String(),
	"5":  time.May.String(),
	"6":  time.June.String(),
	"7":  time.July.String(),
	"8":  time.August.String(),
	"9":  time.September.String(),
	"10": time.October.String(),
	"11": time.November.String(),
	"12": time.December.String(),
}

func (r *queryResolver) Months(ctx context.Context) ([]gen.MonthItem, error) {
	user := auth.GetCurrentUser(ctx)
	result := []gen.MonthItem{}

	months, err := r.Repository.ListMonths(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	for _, month := range months {
		firstDay, lastDay, err := getFirstAndLastDayOfMonth(month.Year, month.Month)

		if err == nil {
			result = append(result, gen.MonthItem{
				ID:    strings.ToLower(monthNames[month.Month][:3] + "-" + month.Year),
				Name:  monthNames[month.Month],
				Year:  month.Year,
				Start: firstDay.Format(time.RFC3339),
				End:   lastDay.Format(time.RFC3339),
			})
		}
	}
	return result, nil
}

func getFirstAndLastDayOfMonth(year, month string) (time.Time, time.Time, error) {
	intYear, err := strconv.Atoi(year)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year: %w", err)
	}

	intMonth, err := strconv.Atoi(month)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month: %w", err)
	}

	firstDay := time.Date(intYear, time.Month(intMonth), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	return firstDay, lastDay, nil
}
