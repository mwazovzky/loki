package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mwazovzky/loki/models"
	"gorm.io/gorm"
)

type AuthHandlers struct {
	users *models.Users
}

func NewAuthHandlers(db *gorm.DB) *AuthHandlers {
	users := models.NewUsers(db)
	return &AuthHandlers{users}
}

type UserLogin struct {
	Email    string
	Password string
}

type Token struct {
	Value string
}

// curl -X POST localhost:3000/api/login -d '{"name":"test","password":"secret"}'
func (ah *AuthHandlers) Login(rw http.ResponseWriter, r *http.Request) {
	ul := UserLogin{}
	err := json.NewDecoder(r.Body).Decode(&ul)
	if err != nil {
		http.Error(rw, "Unable to read request data", http.StatusBadRequest)
		return
	}

	var user models.User
	err = ah.users.FindByEmail(&user, ul.Email)
	if err != nil {
		http.Error(rw, "Wrong email or password", http.StatusBadRequest)
		return
	}

	err = verifyPassword(ul.Password, user.Password)
	if err != nil {
		http.Error(rw, "Wrong email or password", http.StatusBadRequest)
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(rw, "Sometning went wrong", http.StatusInternalServerError)
		return
	}

	user.Token = token
	err = ah.users.Update(&user, user.ID)
	if err != nil {
		http.Error(rw, "Model not found", http.StatusNotFound)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	t := Token{token}
	e := json.NewEncoder(rw)
	err = e.Encode(t)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func verifyPassword(password string, userPassword string) error {
	if password != userPassword {
		return errors.New("wrong password")
	}

	return nil
}

func generateToken() (string, error) {
	len := 10

	b := make([]byte, len)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
