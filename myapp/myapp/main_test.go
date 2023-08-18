package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsValidPhoneNumber(t *testing.T) {
	validNumbers := []string{"1234567890", "9876543210"}
	invalidNumbers := []string{"12345", "abcdefg"}

	for _, num := range validNumbers {
		if !isValidPhoneNumber(num) {
			t.Errorf("Expected %s to be a valid phone number", num)
		}
	}

	for _, num := range invalidNumbers {
		if isValidPhoneNumber(num) {
			t.Errorf("Expected %s to be an invalid phone number", num)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	validEmails := []string{"test@example.com", "user.name123@gmail.com"}
	invalidEmails := []string{"notanemail", "user@.com", "user@domain"}

	for _, email := range validEmails {
		if !isValidEmail(email) {
			t.Errorf("Expected %s to be a valid email", email)
		}
	}

	for _, email := range invalidEmails {
		if isValidEmail(email) {
			t.Errorf("Expected %s to be an invalid email", email)
		}
	}
}

func TestSubmitHandler(t *testing.T) {
	// Create a test HTTP request with form data
	payload := strings.NewReader("full_name=John+Doe&gender=Male&from_date=2023-08-01&to_date=2023-08-10&phone_number=1234567890&email=test@example.com&message=Hello")
	req := httptest.NewRequest("POST", "/submit", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	// Call the handler
	submitHandler(rec, req)

	// Check the response status code and content
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rec.Code)
	}
	expectedResponse := "Message submitted successfully!"
	if body := rec.Body.String(); body != expectedResponse {
		t.Errorf("Expected response %s, but got %s", expectedResponse, body)
	}
}
