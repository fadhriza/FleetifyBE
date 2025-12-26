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

type CreateItemRequest struct {
	Name     string  `json:"name" validate:"required"`
	Stock    int     `json:"stock"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Category string  `json:"category"`
	Unit     string  `json:"unit"`
	MinStock int     `json:"min_stock"`
}

type UpdateItemRequest struct {
	Name     *string  `json:"name"`
	Stock    *int     `json:"stock"`
	Price    *float64 `json:"price"`
	Category *string  `json:"category"`
	Unit     *string  `json:"unit"`
	MinStock *int     `json:"min_stock"`
}

func GetItems(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT items_id, name, stock, price, category, unit, min_stock, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		errors.LogError("Get items query error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch items",
		})
	}
	defer rows.Close()

	var items []models.Items
	for rows.Next() {
		var item models.Items
		err := rows.Scan(
			&item.ItemsId,
			&item.Name,
			&item.Stock,
			&item.Price,
			&item.Category,
			&item.Unit,
			&item.MinStock,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			errors.LogError("Item scan error", err)
			continue
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		errors.LogError("Rows iteration error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process items",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  items,
		"count": len(items),
	})
}

func GetItemById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Item ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item models.Items
	query := `
		SELECT items_id, name, stock, price, category, unit, min_stock, created_at, updated_at
		FROM items
		WHERE items_id = $1
	`

	err := database.DB.QueryRow(ctx, query, id).Scan(
		&item.ItemsId,
		&item.Name,
		&item.Stock,
		&item.Price,
		&item.Category,
		&item.Unit,
		&item.MinStock,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Item not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  item,
	})
}

func CreateItem(c *fiber.Ctx) error {
	var req CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Item name is required",
		})
	}

	if req.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Price must be greater than 0",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	query := `
		INSERT INTO items (name, stock, price, category, unit, min_stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING items_id, name, stock, price, category, unit, min_stock, created_at, updated_at
	`

	var item models.Items
	err := database.DB.QueryRow(ctx, query,
		req.Name,
		req.Stock,
		req.Price,
		req.Category,
		req.Unit,
		req.MinStock,
		now,
		now,
	).Scan(
		&item.ItemsId,
		&item.Name,
		&item.Stock,
		&item.Price,
		&item.Category,
		&item.Unit,
		&item.MinStock,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		errors.LogError("Item creation error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Item created successfully",
		"data":    item,
	})
}

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Item ID is required",
		})
	}

	var req UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingItem models.Items
	checkQuery := `SELECT items_id FROM items WHERE items_id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingItem.ItemsId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Item not found",
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

	if req.Stock != nil {
		updateFields = append(updateFields, fmt.Sprintf("stock = $%d", argPos))
		args = append(args, *req.Stock)
		argPos++
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Price must be greater than 0",
			})
		}
		updateFields = append(updateFields, fmt.Sprintf("price = $%d", argPos))
		args = append(args, *req.Price)
		argPos++
	}

	if req.Category != nil {
		updateFields = append(updateFields, fmt.Sprintf("category = $%d", argPos))
		args = append(args, *req.Category)
		argPos++
	}

	if req.Unit != nil {
		updateFields = append(updateFields, fmt.Sprintf("unit = $%d", argPos))
		args = append(args, *req.Unit)
		argPos++
	}

	if req.MinStock != nil {
		updateFields = append(updateFields, fmt.Sprintf("min_stock = $%d", argPos))
		args = append(args, *req.MinStock)
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
		UPDATE items
		SET %s
		WHERE items_id = $%d
		RETURNING items_id, name, stock, price, category, unit, min_stock, created_at, updated_at
	`, strings.Join(updateFields, ", "), argPos)

	var item models.Items
	err = database.DB.QueryRow(ctx, query, args...).Scan(
		&item.ItemsId,
		&item.Name,
		&item.Stock,
		&item.Price,
		&item.Category,
		&item.Unit,
		&item.MinStock,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		errors.LogError("Item update error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update item",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Item updated successfully",
		"data":    item,
	})
}

func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Item ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingItem models.Items
	checkQuery := `SELECT items_id FROM items WHERE items_id = $1`
	err := database.DB.QueryRow(ctx, checkQuery, id).Scan(&existingItem.ItemsId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Item not found",
		})
	}

	deleteQuery := `DELETE FROM items WHERE items_id = $1`
	_, err = database.DB.Exec(ctx, deleteQuery, id)
	if err != nil {
		errors.LogError("Item deletion error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete item",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Item deleted successfully",
	})
}

