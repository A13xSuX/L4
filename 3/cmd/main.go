package main

import (
	"4_3/internal/handlers"
	"4_3/internal/middleware"
	"4_3/internal/myLogger"
	"4_3/internal/repository"
	"4_3/internal/service"
	"4_3/internal/worker"
	"flag"
	"net/http"
	"time"
)

func main() {
	logger := myLogger.NewLogger()
	go logger.Start()

	portFlag := flag.String("port", "8080", "port to run server")
	flag.Parse()

	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := service.NewEventService(eventRepo, remindCh)
	eventHandler := handlers.NewEventHandler(eventService, logger)

	remindWorker := worker.NewRemindWorker(remindCh, logger)
	go remindWorker.Start()

	archiveInterval := 30 * time.Second
	archiveWorker := worker.NewArchiveWorker(eventRepo, archiveInterval, logger)
	go archiveWorker.Start()

	http.Handle("/status", middleware.LoggingMiddleware(http.HandlerFunc(eventHandler.Status), logger))
	http.Handle("/create_event", middleware.LoggingMiddleware(http.HandlerFunc(eventHandler.CreateEvent), logger))
	http.Handle("/events_for_day", middleware.LoggingMiddleware(http.HandlerFunc(eventHandler.EventsForDay), logger))
	http.Handle("/events_for_week", middleware.LoggingMiddleware(http.HandlerFunc(eventHandler.EventsForWeek), logger))
	http.Handle("/events_for_month", middleware.LoggingMiddleware(http.HandlerFunc(eventHandler.EventsForMonth), logger))
	http.Handle("/delete_event", middleware.LoggingMiddleware(http.HandlerFunc(eventHandler.Delete), logger))
	http.Handle("/update_event", middleware.LoggingMiddleware(http.HandlerFunc(eventHandler.UpdateEvent), logger))

	addr := ":" + *portFlag
	msg := "Server starting on " + addr
	logger.Info(msg)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error("Server failed", err)
		return
	}

}
