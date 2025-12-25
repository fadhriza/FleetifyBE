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
)

type CreatePurchasingDetailRequest struct {
	PurchasingId int64   `json:"purchasing_id" validate:"required"`
	ItemId       int64   `json:"item_id" validate:"required"`
	Qty          int     `json:"qty" validate:"required,gt=0"`
	Subtotal     float64 `json:"subtotal" validate:"required,gt=0"`
}

type UpdatePurchasingDetailRequest struct {
	ItemId   *int64   `json:"item_id"`
	Qty      *int     `json:"qty"`
	Subtotal *float64 `json:"subtotal"`
}

func GetPurchasingDetails(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, purchasing_id, item_id, qty, subtotal
		FROM purchasing_details
		ORDER BY id DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		errors.LogError("Get purchasing details query error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch purchasing details",
		})
	}
	defer rows.Close()

	var details []models.PurchasingDetails
	for rows.Next() {
		var detail models.PurchasingDetails
		err := rows.Scan(
			&detail.Id,
			&detail.PurchasingId,
			&detail.ItemId,
			&detail.Qty,
			&detail.Subtotal,
		)
		if err != nil {
			errors.LogError("Purchasing detail scan error", err)
			continue
		}
		details = append(details, detail)
	}

	if err = rows.Err(); err != nil {
		errors.LogError("Rows iteration error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process purchasing details",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  details,
		"count": len(details),
	})
}

