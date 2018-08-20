// Package controller represents
// the handler which manages the crud
package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/bookfire/model"

	"gopkg.in/mgo.v2"
)

var (
	secretKey = []byte("secret")
)

// UserHandler is in the charge of authentication
// entity
type UserHandler struct {
	db *mgo.Session
}

// NewUserHandler generates
func NewUserHandler(s *mgo.Session) *UserHandler {
	return &UserHandler{s}
}

// Create adds new resources to database
func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dao := model.NewUserDAO(u.db)
	if err := dao.Create(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["sub"] = user.Login
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(5) * 24 * 7).Unix()
	claims["iax"] = time.Now()

	token.Claims = claims
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, tokenString)
}

// Login retrieves the token of the
func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dao := model.NewUserDAO(u.db)
	search, err := dao.Find(user.Login, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["sub"] = search.Login
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(5) * 24 * 7).Unix()
	claims["iax"] = time.Now()

	token.Claims = claims
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, tokenString)
}

// Middleware type represents the alias for  the middleware constructor
type Middleware func(http.HandlerFunc) http.HandlerFunc

// EnabledJwt is a middleware that checks the existent token
func EnabledJwt() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
				func(token *jwt.Token) (interface{}, error) {
					return secretKey, nil
				})

			if err == nil {
				if token.Valid {
					next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprint(w, "Token is not valid")
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Unauthorized access to this resource")
			}
		}
	}
}
