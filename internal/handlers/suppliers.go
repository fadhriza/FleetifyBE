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
	"fleetify/pkg/query"
)

type CreateSupplierRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	SupplierType string `json:"supplier_type"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateSupplierRequest struct {
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	Address      *string `json:"address"`
	Phone        *string `json:"phone"`
	SupplierType *string `json:"supplier_type"`
	IsActive     *bool   `json:"is_active"`
}

func GetSuppliers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := query.ParseQueryParams(c)
	
	searchFields := []string{"name", "email", "address", "phone", "supplier_type"}
	filterFields := map[string]string{
		"supplier_type": "supplier_type",
		"is_active":     "is_active",
	}

	whereClause, whereArgs := query.BuildWhereClause(params, searchFields, filterFields)
	orderClause := query.BuildOrderClause(params, "created_at")
	paginationClause, paginationArgs := query.BuildPaginationClause(params, len(whereArgs)+1)

	countQuery := query.BuildCountQuery("suppliers", whereClause)

	var totalCount int
	err := database.DB.QueryRow(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		errors.LogError("Get suppliers count error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to count suppliers",
		})
	}

	baseQuery := `
		SELECT suppliers_id, name, email, address, phone, supplier_type, is_active, created_at, updated_at
		FROM suppliers
	`
	
	fullQuery := baseQuery + " " + whereClause + " " + orderClause + " " + paginationClause
	allArgs := append(whereArgs, paginationArgs...)

	rows, err := database.DB.Query(ctx, fullQuery, allArgs...)
	if err != nil {
		errors.LogError("Get suppliers query error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch suppliers",
		})
	}
	defer rows.Close()

	var suppliers []models.Suppliers
	for rows.Next() {
		var supplier models.Suppliers
		err := rows.Scan(
			&supplier.SuppliersId,
			&supplier.Name,
			&supplier.Email,
			&supplier.Address,
			&supplier.Phone,
			&supplier.SupplierType,
			&supplier.IsActive,
			&supplier.CreatedAt,
			&supplier.UpdatedAt,
		)
		if err != nil {
			errors.LogError("Supplier scan error", err)
			continue
		}
		suppliers = append(suppliers, supplier)
	}

	if err = rows.Err(); err != nil {
		errors.LogError("Rows iteration error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process suppliers",
		})
	}

	response := query.NewPaginatedResponse(suppliers, totalCount, params.Page, params.Limit)
	return c.JSON(fiber.Map{
		"error": false,
		"data":  response.Data,
		"count": response.Count,
		"page":  response.Page,
		"limit": response.Limit,
		"total_pages": response.TotalPages,
		"total_count": response.TotalCount,
	})
}

func GetSupplierById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var supplier models.Suppliers
	query := `
		SELECT suppliers_id, name, email, address, phone, supplier_type, is_active, created_at, updated_at
		FROM suppliers
		WHERE suppliers_id = $1
	`

	err := database.DB.QueryRow(ctx, query, id).Scan(
		&supplier.SuppliersId,
		&supplier.Name,
		&supplier.Email,
		&supplier.Address,
		&supplier.Phone,
		&supplier.SupplierType,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  supplier,
	})
}

func CreateSupplier(c *fiber.Ctx) error {
	var req CreateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier name is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	query := `
		INSERT INTO suppliers (name, email, address, phone, supplier_type, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING suppliers_id, name, email, address, phone, supplier_type, is_active, created_at, updated_at
	`

	var supplier models.Suppliers
	err := database.DB.QueryRow(ctx, query,
		req.Name,
		req.Email,
		req.Address,
		req.Phone,
		req.SupplierType,
		isActive,
		now,
		now,
	).Scan(
		&supplier.SuppliersId,
		&supplier.Name,
		&supplier.Email,
		&supplier.Address,
		&supplier.Phone,
		&supplier.SupplierType,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		errors.LogError("Supplier creation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create supplier",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Supplier created successfully",
		"data":    supplier,
	})
}

func UpdateSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier ID is required",
		})
	}

	var req UpdateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingSupplier models.Suppliers
	checkQuery := `SELECT suppliers_id FROM suppliers WHERE suppliers_id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingSupplier.SuppliersId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier not found",
		})
	}

	updateFields := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Name != nil {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}

	if req.Email != nil {
		updateFields = append(updateFields, fmt.Sprintf("email = $%d", argPos))
		args = append(args, *req.Email)
		argPos++
	}

	if req.Address != nil {
		updateFields = append(updateFields, fmt.Sprintf("address = $%d", argPos))
		args = append(args, *req.Address)
		argPos++
	}

	if req.Phone != nil {
		updateFields = append(updateFields, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *req.Phone)
		argPos++
	}

	if req.SupplierType != nil {
		updateFields = append(updateFields, fmt.Sprintf("supplier_type = $%d", argPos))
		args = append(args, *req.SupplierType)
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

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE suppliers
		SET %s
		WHERE suppliers_id = $%d
		RETURNING suppliers_id, name, email, address, phone, supplier_type, is_active, created_at, updated_at
	`, strings.Join(updateFields, ", "), argPos)

	var supplier models.Suppliers
	err = database.DB.QueryRow(ctx, query, args...).Scan(
		&supplier.SuppliersId,
		&supplier.Name,
		&supplier.Email,
		&supplier.Address,
		&supplier.Phone,
		&supplier.SupplierType,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		errors.LogError("Supplier update error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update supplier",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Supplier updated successfully",
		"data":    supplier,
	})
}

func DeleteSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingSupplier models.Suppliers
	checkQuery := `SELECT suppliers_id FROM suppliers WHERE suppliers_id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingSupplier.SuppliersId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier not found",
		})
	}

	deleteQuery := `DELETE FROM suppliers WHERE suppliers_id = $1`
	_, err = database.DB.Exec(ctx, deleteQuery, id)
	if err != nil {
		errors.LogError("Supplier deletion error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete supplier",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Supplier deleted successfully",
	})
}

