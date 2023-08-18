package myapp

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestIsValidPhoneNumber_ValidNumbers(t *testing.T) {
	validNumbers := []string{
		"7005353812",
		"8123456789",
		"9012345678",
	}

	for _, number := range validNumbers {
		isValid := isValidPhoneNumber(number)
		if !isValid {
			t.Errorf("Expected isValidPhoneNumber(%s) to be true, but got false", number)
		}
	}
}

func TestIsValidPhoneNumber_InvalidNumbers(t *testing.T) {
	invalidNumbers := []string{
		"1234567890",
		"6001234567",
		"912345",
	}

	for _, number := range invalidNumbers {
		isValid := isValidPhoneNumber(number)
		if isValid {
			t.Errorf("Expected isValidPhoneNumber(%s) to be false, but got true", number)
		}
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	// Test valid Indian number
	indianNumber := "9876543210"
	if !isValidPhoneNumber(indianNumber) {
		t.Errorf("Expected %s to be a valid Indian phone number, but it's not.", indianNumber)
	}

	// Test valid international number with country code
	internationalNumber := "+14155552671" // Example: US phone number
	if isValidPhoneNumber(internationalNumber) {
		t.Errorf("Expected %s to be an invalid phone number, but it's considered valid.", internationalNumber)
	}
}

func TestIsValidPhoneNumber_EdgeCases(t *testing.T) {
	edgeCases := []string{
		"+9188888",      // Invalid length
		"1234567890",    // Invalid prefix
		"abc1234567",    // Contains non-numeric characters
		"+198765432100", // Invalid length
		" 9876543210",   // Leading space
		"9876543210 ",   // Trailing space

	}

	for _, number := range edgeCases {
		isValid := isValidPhoneNumber(number)
		if isValid {
			t.Errorf("Expected isValidPhoneNumber(%s) to be false, but got true", number)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	validEmails := []string{
		"john.doe@example.com",
		"jane.smith@mycompany.net",
		"info@business.org",
		"contact@school.edu",
		"support@startup.biz",
		"customer@example.info",
		"sales@tech.io",
	}

	invalidEmails := []string{
		"invalid-email",
		"email@example",
		"user@invalid",
		"user@example",
		"user@example.invalid",
	}

	for _, email := range validEmails {
		if !isValidEmail(email) {
			t.Errorf("Expected %s to be a valid email, but got invalid", email)
		}
	}

	for _, email := range invalidEmails {
		if isValidEmail(email) {
			t.Errorf("Expected %s to be an invalid email, but got valid", email)
		}
	}
}

func TestIsValidEmail_InvalidEmail(t *testing.T) {
	invalidEmail := "user@gmail.xon"
	result := isValidEmail(invalidEmail)
	if result {
		t.Errorf("Expected isValidEmail(%s) to be false, but got true", invalidEmail)
	}
}

func TestIsValidEmail_EdgeCases(t *testing.T) {
	edgeCases := []string{
		"user@",             // Invalid domain
		"user@example",      // Incomplete domain
		"user@.com",         // Invalid domain
		" @example.com",     // Leading space
		"user@example.com ", // Trailing space
	}

	for _, email := range edgeCases {
		if isValidEmail(email) {
			t.Errorf("Expected %s to be an invalid email, but got valid", email)
		}
	}
}

func TestHandleGetRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected HTTP status %d, but got %d", http.StatusOK, status)
	}

	expectedResponse := "GET request received"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponse, rr.Body.String())
	}
}

func TestHandlePostRequest_ValidData(t *testing.T) {
	formData := strings.NewReader("full_name=Suriya+Dhar&phone_number=7005353812&email=user@gmail.com&FromDate=2023-08-01&ToDate=2023-08-15")
	req, err := http.NewRequest("POST", "/", formData)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePostRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected HTTP status %d, but got %d", http.StatusOK, status)
	}

	expectedResponse := "Data submitted successfully"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponse, rr.Body.String())
	}
}

func TestHandlePostRequest_InvalidPhoneNumber(t *testing.T) {
	formData := strings.NewReader("full_name=Suriya+Dhar&phone_number=5353813&email=user@gmail.com&FromDate=2023-08-01&ToDate=2023-08-15")
	req, err := http.NewRequest("POST", "/", formData)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePostRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected HTTP status %d, but got %d", http.StatusOK, status)
	}

	expectedResponse := "Data submitted successfully"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponse, rr.Body.String())
	}
}

func TestCompareDates_ValidDates(t *testing.T) {
	err := CompareDates("2023-08-01", "2023-08-15")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestCompareDates_InvalidDates(t *testing.T) {
	err := CompareDates("2023-08-15", "2023-08-01")
	if err != ErrBadRequest {
		t.Errorf("Expected bad request error, but got %v", err)
	}
}

func TestCompareDates_EdgeCases(t *testing.T) {
	// Test equal dates
	err := CompareDates("2023-08-15", "2023-08-15")
	if err != nil {
		t.Errorf("Expected no error for equal dates, but got %v", err)
	}

	// Test same dates
	err = CompareDates("2023-08-15", "2023-08-15")
	if err != nil {
		t.Errorf("Expected no error for same dates, but got %v", err)
	}
}

func TestGetMessageByIDHandler_ValidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/messages?id=48", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getMessageByIDHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}

}

func TestHandlePostRequest_MissingRequiredFields(t *testing.T) {
	// Omitting required fields in the form data
	formData := url.Values{
		"full_name": {"John Doe"},
		"FromDate":  {"2023-08-01"},
		"ToDate":    {"2023-08-10"},
	}

	req, err := http.NewRequest("POST", "/messages48", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePostRequest)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, but got %v", http.StatusBadRequest, rr.Code)
	}

	expectedResponse := "Missing required fields\n" // Add a newline character
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponse, rr.Body.String())
	}
}

func TestGetMessageByIDHandler_MissingID(t *testing.T) {
	req, err := http.NewRequest("GET", "/messages123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getMessageByIDHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, but got %v", http.StatusBadRequest, rr.Code)
	}
}

func TestGetMessageByIDHandler_InvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", " ", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getMessageByIDHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, but got %v", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleRequestGet(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	recorder := httptest.NewRecorder()

	HandleRequest(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	expectedResponse := "GET request received"
	if body := recorder.Body.String(); body != expectedResponse {
		t.Errorf("Expected response '%s', but got '%s'", expectedResponse, body)
	}
}


