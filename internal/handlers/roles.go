package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"fleetify/internal/database"
	"fleetify/internal/models"
	"fleetify/pkg/errors"
)

type CreateRoleRequest struct {
	RoleOID         string `json:"role_oid" validate:"required"`
	RoleName        string `json:"role_name" validate:"required"`
	RoleDescription string `json:"role_description"`
}

type UpdateRoleRequest struct {
	RoleName        string `json:"role_name"`
	RoleDescription string `json:"role_description"`
}

func GetRoles(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT roles_id, role_oid, role_name, role_description, created_timestamp, updated_timestamp
		FROM roles
		ORDER BY created_timestamp DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		errors.LogError("Get roles query error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch roles",
		})
	}
	defer rows.Close()

	var roles []models.Roles
	for rows.Next() {
		var role models.Roles
		err := rows.Scan(
			&role.RolesId,
			&role.RoleOID,
			&role.RoleName,
			&role.RoleDescription,
			&role.CreatedTimestamp,
			&role.UpdatedTimestamp,
		)
		if err != nil {
			errors.LogError("Role scan error", err)
			continue
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		errors.LogError("Rows iteration error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process roles",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  roles,
		"count": len(roles),
	})
}

func GetRoleByOID(c *fiber.Ctx) error {
	roleOID := c.Params("oid")
	if roleOID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Role OID parameter is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var role models.Roles
	query := `
		SELECT roles_id, role_oid, role_name, role_description, created_timestamp, updated_timestamp
		FROM roles
		WHERE role_oid = $1
	`

	err := database.DB.QueryRow(ctx, query, roleOID).Scan(
		&role.RolesId,
		&role.RoleOID,
		&role.RoleName,
		&role.RoleDescription,
		&role.CreatedTimestamp,
		&role.UpdatedTimestamp,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Role not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  role,
	})
}

func CreateRole(c *fiber.Ctx) error {
	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if req.RoleOID == "" || req.RoleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Role OID and role name are required",
		})
	}

	if strings.ToUpper(req.RoleOID) == "ADMIN" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Cannot create ADMIN role",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingRole models.Roles
	checkQuery := `SELECT roles_id FROM roles WHERE role_oid = $1`
	err := database.DB.QueryRow(ctx, checkQuery, strings.ToUpper(req.RoleOID)).Scan(&existingRole.RolesId)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Role OID already exists",
		})
	}

	roleID := uuid.New().String()
	now := time.Now()

	insertQuery := `
		INSERT INTO roles (roles_id, role_oid, role_name, role_description, created_timestamp, updated_timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING roles_id, role_oid, role_name, role_description, created_timestamp, updated_timestamp
	`

	var role models.Roles
	err = database.DB.QueryRow(ctx, insertQuery,
		roleID,
		strings.ToUpper(req.RoleOID),
		req.RoleName,
		req.RoleDescription,
		now,
		now,
	).Scan(
		&role.RolesId,
		&role.RoleOID,
		&role.RoleName,
		&role.RoleDescription,
		&role.CreatedTimestamp,
		&role.UpdatedTimestamp,
	)

	if err != nil {
		errors.LogError("Role creation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create role",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Role created successfully",
		"data":    role,
	})
}

func UpdateRole(c *fiber.Ctx) error {
	roleOID := c.Params("oid")
	if roleOID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Role OID parameter is required",
		})
	}

	if strings.ToUpper(roleOID) == "ADMIN" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Cannot update ADMIN role",
		})
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingRole models.Roles
	checkQuery := `SELECT roles_id FROM roles WHERE role_oid = $1`
	err := database.DB.QueryRow(ctx, checkQuery, strings.ToUpper(roleOID)).Scan(&existingRole.RolesId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Role not found",
		})
	}

	updateFields := []string{}
	args := []interface{}{}
	argPos := 1

	if req.RoleName != "" {
		updateFields = append(updateFields, fmt.Sprintf("role_name = $%d", argPos))
		args = append(args, req.RoleName)
		argPos++
	}

	if req.RoleDescription != "" {
		updateFields = append(updateFields, fmt.Sprintf("role_description = $%d", argPos))
		args = append(args, req.RoleDescription)
		argPos++
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No fields to update",
		})
	}

	updateFields = append(updateFields, fmt.Sprintf("updated_timestamp = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	args = append(args, strings.ToUpper(roleOID))

	updateQuery := fmt.Sprintf(`
		UPDATE roles
		SET %s
		WHERE role_oid = $%d
		RETURNING roles_id, role_oid, role_name, role_description, created_timestamp, updated_timestamp
	`, strings.Join(updateFields, ", "), argPos)

	var role models.Roles
	err = database.DB.QueryRow(ctx, updateQuery, args...).Scan(
		&role.RolesId,
		&role.RoleOID,
		&role.RoleName,
		&role.RoleDescription,
		&role.CreatedTimestamp,
		&role.UpdatedTimestamp,
	)

	if err != nil {
		errors.LogError("Role update error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update role",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Role updated successfully",
		"data":    role,
	})
}

func DeleteRole(c *fiber.Ctx) error {
	roleOID := c.Params("oid")
	if roleOID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Role OID parameter is required",
		})
	}

	if strings.ToUpper(roleOID) == "ADMIN" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Cannot delete ADMIN role",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingRole models.Roles
	checkQuery := `SELECT roles_id FROM roles WHERE role_oid = $1`
	err := database.DB.QueryRow(ctx, checkQuery, strings.ToUpper(roleOID)).Scan(&existingRole.RolesId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Role not found",
		})
	}

	deleteQuery := `DELETE FROM roles WHERE role_oid = $1`
	_, err = database.DB.Exec(ctx, deleteQuery, strings.ToUpper(roleOID))
	if err != nil {
		errors.LogError("Role deletion error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete role",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Role deleted successfully",
	})
}

