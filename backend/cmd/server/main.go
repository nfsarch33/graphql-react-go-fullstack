package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/graph"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/graph/generated"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/internal/db"
)

const defaultPort = "8080"

func main() {
	// Initialize database
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Close database on shutdown
	defer database.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create GraphQL resolver with database
	resolvers := &graph.Resolver{
		DB: database,
	}

	// Initialize GraphQL configuration
	config := generated.Config{
		Resolvers: resolvers,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	// Setup router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// CORS middleware for React frontend
	router.Use(corsMiddleware)

	// GraphQL Playground (browser-based IDE for testing)
	router.Handle("/graphql", playground.Handler("GraphQL Playground", "/query"))

	// GraphQL Query endpoint
	router.Handle("/query", srv)

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("Server running on http://localhost:%s", port)
	log.Printf("GraphQL Playground: http://localhost:%s/graphql", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

// corsMiddleware handles CORS for React frontend
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
