package project

import (
	"fmt"
	"time"
)

func validateDateRange(fromDate, toDate string) error {
	fromDateParsed, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return fmt.Errorf("Invalid From Date format")
	}

	toDateParsed, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return fmt.Errorf("Invalid To Date format")
	}

	if toDateParsed.Before(fromDateParsed) || toDateParsed.Equal(fromDateParsed) {
		return fmt.Errorf("To Date should be greater than From Date")
	}

	return nil
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

	// Resume validation
	if !isValidResume(resumeHeader) {
		http.Error(w, "Invalid resume file. Only PDF and PNG files up to 5MB are allowed.", http.StatusBadRequest)
		return
	}

	// Generate a unique filename for the resume
	resumeFileName := fmt.Sprintf("resume_%d%s", time.Now().Unix(), filepath.Ext(resumeHeader.Filename))

	// Save the resume file to the 'resumes' directory
	resumeFilePath := filepath.Join("resumes", resumeFileName)
	outFile, err := os.Create(resumeFilePath)
	if err != nil {
		http.Error(w, "Error saving resume file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resumeFile)
	if err != nil {
		http.Error(w, "Error saving resume file", http.StatusInternalServerError)
		return
	}

	// Assign the resume file name to the contact.Resume field
	contact.Resume = resumeFileName


	func isValidResume(resumeHeader *multipart.FileHeader) bool {
		if resumeHeader != nil {
			if resumeHeader.Size <= 5*1024*1024 && (resumeHeader.Header.Get("Content-Type") == "application/pdf" || resumeHeader.Header.Get("Content-Type") == "image/png") {
				return true
			}
		}
		return false
	}
}