func GetPurchasingDetailsByPurchasingId(c *fiber.Ctx) error {
	purchasingId := c.Params("purchasing_id")
	if purchasingId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, purchasing_id, item_id, qty, subtotal
		FROM purchasing_details
		WHERE purchasing_id = $1
		ORDER BY id
	`

	rows, err := database.DB.Query(ctx, query, purchasingId)
	if err != nil {
		errors.LogError("Get purchasing details query error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch purchasing details",
		})
	}
	defer rows.Close()

	var details []models.PurchasingDetails
	for rows.Next() {
		var detail models.PurchasingDetails
		err := rows.Scan(
			&detail.Id,
			&detail.PurchasingId,
			&detail.ItemId,
			&detail.Qty,
			&detail.Subtotal,
		)
		if err != nil {
			errors.LogError("Purchasing detail scan error", err)
			continue
		}
		details = append(details, detail)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  details,
		"count": len(details),
	})
}

func GetPurchasingDetailById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing detail ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var detail models.PurchasingDetails
	query := `
		SELECT id, purchasing_id, item_id, qty, subtotal
		FROM purchasing_details
		WHERE id = $1
	`

	err := database.DB.QueryRow(ctx, query, id).Scan(
		&detail.Id,
		&detail.PurchasingId,
		&detail.ItemId,
		&detail.Qty,
		&detail.Subtotal,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing detail not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  detail,
	})
}

func CreatePurchasingDetail(c *fiber.Ctx) error {
	var req CreatePurchasingDetailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var purchasingExists int64
	err := database.DB.QueryRow(ctx, "SELECT id FROM purchasings WHERE id = $1", req.PurchasingId).Scan(&purchasingExists)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing not found",
		})
	}

	var itemExists int64
	err = database.DB.QueryRow(ctx, "SELECT id FROM items WHERE id = $1", req.ItemId).Scan(&itemExists)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Item not found",
		})
	}

	query := `
		INSERT INTO purchasing_details (purchasing_id, item_id, qty, subtotal)
		VALUES ($1, $2, $3, $4)
		RETURNING id, purchasing_id, item_id, qty, subtotal
	`

	var detail models.PurchasingDetails
	err = database.DB.QueryRow(ctx, query,
		req.PurchasingId,
		req.ItemId,
		req.Qty,
		req.Subtotal,
	).Scan(
		&detail.Id,
		&detail.PurchasingId,
		&detail.ItemId,
		&detail.Qty,
		&detail.Subtotal,
	)

	if err != nil {
		errors.LogError("Purchasing detail creation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create purchasing detail",
		})
	}

	updateGrandTotalQuery := `
		UPDATE purchasings
		SET grand_total = (
			SELECT COALESCE(SUM(subtotal), 0)
			FROM purchasing_details
			WHERE purchasing_id = $1
		)
		WHERE id = $1
	`
	_, err = database.DB.Exec(ctx, updateGrandTotalQuery, req.PurchasingId)
	if err != nil {
		errors.LogError("Grand total update error", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Purchasing detail created successfully",
		"data":    detail,
	})
}

func UpdatePurchasingDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing detail ID is required",
		})
	}

	var req UpdatePurchasingDetailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingDetail models.PurchasingDetails
	checkQuery := `SELECT id, purchasing_id FROM purchasing_details WHERE id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingDetail.Id, &existingDetail.PurchasingId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing detail not found",
		})
	}

	updateFields := []string{}
	args := []interface{}{}
	argPos := 1

	if req.ItemId != nil {
		var itemExists int64
		err = database.DB.QueryRow(ctx, "SELECT id FROM items WHERE id = $1", *req.ItemId).Scan(&itemExists)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Item not found",
			})
		}
		updateFields = append(updateFields, fmt.Sprintf("item_id = $%d", argPos))
		args = append(args, *req.ItemId)
		argPos++
	}

	if req.Qty != nil {
		if *req.Qty <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Quantity must be greater than 0",
			})
		}
		updateFields = append(updateFields, fmt.Sprintf("qty = $%d", argPos))
		args = append(args, *req.Qty)
		argPos++
	}

	if req.Subtotal != nil {
		if *req.Subtotal <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Subtotal must be greater than 0",
			})
		}
		updateFields = append(updateFields, fmt.Sprintf("subtotal = $%d", argPos))
		args = append(args, *req.Subtotal)
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
		UPDATE purchasing_details
		SET %s
		WHERE id = $%d
		RETURNING id, purchasing_id, item_id, qty, subtotal
	`, strings.Join(updateFields, ", "), argPos)

	var detail models.PurchasingDetails
	err = database.DB.QueryRow(ctx, query, args...).Scan(
		&detail.Id,
		&detail.PurchasingId,
		&detail.ItemId,
		&detail.Qty,
		&detail.Subtotal,
	)

	if err != nil {
		errors.LogError("Purchasing detail update error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update purchasing detail",
		})
	}

	updateGrandTotalQuery := `
		UPDATE purchasings
		SET grand_total = (
			SELECT COALESCE(SUM(subtotal), 0)
			FROM purchasing_details
			WHERE purchasing_id = $1
		)
		WHERE id = $1
	`
	_, err = database.DB.Exec(ctx, updateGrandTotalQuery, detail.PurchasingId)
	if err != nil {
		errors.LogError("Grand total update error", err)
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Purchasing detail updated successfully",
		"data":    detail,
	})
}

func DeletePurchasingDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing detail ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingDetail models.PurchasingDetails
	checkQuery := `SELECT id, purchasing_id FROM purchasing_details WHERE id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingDetail.Id, &existingDetail.PurchasingId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Purchasing detail not found",
		})
	}

	deleteQuery := `DELETE FROM purchasing_details WHERE id = $1`
	_, err = database.DB.Exec(ctx, deleteQuery, id)
	if err != nil {
		errors.LogError("Purchasing detail deletion error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete purchasing detail",
		})
	}

	updateGrandTotalQuery := `
		UPDATE purchasings
		SET grand_total = (
			SELECT COALESCE(SUM(subtotal), 0)
			FROM purchasing_details
			WHERE purchasing_id = $1
		)
		WHERE id = $1
	`
	_, err = database.DB.Exec(ctx, updateGrandTotalQuery, existingDetail.PurchasingId)
	if err != nil {
		errors.LogError("Grand total update error", err)
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Purchasing detail deleted successfully",
	})
}

