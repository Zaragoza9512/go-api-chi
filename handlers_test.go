package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test 1: Login Handler debe retornar un token
func TestLoginHandler(t *testing.T) {
	// Preparar request
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	body, _ := json.Marshal(loginReq)
	
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Preparar response recorder
	rr := httptest.NewRecorder()
	
	// Ejecutar handler
	handler := http.HandlerFunc(LogingHandler)
	handler.ServeHTTP(rr, req)
	
	// Verificar status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler retornó status incorrecto: got %v want %v", status, http.StatusOK)
	}
	
	// Verificar que retorna un token
	var response LogingResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("No se pudo decodear respuesta JSON: %v", err)
	}
	
	if response.Token == "" {
		t.Error("El token no puede estar vacío")
	}
	
	t.Logf("✅ Test pasó - Token generado correctamente")
}

// Test 2: Login con body inválido debe retornar 400
func TestLoginHandlerInvalidJSON(t *testing.T) {
	// Request con JSON inválido
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(LogingHandler)
	handler.ServeHTTP(rr, req)
	
	// Debe retornar 400 Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler retornó status incorrecto: got %v want %v", status, http.StatusBadRequest)
	}
	
	t.Log("✅ Test pasó - JSON inválido rechazado correctamente")
}

// Test 3: Validar estructura de Product
func TestProductStruct(t *testing.T) {
	product := Product{
		ID:          1,
		Name:        "Laptop",
		Description: "Gaming laptop",
		Price:       1500.50,
		Stock:       10,
	}
	
	if product.ID != 1 {
		t.Errorf("Product ID incorrecto: got %v want %v", product.ID, 1)
	}
	
	if product.Price <= 0 {
		t.Error("Product price debe ser mayor a 0")
	}
	
	t.Log("✅ Test pasó - Estructura Product válida")
}
