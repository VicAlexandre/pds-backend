package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/VicAlexandre/pds-backend/internal/handlers"
	"github.com/VicAlexandre/pds-backend/internal/models"
	"github.com/VicAlexandre/pds-backend/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	config Config
}

type Config struct {
	addr string
}

func (app *Application) Mount(conn *sql.DB) http.Handler {
	r := chi.NewRouter()

	/* middleware */
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	/* models */
	userModel := &models.UserModel{
		DB: conn,
	}

	/* handlers */
	authHandler := &handlers.AuthHandler{
		AuthService: services.NewAuthService(userModel),
	}

	/* routes */
	r.Route("/v1", func(r chi.Router) {
		/* health check route */
		r.Get("/health", handlers.HealthCheckHandler)

		/* authentication routes */
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)
		})
		//
		// /* user management routes */
		// r.Group(func(r chi.Router) {
		// 	r.Use(app.authMiddleware)
		//
		// 	r.Get("/me", app.getCurrentUserHandler)
		// 	r.Patch("/me", app.updateCurrentUserHandler)
		// 	r.Delete("/me", app.deleteCurrentUserHandler)
		//
		// 	r.Patch("/me/password", app.changePasswordHandler)
		// })
		//
		// r.Post("/forgot-password", app.forgotPasswordHandler)
		// r.Post("/reset-password", app.resetPasswordHandler)
	})

	return r
}

func (app *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}

	log.Println("HTTP server starting on", app.config.addr)

	return srv.ListenAndServe()
}

func NewConfig(addr string) Config {
	cfg := Config{
		addr: addr,
	}

	return cfg
}

func NewApplication(cfg Config) Application {
	app := Application{
		config: cfg,
	}

	return app
}
