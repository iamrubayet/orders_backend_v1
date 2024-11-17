package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang-orders-app/repository"
	"golang-orders-app/utils"
)

// LoginRequest represents the request body for the login endpoint
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for the login endpoint
type LoginResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserHandler represents the handler for user endpoints
type UserHandler struct {
	UserRepo *repository.UserRepository
}

// NewUserHandler initializes and returns a new UserHandler
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// LoginHandler processes login requests
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	user, err := h.UserRepo.GetUserByUsername(loginReq.Username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil || user.Password != loginReq.Password {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "The user credentials were incorrect.",
			"type":    "error",
			"code":    400,
		})
		return
	}

	// Generate JWT Token
	accessToken, err := utils.GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Create response
	response := LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    432000, // 5 days in seconds
		AccessToken:  accessToken,
		RefreshToken: "REFRESH_TOKEN", // Static for now
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler processes logout requests
func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}

	// Extract token from the 'Bearer' prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Validate the token
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Invalidate the token (this would usually involve blacklisting the token)
	// Since JWT is stateless, the client needs to remove the token manually. Here, we just acknowledge the logout.

	// Respond with success message
	response := map[string]interface{}{
		"message": "Successfully logged out",
		"type":    "success",
		"code":    200,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
