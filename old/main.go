package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"
)

func main() {
    http.HandleFunc("/verify", handleVerifyRequest)
    http.HandleFunc("/export", handleExportRequest)
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

func handleVerifyRequest(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "You need to pass an email address to verify.", 500)
        return
    }
	
    verifyResult := VerifyResult{Email: email}
    verifyResult.Verify()

    json.NewEncoder(w).Encode(verifyResult)
}

func handleExportRequest(w http.ResponseWriter, r *http.Request) {
    emails := r.URL.Query().Get("emails")
    if emails == "" {
        http.Error(w, "No emails to export.", 400)
        return
    }

    w.Header().Set("Content-Disposition", "attachment; filename=valid_emails.txt")
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(strings.Replace(emails, ",", "\n", -1)))
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