package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"fleetify/internal/database"
	"fleetify/internal/models"
	"fleetify/pkg/errors"
	"fleetify/pkg/password"
)

type UpdateUserRequest struct {
	Role     string `json:"role"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	IsActive *bool  `json:"is_active"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT users_id, username, role, full_name, email, phone, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		errors.LogError("Get users query error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch users",
		})
	}
	defer rows.Close()

	var users []models.Users
	for rows.Next() {
		var user models.Users
		err := rows.Scan(
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
			errors.LogError("User scan error", err)
			continue
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		errors.LogError("Rows iteration error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process users",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  users,
		"count": len(users),
	})
}

func GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("uname")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username parameter is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.Users
	query := `
		SELECT users_id, username, role, full_name, email, phone, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := database.DB.QueryRow(ctx, query, username).Scan(
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  user,
	})
}

func UpdateUserByUsername(c *fiber.Ctx) error {
	username := c.Params("uname")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username parameter is required",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.Users
	checkQuery := `SELECT users_id FROM users WHERE username = $1`
	err := database.DB.QueryRow(ctx, checkQuery, username).Scan(&existingUser.UsersId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	updateFields := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Role != "" {
		updateFields = append(updateFields, fmt.Sprintf("role = $%d", argPos))
		args = append(args, req.Role)
		argPos++
	}

	if req.FullName != "" {
		updateFields = append(updateFields, fmt.Sprintf("full_name = $%d", argPos))
		args = append(args, req.FullName)
		argPos++
	}

	if req.Email != "" {
		updateFields = append(updateFields, fmt.Sprintf("email = $%d", argPos))
		args = append(args, req.Email)
		argPos++
	}

	if req.Phone != "" {
		updateFields = append(updateFields, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, req.Phone)
		argPos++
	}

	if req.IsActive != nil {
		updateFields = append(updateFields, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *req.IsActive)
		argPos++
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No fields to update",
		})
	}

	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	args = append(args, username)

	updateQuery := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE username = $%d
		RETURNING users_id, username, role, full_name, email, phone, is_active, created_at, updated_at
	`, strings.Join(updateFields, ", "), argPos)

	var user models.Users
	err = database.DB.QueryRow(ctx, updateQuery, args...).Scan(
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
		errors.LogError("User update error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "User updated successfully",
		"data":    user,
	})
}

func ChangePassword(c *fiber.Ctx) error {
	username := c.Params("uname")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username parameter is required",
		})
	}

	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "New password must be at least 6 characters",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.Users
	checkQuery := `SELECT users_id FROM users WHERE username = $1`
	err := database.DB.QueryRow(ctx, checkQuery, username).Scan(&existingUser.UsersId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	hashedPassword, err := password.Hash(req.NewPassword)
	if err != nil {
		errors.LogError("Password hashing error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process password",
		})
	}

	updateQuery := `
		UPDATE users
		SET password = $1, updated_at = $2
		WHERE username = $3
	`

	_, err = database.DB.Exec(ctx, updateQuery, hashedPassword, time.Now(), username)
	if err != nil {
		errors.LogError("Password update error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Password updated successfully",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	username := c.Params("uname")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username parameter is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.Users
	checkQuery := `SELECT users_id FROM users WHERE username = $1`
	err := database.DB.QueryRow(ctx, checkQuery, username).Scan(&existingUser.UsersId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	deleteQuery := `DELETE FROM users WHERE username = $1`
	_, err = database.DB.Exec(ctx, deleteQuery, username)
	if err != nil {
		errors.LogError("User deletion error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "User deleted successfully",
	})
}


