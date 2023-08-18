package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type ContactMessage struct {
	ID          int
	FullName    string
	Gender      string
	FromDate    string
	ToDate      string
	PhoneNumber string
	Resume      string
	Email       string
	Message     string
}

var db *sql.DB

func main() {
	// Update these with your PostgreSQL connection details
	dbHost := "localhost"
	dbPort := 5432
	dbUser := "postgres"
	dbPassword := "Haripriya@2001"
	dbName := "xenonstack"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY,
			full_name TEXT,
			gender TEXT,
			from_date TEXT,
			to_date TEXT,
			phone_number TEXT,
			resume TEXT,
			email TEXT,
			message TEXT
		);
	`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	http.HandleFunc("/", contactFormHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/messages", getMessagesHandler)

	fmt.Println("Server started on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func contactFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Contact Form</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
			}
			.container {
				width: 50%;
				margin: auto;
				padding: 20px;
				background-color: #fff;
				box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				border-radius: 5px;
			}
			h1 {
				text-align: center;
			}
			form {
				margin-top: 10px;
			}
			label {
				display: block;
				font-weight: bold;
				margin-top: 10px;
			}
			input[type="text"],
			input[type="date"],
			input[type="tel"],
			input[type="email"],
			select,
			textarea {
				width: 100%;
				padding: 10px;
				margin-top: 5px;
				border: 1px solid #ccc;
				border-radius: 5px;
				font-size: 16px;
			}
			select {
				height: 40px;
			}
			textarea {
				height: 100px;
			}
			input[type="submit"] {
				background-color: #007bff;
				color: #fff;
				padding: 10px 20px;
				border: none;
				border-radius: 5px;
				cursor: pointer;
				font-size: 16px;
			}
			
		  </style>

	</head>
<body>
	<div class="container">
		<h1>Contact Us</h1>
		<form action="/submit" method="post" enctype="multipart/form-data">
			
	<label for="full_name">Full Name:</label><br>
	<input type="text" id="full_name" name="full_name" required><br>

	<label for="gender">Gender:</label><br>
	<select id="gender" name="gender" required>
		<option value="Male">Male</option>
		<option value="Female">Female</option>
		<option value="Others">Others</option>
	</select><br>

	<label for="from_date">From Date:</label><br>
	<input type="date" id="from_date" name="from_date" required><br>

	<label for="to_date">To Date:</label><br>
	<input type="date" id="to_date" name="to_date" required><br>

	<label for="phone_number">Phone Number:</label><br>
	<input type="tel" id="phone_number" name="phone_number" pattern="[0-9]+" required><br>

	<label for="resume">Resume (PDF/PNG, max 5MB):</label><br>
	<input type="file" id="resume" name="resume" accept=".pdf,.png" required><br>

	<label for="email">Email:</label><br>
	<input type="email" id="email" name="email" required><br>
	<label for="message">Message:</label><br>
			<textarea id="message" name="message" rows="4" required></textarea><br>

			<input type="submit" value="Submit">
		</form>
	</div>
</body>
</html>
	`

	tmplParsed := template.Must(template.New("contactForm").Parse(tmpl))
	tmplParsed.Execute(w, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(5 * 1024 * 1024) // Max file size 5MB
	if err != nil {
		fmt.Println("Error parsing form data:", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	contact := ContactMessage{
		FullName:    strings.TrimSpace(r.FormValue("full_name")),
		Gender:      r.FormValue("gender"),
		FromDate:    r.FormValue("from_date"),
		ToDate:      r.FormValue("to_date"),
		PhoneNumber: r.FormValue("phone_number"),
		Email:       r.FormValue("email"),
		Message:     r.FormValue("message"),
	}

	if !isValidPhoneNumber(contact.PhoneNumber) {
		http.Error(w, "Please enter a valid Indian phone number", http.StatusBadRequest)
		return
	}

	if !isValidEmail(contact.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Resume validation
	_, resumeHeader, err := r.FormFile("resume")
	if err == nil {
		if !isValidResume(resumeHeader) {
			http.Error(w, "Invalid resume file. Only PDF and PNG files up to 5MB are allowed.", http.StatusBadRequest)
			return
		}
		// Save the resume if needed
		// Example: contact.Resume = saveFileAndGetPath(file)
	}

	fromDateParsed, _ := time.Parse("2006-01-02", contact.FromDate)
	toDateParsed, _ := time.Parse("2006-01-02", contact.ToDate)
	if toDateParsed.Before(fromDateParsed) {
		http.Error(w, "To Date should be greater than From Date", http.StatusBadRequest)
		return
	}

	if len(contact.FullName) >= 30 {
		http.Error(w, "Full Name should be less than 30 characters", http.StatusBadRequest)
		return
	}

	sqlStatement := "INSERT INTO messages (full_name, gender, from_date, to_date, phone_number, resume, email, message) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err = db.Exec(sqlStatement, contact.FullName, contact.Gender, contact.FromDate, contact.ToDate, contact.PhoneNumber, contact.Resume, contact.Email, contact.Message)
	if err != nil {
		fmt.Println("Error storing data in the database:", err)
		http.Error(w, "Error storing data in the database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message submitted successfully!")
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, full_name, gender, from_date, to_date, phone_number, email, message FROM messages")
	if err != nil {
		http.Error(w, "Error retrieving data from the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []ContactMessage
	for rows.Next() {
		var msg ContactMessage
		err := rows.Scan(&msg.ID, &msg.FullName, &msg.Gender, &msg.FromDate, &msg.ToDate, &msg.PhoneNumber, &msg.Email, &msg.Message)
		if err != nil {
			http.Error(w, "Error scanning database rows", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	fmt.Fprintf(w, "<h1>Contact Messages</h1>")
	for _, msg := range messages {
		fmt.Fprintf(w, "<p><strong>Full Name:</strong> %s<br><strong>Gender:</strong> %s<br><strong>From Date:</strong> %s<br><strong>To Date:</strong> %s<br><strong>Phone Number:</strong> %s<br><strong>Email:</strong> %s<br><strong>Message:</strong> %s</p><hr>", msg.FullName, msg.Gender, msg.FromDate, msg.ToDate, msg.PhoneNumber, msg.Email, msg.Message)
	}
}

func isValidPhoneNumber(phoneNumber string) bool {
	match, err := regexp.MatchString(`^\d{10}$`, phoneNumber)
	return err == nil && match
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func isValidResume(resumeHeader *multipart.FileHeader) bool {
	if resumeHeader != nil {
		if resumeHeader.Size <= 5*1024*1024 && (resumeHeader.Header.Get("Content-Type") == "application/pdf" || resumeHeader.Header.Get("Content-Type") == "image/png") {
			return true
		}
	}
	return false
}
