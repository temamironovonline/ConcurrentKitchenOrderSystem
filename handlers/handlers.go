package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"main.go/chef"
	"main.go/dish"
)

const (
	// DefaultTimeOrder хранит в себе стандартный таймаут на одно блюдо
	DefaultTimeoutOrder = 30 * time.Second
)

// HTTPHandlers структура, отвечающая за хендлеры
type HTTPHandlers struct {
	ctx                    context.Context
	c                      *chef.Chef
	listAvailabilityDishes *dish.ListAvailabilityDishes
}

// NewHTTPHandlers создает экземпляр структуры HTTPHandlers
func NewHTTPHandlers(ctx context.Context, c *chef.Chef, listAvailabilityDishes *dish.ListAvailabilityDishes) *HTTPHandlers {
	return &HTTPHandlers{
		ctx:                    ctx,
		c:                      c,
		listAvailabilityDishes: listAvailabilityDishes,
	}
}

func writeAnswerJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("внутренняя ошибка конвертации ответа JSON")
	}
}

func writeAnswerError(w http.ResponseWriter, text string, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonErr := NewErrorHTTP(text, time.Now())
	if err := json.NewEncoder(w).Encode(&jsonErr); err != nil {
		log.Println("внутренняя ошибка конвертации ответа ошибки")
	}
}

// HandleCreateOrder хендлер создания заказа
func (h *HTTPHandlers) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var listDTO ListDishesDTO
	if err := json.NewDecoder(r.Body).Decode(&listDTO); err != nil {
		writeAnswerError(w, "не удается перевести полученные  данные в json", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	for _, v := range listDTO.List {
		if _, ok := h.listAvailabilityDishes.Get(v); !ok {
			writeAnswerError(w, "такого блюда в меню нет", http.StatusBadRequest)
			return
		}
	}

	id, err := h.c.CreateOrder(h.ctx, DefaultTimeoutOrder, listDTO.List)
	if err != nil {
		writeAnswerError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ans := NewReturnID(id, time.Now())
	writeAnswerJSON(w, ans, http.StatusCreated)
}

func getPathValue(r *http.Request) (int, bool) {
	value := r.PathValue("id")

	id, err := strconv.Atoi(value)

	if err != nil {
		return 0, false
	}

	return id, true
}

// HandleGetIntermediateResult хендлер для получения информации о конкретном блюде
func (h *HTTPHandlers) HandleGetIntermediateResult(w http.ResponseWriter, r *http.Request) {
	id, ok := getPathValue(r)
	if !ok {
		writeAnswerError(w, "ошибка конвертации в id", http.StatusBadRequest)
		return
	}

	snap, err := h.c.GetOrder(id)
	if err != nil {
		writeAnswerError(w, err.Error(), http.StatusNotFound)
		return
	}

	orderDTO := NewOrderDTO(snap)

	writeAnswerJSON(w, orderDTO, http.StatusOK)
}

// HandleCancelOrder хендлер для отмены заказа
func (h *HTTPHandlers) HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := getPathValue(r)
	if !ok {
		writeAnswerError(w, "ошибка конвертации в id", http.StatusBadRequest)
		return
	}

	snap, err := h.c.CancelOrder(id)
	if err != nil {
		writeAnswerError(w, err.Error(), http.StatusNotFound)
		return
	}

	orderDTO := NewOrderDTO(snap)
	writeAnswerJSON(w, orderDTO, http.StatusOK)
}

// HandleGetStatistics используется для получения статистики
func (h *HTTPHandlers) HandleGetStatistics(w http.ResponseWriter, r *http.Request) {
	stats := h.c.GetStatistics()

	if stats == nil {
		writeAnswerError(w, "статистика пока не создана", http.StatusInternalServerError)
		return
	}

	answer := NewStats(stats)

	writeAnswerJSON(w, answer, http.StatusOK)
}
