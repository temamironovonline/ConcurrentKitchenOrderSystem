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

type HTTPHandlers struct {
	ctx                    context.Context
	c                      *chef.Chef
	listAvailabilityDishes *dish.ListAvailabilityDishes
}

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

}

func writeAnswerError(w http.ResponseWriter, text string, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonErr := NewErrorHTTP(text, time.Now())
	if err := json.NewEncoder(w).Encode(&jsonErr); err != nil {
		log.Println("внутренняя ошибка конвертации ответа ошибки")
	}
}

func (h *HTTPHandlers) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var listDTO ListDishesDTO
	if err := json.NewDecoder(r.Body).Decode(&listDTO); err != nil {
		writeAnswerError(w, "не удается перевести полученные  данные в json", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	for _, v := range listDTO.List {
		if _, ok := h.listAvailabilityDishes.Get(v); !ok {
			writeAnswerError(w, "такого блюда в меню нет", http.StatusBadRequest)
			return
		}
	}

	id := h.c.CreateOrder(h.ctx, 10, listDTO.List)

	ans := NewCreateOrderAnswerDTO(id, time.Now())
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

	orderDTO := newOrderDTO(snap)

	writeAnswerJSON(w, orderDTO, http.StatusOK)
}
