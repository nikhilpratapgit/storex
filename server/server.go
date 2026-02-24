package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nikhilpratapgit/storex/handler"
	"github.com/nikhilpratapgit/storex/middleware"
	"github.com/nikhilpratapgit/storex/utils"
)

type Server struct {
	chi.Router
	server *http.Server
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/v1", func(v1 chi.Router) {
		v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, map[string]string{
				"status": "Server is Running",
			})
		})
		//public routes
		v1.Post("/register", handler.RegisterUser)
		v1.Post("/login", handler.LoginUser)
		// auth required
		v1.Group(func(v1 chi.Router) {
			v1.Use(middleware.Auth)
			v1.Delete("/logout", handler.Logout)
			v1.Group(func(v1 chi.Router) {
				v1.Use(middleware.RoleMiddleware("admin", "asset-manager"))
				// role based
				v1.Post("/asset", handler.CreateAsset)
				v1.Get("/assets", handler.ShowAssets)

			})

		})

	})
	return &Server{
		Router: router,
	}
}
