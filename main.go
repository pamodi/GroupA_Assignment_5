package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Function created by Balangoda Pamodi : 500229522
const (
	host     = "localhost"
	port     = 5432
	user     = ""
	password = ""
	dbname   = "user_management_db"
)

// JWT Key
var jwtKey = []byte("A7F077BEE149BB313832C810E5547FB647102356F375795FFB694CFB1A43E37B")

// Struct to represent claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Struct to represent token
type Token struct {
	Token    string
	ExpireAt time.Time
}

func main() {
	// Setup database
	db := SetupDatabase()
	defer db.Close()

	// Define API endpoints
	http.HandleFunc("/token", GenerateTokenHandler())

	fmt.Println("Server started on :8012")
	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8012", nil))
}

// Setup database
func SetupDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")
	return db
}

// Generate token handler
func GenerateTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user email from request header
		formValue := r.PostFormValue("Email")
		if formValue == "" {
			http.Error(w, "Email header not found", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		tokenString, err := GenerateJWTToken(formValue)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Respond with JWT token
		response := Token{Token: tokenString, ExpireAt: time.Now().Add(2 * time.Hour)}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Generate JWT token
func GenerateJWTToken(email string) (string, error) {
	// Create JWT token claims
	claims := jwt.MapClaims{
		"email":  email,
		"expiry": time.Now().Add(2 * time.Hour).Unix(), // Token expires in 2 minutes
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

// Function created by Tejaswi Cheripally: 500229934

// Authentication middleware
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		email := r.FormValue("email")
		claims := jwt.MapClaims{
			"email":  email,
			"expiry": time.Now().Add(2 * time.Hour).Unix(),
		}

		// Extract JWT token from Authorization header
		tokenString := ExtractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Authorization header not found", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if "expiry" claim exists and is valid
		expClaim, ok := claims["expiry"].(float64)
		if !ok || expClaim == 0 {
			http.Error(w, "Expired token", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
//Function created by Syed Abdul Qadeer: 500228186

// Function to extract token from header
func ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return ""
	}
	return tokenParts[1]
}
//Function created by Samhita Dubbaka: 500225971

// Struct to represent invitation code
type InvitationCode struct {
	Code string json:"code"
	Used bool   json:"used"
}

// Generate invitation code
func GenerateInvitationCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")

		// Generate random invitation code
		codeBytes := make([]byte, 16)
		_, err := rand.Read(codeBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		code := base64.URLEncoding.EncodeToString(codeBytes)
		fmt.Println(code)
		// Insert the code into the database
		_, err = db.Exec("INSERT INTO invitation_codes (code, email, used) VALUES ($1, $2, false)", code, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the generated code
		invitationCode := InvitationCode{Code: code}
		json.NewEncoder(w).Encode(invitationCode)
		w.WriteHeader(http.StatusCreated)
	}
}