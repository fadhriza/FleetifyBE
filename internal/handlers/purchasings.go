package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"fleetify/internal/database"
	"fleetify/internal/models"
	"fleetify/pkg/errors"
)

type PurchasingDetailRequest struct {
	ItemId int64 `json:"item_id" validate:"required"`
	Qty    int   `json:"qty" validate:"required,gt=0"`
}

type CreatePurchasingRequest struct {
	Date       string                 `json:"date" validate:"required"`
	SupplierId int64                  `json:"supplier_id" validate:"required"`
	UserId     int64                  `json:"user_id" validate:"required"`
	Status     string                 `json:"status"`
	Notes      string                 `json:"notes"`
	Details    []PurchasingDetailRequest `json:"details" validate:"required,min=1"`
}

type UpdatePurchasingRequest struct {
	Date       *string `json:"date"`
	SupplierId *int64  `json:"supplier_id"`
	UserId     *int64  `json:"user_id"`
	Status     *string `json:"status"`
	Notes      *string `json:"notes"`
}

type PurchasingResponse struct {
	models.Purchasings
	SupplierName string `json:"supplier_name"`
	UserName     string `json:"user_name"`
	Details      []models.PurchasingDetails `json:"details"`
}

func GetPurchasings(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT p.id, p.date, p.supplier_id, p.user_id, p.grand_total, p.status, p.notes, p.created_at,
		       s.name as supplier_name, u.full_name as user_name
		FROM purchasings p
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		LEFT JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		errors.LogError("Get purchasings query error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch purchasings",
		})
	}
	defer rows.Close()

	var purchasings []PurchasingResponse
	for rows.Next() {
		var p PurchasingResponse
		var supplierName sql.NullString
		var userName sql.NullString
		err := rows.Scan(
			&p.Id,
			&p.Date,
			&p.SupplierId,
			&p.UserId,
			&p.GrandTotal,
			&p.Status,
			&p.Notes,
			&p.CreatedAt,
			&supplierName,
			&userName,
		)
		if err != nil {
			errors.LogError("Purchasing scan error", err)
			continue
		}
		p.SupplierName = supplierName.String
		p.UserName = userName.String

		detailsQuery := `
			SELECT id, purchasing_id, item_id, qty, subtotal
			FROM purchasing_details
			WHERE purchasing_id = $1
		`
		detailsRows, err := database.DB.Query(ctx, detailsQuery, p.Id)
		if err == nil {
			defer detailsRows.Close()
			for detailsRows.Next() {
				var detail models.PurchasingDetails
				detailsRows.Scan(
					&detail.Id,
					&detail.PurchasingId,
					&detail.ItemId,
					&detail.Qty,
					&detail.Subtotal,
				)
				p.Details = append(p.Details, detail)
			}
		}

		purchasings = append(purchasings, p)
	}

	if err = rows.Err(); err != nil {
		errors.LogError("Rows iteration error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process purchasings",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  purchasings,
		"count": len(purchasings),
	})
}

func GetPurchasingById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p PurchasingResponse
	var supplierName sql.NullString
	var userName sql.NullString
	query := `
		SELECT p.id, p.date, p.supplier_id, p.user_id, p.grand_total, p.status, p.notes, p.created_at,
		       s.name as supplier_name, u.full_name as user_name
		FROM purchasings p
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.id = $1
	`

	err := database.DB.QueryRow(ctx, query, id).Scan(
		&p.Id,
		&p.Date,
		&p.SupplierId,
		&p.UserId,
		&p.GrandTotal,
		&p.Status,
		&p.Notes,
		&p.CreatedAt,
		&supplierName,
		&userName,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing not found",
		})
	}
	p.SupplierName = supplierName.String
	p.UserName = userName.String

	detailsQuery := `
		SELECT id, purchasing_id, item_id, qty, subtotal
		FROM purchasing_details
		WHERE purchasing_id = $1
	`
	detailsRows, err := database.DB.Query(ctx, detailsQuery, id)
	if err == nil {
		defer detailsRows.Close()
		for detailsRows.Next() {
			var detail models.PurchasingDetails
			detailsRows.Scan(
				&detail.Id,
				&detail.PurchasingId,
				&detail.ItemId,
				&detail.Qty,
				&detail.Subtotal,
			)
			p.Details = append(p.Details, detail)
		}
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  p,
	})
}

