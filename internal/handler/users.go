package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qm3llz/tasksWebApi/internal/auth"
	"github.com/qm3llz/tasksWebApi/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(ctx context.Context, username string, hash []byte) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type UserHandler struct {
	repo UserRepo
}

func NewUserHandler(repo UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "status internal server error", http.StatusInternalServerError)
		return
	}

	err = u.repo.Create(r.Context(), user.Username, hash)
	if err != nil {
		http.Error(w, "status internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	myUser, err := u.repo.GetByUsername(r.Context(), user.Username)
	if err != nil {
		http.Error(w, "status internal server error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(myUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "username or password incorrect", http.StatusUnauthorized)
		return
	}
	
	token, err := auth.GenerateToken(myUser.ID)
	if err != nil {
		http.Error(w, "status internal server error", http.StatusInternalServerError)
		return
	}
	
	fmt.Fprint(w, token)
}
