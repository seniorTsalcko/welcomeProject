package auth

import (
	"encoding/json"
	"net/http"
	"welcomeProject/pkg/utils"
)

type AuthHandlers struct {
	service *AuthService
}

func NewAuthHandlers(service *AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

func (h *AuthHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid json format")
		return
	}

	if err := h.service.repo.createUser(&user); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, user)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid json format")
		return
	}

	user, err := h.service.repo.GetUserUserByEmail(input.Email)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := h.service.CheckPassword(user, input.Password); err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := h.service.GenerateToken(user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.JSONResponse(w, http.StatusOK, Token{Token: token})
}
