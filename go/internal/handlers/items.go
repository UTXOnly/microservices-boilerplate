package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/microservices-boilerplate/go/internal/models"
	"github.com/microservices-boilerplate/go/internal/services"
)

// ItemsHandler handles items API requests.
type ItemsHandler struct {
	svc *services.ItemService
}

// NewItemsHandler creates a new ItemsHandler.
func NewItemsHandler(svc *services.ItemService) *ItemsHandler {
	return &ItemsHandler{svc: svc}
}

// List returns paginated items.
func (h *ItemsHandler) List(c echo.Context) error {
	skip, _ := strconv.Atoi(c.QueryParam("skip"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 100
	}

	items, err := h.svc.List(c.Request().Context(), skip, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

// Get returns a single item.
func (h *ItemsHandler) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	item, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "item not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, item)
}

// Create creates a new item.
func (h *ItemsHandler) Create(c echo.Context) error {
	var data models.ItemCreate
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if data.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	item, err := h.svc.Create(c.Request().Context(), &data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, item)
}

// Update partially updates an item.
func (h *ItemsHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var data models.ItemUpdate
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item, err := h.svc.Update(c.Request().Context(), id, &data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, item)
}

// Delete removes an item.
func (h *ItemsHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
