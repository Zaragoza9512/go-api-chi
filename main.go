package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setupDB() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

func setupRouter(db *sql.DB, jwtSecretKey string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(MetricsMiddleware)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Post("/login", LoginHandler(db, jwtSecretKey))
	})

	r.Route("/productos", func(r chi.Router) {
		r.Use(AuthMiddleware(jwtSecretKey))
		r.Post("/", CreateProductHandler(db))
		r.Get("/", GetProductsHandler(db))
		r.Get("/{id}", GetProductByIDHandler(db))
		r.Put("/{id}", UpdateProductHandler(db))
		r.Delete("/{id}", DeleteProductHandler(db))
	})

	r.Handle("/metrics", promhttp.Handler())
	return r
}

func main() {
	// Carga el .env solo en desarrollo (en producción las vars ya están en el sistema)
	_ = godotenv.Load()

	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET no está definido en las variables de entorno")
	}

	db := setupDB()
	defer db.Close()

	router := setupRouter(db, jwtSecretKey)
	log.Println("Servidor escuchando en :8080...")
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
