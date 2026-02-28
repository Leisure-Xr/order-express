package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"order-express/backend/internal/config"
	"order-express/backend/internal/models"
	jwtpkg "order-express/backend/pkg/jwt"

	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	creds.Username = strings.TrimSpace(creds.Username)
	if creds.Username == "" || creds.Password == "" {
		BadRequest(w, "username and password are required")
		return
	}

	var user models.User
	tx := h.db.Raw(
		`SELECT id, username, name, password, role, COALESCE(phone,'') AS phone, COALESCE(avatar,'') AS avatar FROM users WHERE username = ?`,
		creds.Username,
	).Scan(&user)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		Fail(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		Fail(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := jwtpkg.GenerateToken(h.cfg.JWTSecret, h.cfg.JWTExpiryHours, user.ID, user.Username, user.Role)
	if err != nil {
		Fail(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	OK(w, map[string]any{
		"token": token,
		"user": models.UserPublic{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Role:     user.Role,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT is stateless; client should discard the token
	OK(w, nil)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		Unauthorized(w)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwtpkg.ParseToken(h.cfg.JWTSecret, tokenStr)
	if err != nil {
		Unauthorized(w)
		return
	}

	var user models.UserPublic
	tx := h.db.Raw(
		`SELECT id, username, name, role, COALESCE(phone,'') AS phone, COALESCE(avatar,'') AS avatar FROM users WHERE id = ?`,
		claims.UserID,
	).Scan(&user)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		Unauthorized(w)
		return
	}

	OK(w, user)
}
