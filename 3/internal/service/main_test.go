package service

import (
	"2_18/internal/models"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	// Arrange
	cal := NewCalendar()
	event := models.Event{
		UserId:      1,
		Date:        "2025-11-05",
		Description: "Fun",
		Priority:    "high",
	}

	// Act
	cal.CreateEvent(1, "2025-11-05", event)

	//Assert
	events := cal.GetEventsForDay(1, "2025-11-05")

	//check the number of events
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	//check the contains of event
	if events[0].Description != "Fun" {
		t.Errorf("Expected `Fun`, got `%s`", events[0].Description)
	}
}

func TestGetEventsForDayEmpty(t *testing.T) {
	//Arrange
	cal := NewCalendar()

	//Act
	events := cal.GetEventsForDay(1, "2025-11-2025")

	//Assert
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

func TestGetEventsForDayTwoUsers(t *testing.T) {
	//Arrange
	cal := NewCalendar()
	//first user
	event1 := models.Event{
		UserId:      1,
		Date:        "2025-11-05",
		Description: "First",
		Priority:    "high",
	}
	//second user
	event2 := models.Event{
		UserId:      2,
		Date:        "2025-11-05",
		Description: "Second",
		Priority:    "high",
	}

	cal.CreateEvent(1, "2025-11-05", event1)
	cal.CreateEvent(2, "2025-11-05", event2)

	//Act
	user1Event := cal.GetEventsForDay(1, "2025-11-05")
	user2Event := cal.GetEventsForDay(2, "2025-11-05")

	//Assert
	if len(user1Event) != 1 {
		t.Errorf("Expected 1 event, got %d", len(user1Event))
	}
	if len(user2Event) != 1 {
		t.Errorf("Expected 1 event, got %d", len(user2Event))
	}
	if len(user1Event) > 0 && len(user2Event) > 0 {
		if user1Event[0].Description == user2Event[0].Description {
			t.Error("Events got mixed up")
		}
	}
}

func TestGetEventsForWeek(t *testing.T) {
	//Arrange
	cal := NewCalendar()
	events := []struct {
		date string
		desc string
	}{
		{"2025-11-03", "Понедельник"}, // Начало недели
		{"2025-11-05", "Среда"},       // Текущий день
		{"2025-11-07", "Пятница"},     // Конец недели
		{"2025-11-09", "Воскресенье"}, // Последний день недели
	}

	for _, e := range events {
		event := models.Event{
			UserId:      1,
			Date:        e.date,
			Description: e.desc,
			Priority:    "medium",
		}
		cal.CreateEvent(1, e.date, event)
	}

	//Act
	eventsWeek, err := cal.GetEventsForWeek(1, "2025-11-03")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	//Assert
	if len(eventsWeek) != 4 {
		t.Errorf("Expected 4 events, got %d", len(eventsWeek))
	}
}

func TestGetEventsForWeekBoundaries(t *testing.T) {
	cal := NewCalendar()

	events := []struct {
		date string
		desc string
	}{
		{"2025-11-02", "Прошлая неделя"},
		{"2025-11-03", "Наша неделя"},
		{"2025-11-09", "Наша неделя"},
		{"2025-11-10", "Следующая неделя"},
	}

	for _, e := range events {
		event := models.Event{
			UserId:      1,
			Date:        e.date,
			Description: e.desc,
			Priority:    "medium",
		}
		cal.CreateEvent(1, e.date, event)
	}

	// Получаем неделю с 3 по 9 ноября
	eventsWeek, err := cal.GetEventsForWeek(1, "2025-11-05")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(eventsWeek) != 2 {
		t.Errorf("Expected 2 events for the week (3rd and 9th), got %d", len(eventsWeek))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	//Arrange
	cal := NewCalendar()
	events := []struct {
		date string
		desc string
	}{
		{"2025-11-03", "Понедельник"},
		{"2025-11-12", "Среда"},
		{"2025-11-28", "Пятница"},
		{"2025-11-30", "Воскресенье"},
	}

	for _, e := range events {
		event := models.Event{
			UserId:      1,
			Date:        e.date,
			Description: e.desc,
			Priority:    "medium",
		}
		cal.CreateEvent(1, e.date, event)
	}

	//Act
	eventsMonth, err := cal.GetEventsForMonth(1, "2025-11-03")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	//Assert
	if len(eventsMonth) != 4 {
		t.Errorf("Expected 4 events, got %d", len(eventsMonth))
	}
}

func TestGetEventsForMonthBoundaries(t *testing.T) {
	//Arrange
	cal := NewCalendar()
	events := []struct {
		date string
		desc string
	}{
		{"2025-10-03", "Другой месяц"},
		{"2025-11-12", "Среда"},
		{"2025-11-28", "Пятница"},
		{"2025-12-30", "Другой месяц"},
	}

	for _, e := range events {
		event := models.Event{
			UserId:      1,
			Date:        e.date,
			Description: e.desc,
			Priority:    "medium",
		}
		cal.CreateEvent(1, e.date, event)
	}

	//Act
	eventsMonth, err := cal.GetEventsForMonth(1, "2025-11-03")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	//Assert
	if len(eventsMonth) != 2 {
		t.Errorf("Expected 2 events, got %d", len(eventsMonth))
	}
}

func TestDeleteEvent(t *testing.T) {
	//Arrange
	cal := NewCalendar()
	event1 := models.Event{
		UserId:      1,
		Date:        "2025-11-05",
		Description: "First",
		Priority:    "high",
	}
	cal.CreateEvent(1, "2025-11-5", event1)

	//Act
	eventBefore := cal.GetEventsForDay(1, "2025-11-5")
	if len(eventBefore) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(eventBefore))
	}
	success := cal.DeleteEvent(1, "2025-11-5", "First")
	if !success {
		t.Errorf("Not delete")
	}

	eventAfter := cal.GetEventsForDay(1, "2025-11-5")
	if len(eventAfter) != 0 {
		t.Errorf("Expected 0 event, got %d", len(eventAfter))
	}
}

