package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings" // Necesario para strings.Contains en Delete/Update

	// Necesario para fmt.Errorf o logging
	"github.com/go-chi/chi/v5"
)

// Product: Estructura de datos del producto
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// ====================================================================
// Handlers (Manejadores de Peticiones HTTP)
// ====================================================================

// POST /productos: Crea un nuevo producto
func CreateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product

		// 1. Decodificar el cuerpo JSON
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "JSON inválido o campos faltantes", http.StatusBadRequest)
			return
		}

		// 2. Llamada al DAO para crear el producto
		createdProduct, err := CreateProduct(db, product)
		if err != nil {
			log.Printf("DB error al crear producto: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// 3. Respuesta de éxito 201 Created
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdProduct)
	}
}

// GET /productos: Obtiene la lista de todos los productos
func GetProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Llamada al DAO para obtener todos los productos
		products, err := GetProducts(db)
		if err != nil {
			log.Printf("DB error al obtener productos: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// 2. Respuesta de éxito 200 OK
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

// GET /productos/{id}: Obtiene un producto específico
func GetProductByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 🌟 CLAVE CHI: Extracción del parámetro ID sin strings.Split
		idStr := chi.URLParam(r, "id")

		// 1. Convertir el ID a entero
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "El ID debe ser un número entero válido.", http.StatusBadRequest)
			return
		}

		// 2. Llamada al DAO para obtener el producto
		product, err := GetProductByID(db, id)

		if err != nil {
			// Manejar 404 Not Found (cuando el DAO devuelve sql.ErrNoRows)
			if err == sql.ErrNoRows {
				http.Error(w, "Producto no encontrado", http.StatusNotFound)
				return
			}
			// Manejar 500 Internal Server Error
			log.Printf("DB error al obtener producto: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// 3. Respuesta de éxito 200 OK
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// PUT /productos/{id}: Actualiza un producto existente
func UpdateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 🌟 CLAVE CHI: Extracción del parámetro ID
		idStr := chi.URLParam(r, "id")

		// 1. Convertir el ID a entero
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "El ID debe ser un número entero válido.", http.StatusBadRequest)
			return
		}

		var product Product
		// 2. Decodificar el cuerpo JSON
		err = json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "JSON inválido o campos faltantes", http.StatusBadRequest)
			return
		}

		// Aseguramos que el ID de la URL se use para la actualización
		product.ID = id

		// 3. Llamada al DAO para actualizar
		err = UpdateProduct(db, product)
		if err != nil {
			// Usamos la lógica de 404 si el DAO devuelve el error específico
			if strings.Contains(err.Error(), "no encontrado") {
				http.Error(w, "Producto no encontrado.", http.StatusNotFound)
				return
			}
			log.Printf("DB error al actualizar producto: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// 4. Respuesta de éxito 200 OK (Devolver el producto actualizado)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// DELETE /productos/{id}: Elimina un producto
func DeleteProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 🌟 CLAVE CHI: Extracción del parámetro ID
		idStr := chi.URLParam(r, "id")

		// 1. Convertir el ID a entero
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "El ID debe ser un número entero válido.", http.StatusBadRequest)
			return
		}

		// 2. Llamada al DAO para eliminar
		err = DeleteProduct(db, id)
		if err != nil {
			// Usamos la lógica de 404 si el DAO devuelve el error específico
			if strings.Contains(err.Error(), "no encontrado") {
				http.Error(w, "Producto no encontrado.", http.StatusNotFound)
				return
			}
			log.Printf("DB error al eliminar producto: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// 3. Respuesta de éxito 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}
