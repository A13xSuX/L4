package main

import (
	"4_3/internal/models"
	"4_3/internal/service"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var cal = service.NewCalendar()

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed) //устанавливает правильный статус
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"result" : "server is running"}`)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Вызов оригинального handler'а
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		fmt.Printf("Method: %s, URL: %s, Time: %v\n", r.Method, r.URL.Path, duration)
	})
}

// CRUD
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // ответ будет в json формате
	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userId := r.FormValue("user_id")
	date := r.FormValue("date")
	description := r.FormValue("description")
	priority := r.FormValue("priority")
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
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "date is empty"}`)
		return
	}
	if description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "description is empty"}`)
		return
	}
	if priority == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "priority is empty"}`)
		return
	}
	event := models.Event{
		UserID:      userIdInt,
		EventTime:   date,
		Description: description,
		Priority:    priority,
	}
	cal.CreateEvent(userIdInt, date, event)
	fmt.Fprintf(w, `{"result" : "event created"}`)
}

func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
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
	events := cal.GetEventsForDay(userIdInt, date)
	errJson := json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if errJson != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
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

	events, err := cal.GetEventsForWeek(userIdInt, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "GetEventsForWeek"}`)
		return
	}
	errJson := json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if errJson != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
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

	events, err := cal.GetEventsForMonth(userIdInt, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "GetEventsForMonth"}`)
		return
	}
	errJson := json.NewEncoder(w).Encode(map[string]interface{}{"result": events})
	if errJson != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userId := r.FormValue("user_id")
	date := r.FormValue("date")
	description := r.FormValue("description")

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
	if description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "description is empty"}`)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId to int"}`)
		return
	}
	if userIdInt == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is Zero"}`)
		return
	}

	success := cal.DeleteEvent(userIdInt, date, description)
	if !success {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "event not found"}`)
		return
	}

	errJson := json.NewEncoder(w).Encode(map[string]string{"result": "event deleted"})
	if errJson != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	userId := r.FormValue("user_id")
	oldDate := r.FormValue("old_date")
	oldDescription := r.FormValue("old_description")
	oldPriority := r.FormValue("old_priority")
	newDate := r.FormValue("new_date")
	newDescription := r.FormValue("new_description")
	newPriority := r.FormValue("new_priority")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "userId is empty"}`)
		return
	}
	if oldDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "date is empty"}`)
		return
	}
	if oldDescription == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "description is empty"}`)
		return
	}
	if oldPriority == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "priority is empty"}`)
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "UserId to int"}`)
		return
	}
	if userIdInt == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error" : "UserId is Zero"}`)
		return
	}

	if newDate == "" && newDescription == "" && newPriority == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "at least one new field must be provided"}`)
		return
	}

	success := cal.UpdateEvent(userIdInt, oldDate, oldDescription, oldPriority, newDate, newDescription, newPriority)
	if !success {
		w.WriteHeader(http.StatusNotFound) // 404 - событие не найдено
		fmt.Fprintf(w, `{"error": "event not found"}`)
		return
	}

	errJson := json.NewEncoder(w).Encode(map[string]string{"result": "event updated"})
	if errJson != nil {
		http.Error(w, `{"error" : "internal server error"}`, http.StatusInternalServerError)
		return
	}

}

func main() {
	portFlag := flag.String("port", "8080", "port to run server")
	flag.Parse()
	http.Handle("/status", loggingMiddleware(http.HandlerFunc(statusHandler)))
	http.Handle("/create_event", loggingMiddleware(http.HandlerFunc(createEventHandler)))
	http.Handle("/events_for_day", loggingMiddleware(http.HandlerFunc(eventsForDayHandler)))
	http.Handle("/events_for_week", loggingMiddleware(http.HandlerFunc(eventsForWeekHandler)))
	http.Handle("/events_for_month", loggingMiddleware(http.HandlerFunc(eventsForMonthHandler)))
	http.Handle("/delete_event", loggingMiddleware(http.HandlerFunc(deleteEventHandler)))
	http.Handle("/update_event", loggingMiddleware(http.HandlerFunc(updateEventHandler)))

	addr := ":" + *portFlag
	fmt.Printf("Server starting on port %s\n", *portFlag)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
