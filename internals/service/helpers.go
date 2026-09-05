package service

import (
	"F5/internals/apperrors"
	"time"
)

func parseDate(dateStr *string) (*time.Time, error) {
	if dateStr == nil {
		return nil, nil
	}

	parsedDate, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return nil, apperrors.ErrInvalidDueDate
	}

	return &parsedDate, nil
}

func unParseDate(dateTime *time.Time) *string {
	if dateTime == nil {
		return nil
	}

	unparsedDate := dateTime.Format("2006-01-02")

	return &unparsedDate
}
