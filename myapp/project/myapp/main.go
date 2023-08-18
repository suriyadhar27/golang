package myapp

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

func isValidPhoneNumber(number string) bool {
	// The regular expression to match Indian phone numbers (10 digits starting with 7, 8, or 9)
	match, _ := regexp.MatchString(`^(?:\+91|0?91|\d{0,2}-?)?[7-9]\d{9}$`, number)
	return match
}

//func isValidPhoneNumber(number string) bool {
// A simpler regular expression to match 10-digit Indian phone numbers that start with 7, 8, or 9
//match, _ := regexp.MatchString(`^[789]\d{9}$`, number)
//return match
//}

// isValidEmail checks if a given string is a valid business email address.
func isValidEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.(com|net|org|edu|biz|info|io|your-extension-here)$`

	match, _ := regexp.MatchString(emailPattern, email)
	return match
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	// Handle GET request logic here
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "GET request received")
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	fullName := r.PostFormValue("full_name")
	phoneNumber := r.PostFormValue("phone_number")
	email := r.PostFormValue("email")
	fromDate := r.PostFormValue("FromDate")
	toDate := r.PostFormValue("ToDate")

	if len(fullName) == 0 || len(phoneNumber) == 0 || len(email) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	//Comment out the email validation for testing purposes
	if !isValidEmail(email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// validate date
	dateComparisonErr := CompareDates(fromDate, toDate)
	if dateComparisonErr != nil {
		http.Error(w, dateComparisonErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data submitted successfully"))
}

// HandleRequest handles incoming HTTP requests.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetRequest(w, r)
	case http.MethodPost:
		handlePostRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

var ErrBadRequest = errors.New("bad request")

func CompareDates(fromDate, toDate string) error {
	fromDateParsed, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return err
	}
	toDateParsed, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return err
	}
	if toDateParsed.Before(fromDateParsed) {
		return ErrBadRequest
	}
	return nil
}

func getMessageByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is missing", http.StatusBadRequest)
		return
	}
}
