package project

import (
	"testing"
)

func TestValidateDateRange(t *testing.T) {
	testCases := []struct {
		fromDate       string
		toDate         string
		expectedErrMsg string
	}{
		{"2023-08-15", "2023-08-20", ""},                                         // Valid date range
		{"2023-08-20", "2023-08-15", "To Date should be greater than From Date"}, // Invalid date range
		{"2023-08-15", "2023-08-15", "To Date should be greater than From Date"}, // Same dates
		{"invalid", "2023-08-20", "Invalid From Date format"},                    // Invalid from date format
		{"2023-08-15", "invalid", "Invalid To Date format"},                      // Invalid to date format
	}

	for _, tc := range testCases {
		err := validateDateRange(tc.fromDate, tc.toDate)
		if tc.expectedErrMsg == "" {
			// Expecting no error
			if err != nil {
				t.Errorf("For dates %s to %s, expected no error but got: %v", tc.fromDate, tc.toDate, err)
			}
		} else {
			// Expecting an error with the specified error message
			if err == nil || err.Error() != tc.expectedErrMsg {
				t.Errorf("For dates %s to %s, expected error message '%s' but got: %v", tc.fromDate, tc.toDate, tc.expectedErrMsg, err)
			}
		}
	}
}
