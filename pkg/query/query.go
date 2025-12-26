package query

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type QueryParams struct {
	Page     int
	Limit    int
	Sort     string
	SortDir  string
	Search   string
	Filters  map[string]string
}

func ParseQueryParams(c *fiber.Ctx) QueryParams {
	params := QueryParams{
		Page:    1,
		Limit:   10,
		Sort:    "created_at",
		SortDir: "DESC",
		Filters: make(map[string]string),
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			params.Limit = limit
		}
	}

	if sort := c.Query("sort"); sort != "" {
		if strings.HasPrefix(sort, "-") {
			params.Sort = sort[1:]
			params.SortDir = "DESC"
		} else {
			params.Sort = sort
			params.SortDir = "ASC"
		}
	}

	if search := c.Query("search"); search != "" {
		params.Search = strings.TrimSpace(search)
	}

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if strings.HasPrefix(keyStr, "filter[") && strings.HasSuffix(keyStr, "]") {
			field := keyStr[7 : len(keyStr)-1]
			params.Filters[field] = string(value)
		} else if strings.HasPrefix(keyStr, "filter_") {
			field := keyStr[7:]
			params.Filters[field] = string(value)
		}
	})

	return params
}

func BuildWhereClause(params QueryParams, searchFields []string, filterFields map[string]string) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argPos := 1

	if params.Search != "" && len(searchFields) > 0 {
		searchConditions := []string{}
		for _, field := range searchFields {
			searchConditions = append(searchConditions, fmt.Sprintf("%s ILIKE $%d", field, argPos))
		}
		args = append(args, "%"+params.Search+"%")
		argPos++
		conditions = append(conditions, "("+strings.Join(searchConditions, " OR ")+")")
	}

	for field, value := range params.Filters {
		if dbField, ok := filterFields[field]; ok {
			conditions = append(conditions, fmt.Sprintf("%s = $%d", dbField, argPos))
			args = append(args, value)
			argPos++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}

func BuildOrderClause(params QueryParams, defaultSort string) string {
	sortField := params.Sort
	if sortField == "" {
		sortField = defaultSort
	}
	return fmt.Sprintf("ORDER BY %s %s", sortField, params.SortDir)
}

func BuildPaginationClause(params QueryParams, argPos int) (string, []interface{}) {
	offset := (params.Page - 1) * params.Limit
	return fmt.Sprintf("LIMIT $%d OFFSET $%d", argPos, argPos+1), []interface{}{params.Limit, offset}
}

func BuildCountQuery(tableName string, whereClause string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereClause)
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Count      int         `json:"count"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	TotalCount int         `json:"total_count"`
}

func NewPaginatedResponse(data interface{}, totalCount int, page int, limit int) PaginatedResponse {
	totalPages := (totalCount + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	
	var dataCount int
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		dataCount = v.Len()
	}
	
	return PaginatedResponse{
		Data:       data,
		Count:      dataCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}
}

