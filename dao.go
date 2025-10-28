package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Nota: La estructura 'Product' (producto) se define en handlers.go.

// ====================================================================
// DAO (Data Access Object)
// Lógica que interactúa directamente con la base de datos (PostgreSQL).
// ====================================================================

// CreateProduct (Crear Producto): Inserta un nuevo producto y devuelve el producto con el ID asignado.
func CreateProduct(db *sql.DB, product Product, userID int) (Product, error) {

	// ⬇️ CAMBIO 2: Incluir la nueva columna (creator_id) y el nuevo placeholder ($5)
	sqlStatement := `
		INSERT INTO products (name, description, price, stock, creator_id) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id int
	err := db.QueryRow(
		sqlStatement,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		userID, // ⬅️ CAMBIO 3: Pasar el userID como quinto argumento de la consulta
	).Scan(&id)

	if err != nil {
		return Product{}, fmt.Errorf("error al ejecutar INSERT en DB: %w", err)
	}

	product.ID = id
	// Nota: Si quieres devolver el UserID, debes añadirlo al struct Product.
	// product.CreatorID = userID

	return product, nil
}

// GetProducts (Obtener Todos): Consulta y devuelve todos los productos.
func GetProducts(db *sql.DB) ([]Product, error) {
	sqlStatement := `SELECT id, name, description, price, stock FROM products ORDER BY id`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar SELECT ALL en DB: %w", err)
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		// Escanea los resultados de la fila actual
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock)
		if err != nil {
			log.Printf("Error al escanear fila de producto: %v", err)
			continue
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error después de iterar filas: %w", err)
	}

	return products, nil
}

// GetProductByID (Obtener por ID): Consulta y devuelve un producto específico por su ID.
func GetProductByID(db *sql.DB, id int) (Product, error) {
	sqlStatement := `SELECT id, name, description, price, stock FROM products WHERE id = $1`
	var p Product

	// QueryRow se usa para cuando se espera una sola fila.
	err := db.QueryRow(sqlStatement, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock)

	if err != nil {
		// sql.ErrNoRows es manejado directamente por el handler para devolver 404
		return Product{}, err
	}

	return p, nil
}

// UpdateProduct (Actualizar Producto): Actualiza un producto existente.
func UpdateProduct(db *sql.DB, product Product) error {
	sqlStatement := `
		UPDATE products
		SET name = $2, description = $3, price = $4, stock = $5
		WHERE id = $1`

	result, err := db.Exec(
		sqlStatement,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	)
	if err != nil {
		return fmt.Errorf("error al ejecutar UPDATE en DB: %w", err)
	}

	// LÓGICA DE 404: Verificar si se afectó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al leer filas afectadas: %w", err)
	}
	if rowsAffected == 0 {
		// Devolvemos un error específico para que el handler lo mapee a 404
		return fmt.Errorf("producto con ID %d no encontrado", product.ID)
	}

	return nil
}

// DeleteProduct (Eliminar Producto): Elimina un producto por su ID.
func DeleteProduct(db *sql.DB, id int) error {
	sqlStatement := `DELETE FROM products WHERE id = $1`

	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("error al ejecutar DELETE en DB: %w", err)
	}

	// LÓGICA DE 404: Verificar si se afectó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al leer filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		// Devolvemos un error específico para que el handler lo mapee a 404
		return fmt.Errorf("producto con ID %d no encontrado", id)
	}

	return nil
}

/*
 * CLASE: ENRUTAMIENTO PROFESIONAL (DAO - ROBUSTEZ)
 *
 * OBJETIVO: Mejorar la calidad y el manejo de errores del Data Access Object (DAO).
 *
 * CAMBIOS CLAVE EN DAO:
 * 1. Manejo de 404 (Not Found): Las funciones UpdateProduct y DeleteProduct ahora verifican
 * 'result.RowsAffected()'. Si es 0, devuelven un error específico ("no encontrado").
 * Esto permite al handler superior devolver el código HTTP 404 correcto.
 * 2. Encapsulamiento de Errores: Uso de 'fmt.Errorf("mensaje útil: %w", err)' para envolver
 * los errores de SQL. Esto facilita el rastreo de fallos y el testing profesional.
 * 3. Firma de Funciones: Se estandarizó 'CreateProduct' para devolver el objeto 'Product'
 * completo (con ID) en lugar de solo el ID, mejorando la coherencia de la API.
 */
