package handlers

import (
	"4_3/internal/dto"
	"4_3/internal/myLogger"
	"4_3/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type EventHandler struct {
	eventService *service.EventService
	logger       *myLogger.Logger
}

func NewEventHandler(eventService *service.EventService, logger *myLogger.Logger) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		logger:       logger,
	}
}

func (h *EventHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed) //устанавливает правильный статус
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"result" : "server is running"}`)
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // ответ будет в json формате
	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userId := r.FormValue("user_id")
	title := r.FormValue("title")
	date := r.FormValue("date")
	description := r.FormValue("description")
	priority := r.FormValue("priority")
	remindAt := r.FormValue("remind_at")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is empty"}`)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "user_id to number"}`)
		return
	}
	if userIdInt == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is Zero"}`)
		return
	}
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "title empty"}`)
		return
	}
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "date is empty"}`)
		return
	}
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "invalid date"}`)
		return
	}
	if priority == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "priority is empty"}`)
		return
	}
	var remindAtTime *time.Time
	if remindAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, remindAt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error" : "invalid remind_at"}`)
			return
		}
		remindAtTime = &parsedTime
	}
	eventReq := dto.CreateEventRequest{
		UserID:      userIdInt,
		Title:       title,
		EventTime:   dateTime,
		Description: &description,
		Priority:    priority,
		RemindAt:    remindAtTime,
	}
	id, err := h.eventService.Create(eventReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("error in moment create new event", err)
		fmt.Fprintf(w, `{"error" : "failed create new event"}`)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": "event created",
		"id": id})
	if err != nil {
		h.logger.Error("failed to encode response", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func (h *EventHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userId := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is empty"}`)
		return
	}
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "date is empty"}`)
		return
	}
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "invalid date"}`)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "user_id to number"}`)
		return
	}
	if userIdInt == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is Zero"}`)
		return
	}

	events, err := h.eventService.GetEventsForDay(userIdInt, dateTime)
	if err != nil {
		h.logger.Error("Failed to get events for day", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
	errJson := json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if errJson != nil {
		h.logger.Error("failed to encode response", errJson)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

}

func (h *EventHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userId := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is empty"}`)
		return
	}
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "date is empty"}`)
		return
	}
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "invalid date"}`)
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "user_id to number"}`)
		return
	}
	if userIdInt == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is Zero"}`)
		return
	}

	events, err := h.eventService.GetEventsForWeek(userIdInt, dateTime)
	if err != nil {
		h.logger.Error("failed to get events for week", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
	errJson := json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if errJson != nil {
		h.logger.Error("failed to encode response", errJson)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func (h *EventHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userId := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is empty"}`)
		return
	}
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "date is empty"}`)
		return
	}
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "invalid date"}`)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "user_id to number"}`)
		return
	}
	if userIdInt == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is Zero"}`)
		return
	}

	events, err := h.eventService.GetEventsForMonth(userIdInt, dateTime)
	if err != nil {
		h.logger.Error("failed to get events for month", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
	errJson := json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if errJson != nil {
		h.logger.Error("failed to encode response", errJson)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")
	title := r.FormValue("title")
	date := r.FormValue("date")
	description := r.FormValue("description")
	priority := r.FormValue("priority")
	remindAt := r.FormValue("remind_at")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "id is empty"}`)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "id to int"}`)
		return
	}
	if idInt <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "id less or equal zero"}`)
		return
	}

	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "title empty"}`)
		return
	}
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "date empty"}`)
		return
	}
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "invalid date"}`)
		return
	}
	if priority == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "priority empty"}`)
		return
	}
	if remindAt == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "remind_at empty"}`)
		return
	}
	remindAtTime, err := time.Parse(time.RFC3339, remindAt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "invalid remind_at"}`)
		return
	}

	newEvent := dto.UpdateEventRequest{
		Title:       title,
		EventTime:   dateTime,
		Description: &description,
		Priority:    priority,
		RemindAt:    &remindAtTime,
	}
	err = h.eventService.Update(idInt, newEvent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("failed to update event", err)
		fmt.Fprintf(w, `{"error" : "failed to update event"}`)
		return
	}

	errJson := json.NewEncoder(w).Encode(map[string]string{"result": "event updated"})
	if errJson != nil {
		h.logger.Error("failed to encode response", errJson)
		http.Error(w, `{"error" : "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "id is empty"}`)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "id to int"}`)
		return
	}
	if idInt <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "id less or equal zero"}`)
		return
	}

	err = h.eventService.Delete(idInt)
	if err != nil {
		h.logger.Error("Failed to delete event", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	errJson := json.NewEncoder(w).Encode(map[string]string{"result": "event deleted"})
	if errJson != nil {
		h.logger.Error("failed to encode response", errJson)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}
