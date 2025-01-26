package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	emailVerifier "github.com/AfterShip/email-verifier"
)

var verifier = emailVerifier.NewVerifier()

// EmailVerificationResult holds the verification outcome
type EmailVerificationResult struct {
	Email   string `json:"email"`
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// VerifyEmailHandler handles email verification requests
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing email parameter", http.StatusBadRequest)
		return
	}

	ret, err := verifier.Verify(email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verification error: %s", err), http.StatusInternalServerError)
		return
	}

	if !ret.Syntax.Valid {
		http.Error(w, "Email address syntax is invalid", http.StatusBadRequest)
		return
	}


	result := EmailVerificationResult{
		Email:   email,
		Valid:   ret.Syntax.Valid,
		Message: "Invalid Email",
	}
	if result.Valid {
		result.Message = "Valid Email"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ExportCSVHandler allows users to download valid emails as a CSV file
func ExportCSVHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	emails := r.URL.Query().Get("emails")
	if emails == "" {
		http.Error(w, "No emails provided", http.StatusBadRequest)
		return
	}

	validEmails := strings.Split(emails, ",")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=valid_emails.csv")

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	_ = csvWriter.Write([]string{"Valid Emails"})
	for _, email := range validEmails {
		_ = csvWriter.Write([]string{email})
	}
}

func main() {
	router := httprouter.New()
	router.GET("/verify", VerifyEmailHandler)
	router.GET("/export", ExportCSVHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
