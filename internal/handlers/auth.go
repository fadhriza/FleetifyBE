package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"fleetify/internal/database"
	"fleetify/internal/models"
	"fleetify/pkg/errors"
	"fleetify/pkg/jwt"
	"fleetify/pkg/password"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token     string      `json:"token"`
	User      interface{} `json:"user"`
	ExpiresIn string      `json:"expires_in"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if req.Username == "" || req.Password == "" || req.Role == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username, password, role, and full_name are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Password must be at least 6 characters",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.Users
	checkQuery := `SELECT users_id FROM users WHERE username = $1`
	err := database.DB.QueryRow(ctx, checkQuery, req.Username).Scan(&existingUser.UsersId)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Username already exists",
		})
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		errors.LogError("Password hashing error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process password",
		})
	}

	userID := uuid.New().String()
	now := time.Now()

	insertQuery := `
		INSERT INTO users (users_id, username, password, role, full_name, email, phone, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING users_id, username, role, full_name, email, phone, is_active, created_at, updated_at
	`

	var user models.Users
	err = database.DB.QueryRow(ctx, insertQuery,
		userID,
		req.Username,
		hashedPassword,
		req.Role,
		req.FullName,
		req.Email,
		req.Phone,
		true,
		now,
		now,
	).Scan(
		&user.UsersId,
		&user.Username,
		&user.Role,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		errors.LogError("User registration error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to register user",
		})
	}

	token, err := jwt.GenerateToken(user.UsersId, user.Username, user.Role)
	if err != nil {
		errors.LogError("Token generation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate token",
		})
	}

	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"data": AuthResponse{
			Token:     token,
			User:      user,
			ExpiresIn: "24h",
		},
	})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username and password are required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.Users
	query := `SELECT users_id, username, password, role, full_name, email, phone, is_active FROM users WHERE username = $1`
	err := database.DB.QueryRow(ctx, query, req.Username).Scan(
		&user.UsersId,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.IsActive,
	)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Account is inactive",
		})
	}

	if !password.Verify(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
	}

	token, err := jwt.GenerateToken(user.UsersId, user.Username, user.Role)
	if err != nil {
		errors.LogError("Token generation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate token",
		})
	}

	user.Password = ""

	return c.JSON(fiber.Map{
		"error": false,
		"data": AuthResponse{
			Token:     token,
			User:      user,
			ExpiresIn: "24h",
		},
	})
}

func GetToken(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Claims)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.Users
	query := `SELECT users_id, username, role, full_name, email, phone, is_active FROM users WHERE users_id = $1`
	err := database.DB.QueryRow(ctx, query, claims.UserID).Scan(
		&user.UsersId,
		&user.Username,
		&user.Role,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.IsActive,
	)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Account is inactive",
		})
	}

	token, err := jwt.GenerateToken(user.UsersId, user.Username, user.Role)
	if err != nil {
		errors.LogError("Token generation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data": AuthResponse{
			Token:     token,
			User:      user,
			ExpiresIn: "24h",
		},
	})
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Logged out successfully",
	})
}

