package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"order-service/internal/database"
	"testing"
	"time"
)

// простой mock сервиса, реализующий интерфейс service.OrderService
type MockOrderService struct {
	orders map[string]database.Order
}

func NewMockOrderService() *MockOrderService {
	foundOrder := makeBenchmarkOrder()
	foundOrder.OrderUID = "found123"

	benchOrder := makeBenchmarkOrder()
	benchOrder.OrderUID = "test_bench"

	return &MockOrderService{
		orders: map[string]database.Order{
			"found123":   foundOrder,
			"test_bench": benchOrder,
		},
	}
}

func (m *MockOrderService) GetOrder(orderUID string) (database.Order, error) {
	order, exists := m.orders[orderUID]
	if !exists {
		return database.Order{}, fmt.Errorf("Заказ не найден!")
	}
	return order, nil
}

func (m *MockOrderService) ProcessOrder(message []byte) error {
	return nil
}

func (m *MockOrderService) ValidateOrder(order database.Order) error {
	return nil
}

func (m *MockOrderService) GetCacheSize() int {
	return len(m.orders)
}

func (m *MockOrderService) CheckDBConnection() error {
	return nil
}

func (m *MockOrderService) RunBenchmark(orderUID string) (map[string]time.Duration, error) {
	return map[string]time.Duration{"cache": time.Millisecond, "db": time.Second}, nil
}

func (m *MockOrderService) PrintCacheContents() {
	// пустая реализация для тестов
}

func TestOrderHandlerFound(t *testing.T) {
	service := NewMockOrderService()
	handler := orderHandler(service)

	// создаем тестовый запрос
	req := httptest.NewRequest("GET", "/order/found123", nil)
	w := httptest.NewRecorder()

	// вызываем хендлер
	handler(w, req)

	// проверяем ответ
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	// проверяем содержимое ответа
	var response OrderResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if response.OrderUID != "found123" {
		t.Errorf("Ожидался OrderUID 'found123', получен '%s'", response.OrderUID)
	}
	if response.TrackNumber != "TRACK001" {
		t.Errorf("Ожидался TrackNumber 'TRACK001', получен '%s'", response.TrackNumber)
	}
}

func TestOrderHandlerNotFound(t *testing.T) {
	service := NewMockOrderService()
	handler := orderHandler(service)

	// Запрос несуществующего заказа
	req := httptest.NewRequest("GET", "/order/notfound999", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Должен вернуть 404
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус 404, получен %d", w.Code)
	}

	// Проверяем сообщение об ошибке
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if response["error"] != "Order not found" {
		t.Errorf("Ожидалась ошибка 'Order not found', получено '%v'", response["error"])
	}
}

func TestCacheHandler(t *testing.T) {
	service := NewMockOrderService()
	handler := cacheHandler(service)

	req := httptest.NewRequest("GET", "/cache", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if response["cache_size"] != float64(2) {
		t.Errorf("Ожидался размер кэша 1, получен %v", response["cache_size"])
	}
}

func TestHealthHandler(t *testing.T) {
	service := NewMockOrderService()
	handler := healthHandler(service)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Ожидался статус 'healthy', получен '%v'", response["status"])
	}
	if response["cache_size"] != float64(2) {
		t.Errorf("Ожидался размер кэша 1, получен %v", response["cache_size"])
	}
}

func BenchmarkOrderHandler(b *testing.B) {
	mockService := NewMockOrderService()
	h := orderHandler(mockService)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/order/test_bench", nil)
		w := httptest.NewRecorder()

		h(w, req)
		if w.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkJsonEncodeOrder(b *testing.B) {
	order := makeBenchmarkOrder()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(&order)
		if err != nil {
			b.Fatalf("encode failed %v", err)
		}
	}
}

func BenchmarkJsonEncodeOrderResponse(b *testing.B) {
	order := makeBenchmarkOrder()
	resp := toOrderResponse(order)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(&resp)
		if err != nil {
			b.Fatalf("encode failed %v", err)
		}
	}
}

func makeBenchmarkOrder() database.Order {
	return database.Order{
		OrderUID:    "test_bench",
		TrackNumber: "TRACK001",
		Entry:       "WBIL",
		Delivery: database.Delivery{
			Name:    "Test User",
			Phone:   "+79161234567",
			Zip:     "123456",
			City:    "Moscow",
			Address: "Street 123",
			Region:  "Moscow",
			Email:   "test@example.com",
		},
		Payment: database.Payment{
			Transaction:  "b563feb7-b2b8-4b6a-9f5d-123456789abc",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []database.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest123",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		CustomerID:        "test-customer123",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
		InternalSignature: "",
	}
}