func CreatePurchasing(c *fiber.Ctx) error {
	var req CreatePurchasingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if len(req.Details) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "At least one detail item is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		errors.LogError("Transaction begin error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to start transaction",
		})
	}
	defer tx.Rollback(ctx)

	purchasingDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	var supplierExists int64
	err = tx.QueryRow(ctx, "SELECT id FROM suppliers WHERE id = $1", req.SupplierId).Scan(&supplierExists)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Supplier not found",
		})
	}

	var userExists int64
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", req.UserId).Scan(&userExists)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	grandTotal := 0.0
	for _, detail := range req.Details {
		var itemPrice float64
		err = tx.QueryRow(ctx, "SELECT price FROM items WHERE id = $1", detail.ItemId).Scan(&itemPrice)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": fmt.Sprintf("Item with ID %d not found", detail.ItemId),
			})
		}
		subtotal := itemPrice * float64(detail.Qty)
		grandTotal += subtotal
	}

	status := "pending"
	if req.Status != "" {
		status = req.Status
	}

	now := time.Now()
	var purchasingId int64
	insertQuery := `
		INSERT INTO purchasings (date, supplier_id, user_id, grand_total, status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err = tx.QueryRow(ctx, insertQuery,
		purchasingDate,
		req.SupplierId,
		req.UserId,
		grandTotal,
		status,
		req.Notes,
		now,
	).Scan(&purchasingId)

	if err != nil {
		errors.LogError("Purchasing creation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create purchasing",
		})
	}

	for _, detail := range req.Details {
		var itemPrice float64
		err = tx.QueryRow(ctx, "SELECT price FROM items WHERE id = $1", detail.ItemId).Scan(&itemPrice)
		if err != nil {
			errors.LogError("Item price query error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": fmt.Sprintf("Failed to get price for item ID %d", detail.ItemId),
			})
		}
		subtotal := itemPrice * float64(detail.Qty)
		detailQuery := `
			INSERT INTO purchasing_details (purchasing_id, item_id, qty, subtotal)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(ctx, detailQuery, purchasingId, detail.ItemId, detail.Qty, subtotal)
		if err != nil {
			errors.LogError("Purchasing detail creation error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to create purchasing details",
			})
		}
	}

	if err = tx.Commit(ctx); err != nil {
		errors.LogError("Transaction commit error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to commit transaction",
		})
	}

	var purchasing models.Purchasings
	getQuery := `
		SELECT id, date, supplier_id, user_id, grand_total, status, notes, created_at
		FROM purchasings
		WHERE id = $1
	`
	database.DB.QueryRow(ctx, getQuery, purchasingId).Scan(
		&purchasing.Id,
		&purchasing.Date,
		&purchasing.SupplierId,
		&purchasing.UserId,
		&purchasing.GrandTotal,
		&purchasing.Status,
		&purchasing.Notes,
		&purchasing.CreatedAt,
	)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Purchasing created successfully",
		"data":    purchasing,
	})
}

func UpdatePurchasing(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing ID is required",
		})
	}

	var req UpdatePurchasingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingPurchasing models.Purchasings
	checkQuery := `SELECT id FROM purchasings WHERE id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingPurchasing.Id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing not found",
		})
	}

	updateFields := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Date != nil {
		purchasingDate, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid date format. Use YYYY-MM-DD",
			})
		}
		updateFields = append(updateFields, fmt.Sprintf("date = $%d", argPos))
		args = append(args, purchasingDate)
		argPos++
	}

	if req.SupplierId != nil {
		var supplierExists int64
		err = database.DB.QueryRow(ctx, "SELECT id FROM suppliers WHERE id = $1", *req.SupplierId).Scan(&supplierExists)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Supplier not found",
			})
		}
		updateFields = append(updateFields, fmt.Sprintf("supplier_id = $%d", argPos))
		args = append(args, *req.SupplierId)
		argPos++
	}

	if req.UserId != nil {
		var userExists int64
		err = database.DB.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", *req.UserId).Scan(&userExists)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "User not found",
			})
		}
		updateFields = append(updateFields, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, *req.UserId)
		argPos++
	}

	if req.Status != nil {
		updateFields = append(updateFields, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *req.Status)
		argPos++
	}

	if req.Notes != nil {
		updateFields = append(updateFields, fmt.Sprintf("notes = $%d", argPos))
		args = append(args, *req.Notes)
		argPos++
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No fields to update",
		})
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE purchasings
		SET %s
		WHERE id = $%d
		RETURNING id, date, supplier_id, user_id, grand_total, status, notes, created_at
	`, strings.Join(updateFields, ", "), argPos)

	var purchasing models.Purchasings
	err = database.DB.QueryRow(ctx, query, args...).Scan(
		&purchasing.Id,
		&purchasing.Date,
		&purchasing.SupplierId,
		&purchasing.UserId,
		&purchasing.GrandTotal,
		&purchasing.Status,
		&purchasing.Notes,
		&purchasing.CreatedAt,
	)

	if err != nil {
		errors.LogError("Purchasing update error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update purchasing",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Purchasing updated successfully",
		"data":    purchasing,
	})
}

func DeletePurchasing(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingPurchasing models.Purchasings
	checkQuery := `SELECT id FROM purchasings WHERE id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingPurchasing.Id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing not found",
		})
	}

	deleteQuery := `DELETE FROM purchasings WHERE id = $1`
	_, err = database.DB.Exec(ctx, deleteQuery, id)
	if err != nil {
		errors.LogError("Purchasing deletion error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete purchasing",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Purchasing deleted successfully",
	})
}

