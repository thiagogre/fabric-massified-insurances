package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id    string
	Token string
}

const (
	dbPath = "../organizations/fabric-ca/org1/fabric-ca-server.db"
)

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func GetUserById(db *sql.DB, id string) (*User, error) {
	query := "SELECT id, token FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.Id, &user.Token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Auth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[AUTH] Received a request")

	var authRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	user, err := GetUserById(db, authRequest.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !CheckPasswordHash(authRequest.Password, user.Token) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
	}

	jsonResponse, err := json.Marshal(authRequest)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
