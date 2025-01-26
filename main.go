package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"
	"fmt"
	"encoding/csv"

	emailVerifier "github.com/AfterShip/email-verifier"
)

var verifier = emailVerifier.NewVerifier()

// EmailVerificationResult holds the verification outcome
type EmailVerificationResult struct {
    Email   string `json:"email"`
    Valid   bool   `json:"valid"`
    Message string `json:"message"`
}

func main() {
    http.HandleFunc("/verify", VerifyEmailHandler)
    http.HandleFunc("/export", ExportCSVHandler)
    http.HandleFunc("/", serveHTML)

    httpPort := os.Getenv("PORT")
    if len(httpPort) == 0 {
        httpPort = "8080"
    }

    log.Println("Starting server at http://localhost:" + httpPort)
    err := http.ListenAndServe(":"+httpPort, logRequest(http.DefaultServeMux))
    if err != nil {
        panic(err)
    }
}

// VerifyEmailHandler handles email verification requests
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("RET: %+v\n", ret)

    result := EmailVerificationResult{
        Email:   email,
        Valid:   ret.Syntax.Valid && ret.HasMxRecords,
        Message: "Invalid Email",
    }
	log.Printf("Email: %s, Valid: %t\n", email, result.Valid)
    if result.Valid {
        result.Message = "Valid Email"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// ExportCSVHandler allows users to download valid emails as a CSV file
func ExportCSVHandler(w http.ResponseWriter, r *http.Request) {
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

func serveHTML(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func logRequest(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
        handler.ServeHTTP(w, r)
    })
}