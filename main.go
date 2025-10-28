package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq"
)

// Constantes de conexión
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "ecom_db"
)

// Clave Secreta JWT
const jwtSecretKey = "ClaveSecretaParaFirmarJWT123456789!@#$"

// Función auxiliar para crear la conexión a la DB
func setupDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a la DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error al hacer ping a la DB: %v", err)
	}

	log.Println("Conexión a la base de datos establecida exitosamente.")
	return db
}

// Router que incluye TODAS las rutas y la lógica de seguridad
func setupRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Middleware CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// --- 1. RUTAS PÚBLICAS Y DE AUTENTICACIÓN ---
	r.Group(func(r chi.Router) {
		// La ruta de LOGIN debe ser pública
		r.Post("/login", LogingHandler)
	})

	// --- 2. RUTAS PROTEGIDAS CON AUTH MIDDLEWARE (Productos) ---
	r.Route("/productos", func(r chi.Router) {

		// ⬇️ Aplicamos el VIGILANTE JWT a todas las rutas que siguen
		r.Use(AuthMiddleware(jwtSecretKey))

		// Rutas CRUD Protegidas. Solo se accede si el JWT es válido
		r.Post("/", CreateProductHandler(db))
		r.Get("/", GetProductsHandler(db))
		r.Get("/{id}", GetProductByIDHandler(db))
		r.Put("/{id}", UpdateProductHandler(db))
		r.Delete("/{id}", DeleteProductHandler(db))
	})

	return r
}

// FUNCIÓN MAIN FINAL
func main() {
	// 1. Conexión a la DB
	db := setupDB()
	defer db.Close()

	// 2. Configuración de Rutas y Seguridad
	router := setupRouter(db) // Solo una vez

	log.Println("Servidor escuchando en :8080...")
	// 3. Iniciar el Servidor
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

/*
 * CLASE: ENRUTAMIENTO PROFESIONAL (FASE 3 - CHI)
 *
 * OBJETIVO: Reemplazar el router básico (http.ServeMux) por 'github.com/go-chi/chi/v5'
 * para obtener código modular, limpio y escalable.
 *
 * CAMBIOS CLAVE:
 * 1. setupRouter(): Esta función centraliza toda la configuración del enrutamiento.
 * 2. r.METODO(ruta, handler): Chi registra rutas por método HTTP (r.Post, r.Get, etc.),
 * eliminando la necesidad de usar 'if r.Method == "..."' dentro de los handlers.
 * 3. Rutas Dinámicas: Se usa la sintaxis '/productos/{id}' para definir variables de URL.
 *
 * PRÓXIMOS PASOS:
 * - Se integrará el Middleware (r.Use) para Seguridad (CORS) y Logging.
 * - Se completará la función main() para iniciar el servidor con el router 'chi'.
 */