func TestUpdateEvent(t *testing.T) {
	//Arrange
	cal := NewCalendar()

	event := models.Event{
		UserId:      1,
		Date:        "2025-11-05",
		Description: "First",
		Priority:    "high",
	}

	cal.CreateEvent(1, "2025-11-05", event)

	eventBeforeUpdate := cal.GetEventsForDay(1, "2025-11-05")
	if len(eventBeforeUpdate) != 1 {
		t.Fatalf("Expected 1 event,got %d", len(eventBeforeUpdate))
	}

	//Act
	success := cal.UpdateEvent(1, "2025-11-05", "First", "high", "2025-11-05", "Bad", "low")

	//Assert
	if !success {
		t.Fatal("Update should return true for successful update")
	}
	eventAfterUpdate := cal.GetEventsForDay(1, "2025-11-05")
	if len(eventAfterUpdate) != 1 {
		t.Fatalf("Expected 1 event,got %d", len(eventBeforeUpdate))
	}
	if eventAfterUpdate[0].Description != "Bad" {
		t.Errorf("Expected description - Bad, but got description - %s", eventAfterUpdate[0].Description)
	}
	if eventAfterUpdate[0].Priority != "low" {
		t.Errorf("Expected priority - low, but got priority - %s", eventAfterUpdate[0].Priority)
	}
}

func TestUpdateEventChangeDate(t *testing.T) {
	cal := NewCalendar()
	event := models.Event{
		UserId:      1,
		Date:        "2025-11-05",
		Description: "Meeting",
		Priority:    "high",
	}
	cal.CreateEvent(1, "2025-11-05", event)

	success := cal.UpdateEvent(1, "2025-11-05", "Meeting", "high", "2025-11-06", "", "")
	if !success {
		t.Fatal("Update should succeed")
	}

	// Проверяем, что событие переместилось
	oldDateEvents := cal.GetEventsForDay(1, "2025-11-05")
	newDateEvents := cal.GetEventsForDay(1, "2025-11-06")

	if len(oldDateEvents) != 0 {
		t.Error("Should be no events in old date")
	}
	if len(newDateEvents) != 1 {
		t.Error("Should be 1 event in new date")
	}
}
