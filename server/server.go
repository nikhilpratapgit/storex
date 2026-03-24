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
	router.Use(middleware.CORSMiddleware)
	router.Route("/v1", func(v1 chi.Router) {
		v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, map[string]string{
				"status": "server is running",
			})
		})

		v1.Route("/auth", func(v1 chi.Router) {
			v1.Use(middleware.CORSMiddleware)
			v1.Post("/register", handler.RegisterUser)
			v1.Post("/login", handler.LoginUser)
			v1.Post("/firebase", handler.FirebaseRegister)
			v1.Post("/google", handler.GoogleLogin)
			v1.Post("/complete-user-profile", handler.CompleteUserProfile) // complete user profile
		})
		v1.Group(func(v1 chi.Router) {
			v1.Use(middleware.Auth)
			v1.Delete("/logout", handler.Logout)
			v1.Route("/users", func(v1 chi.Router) {
				v1.Get("/{id}", handler.FetchUser)
			})

			v1.Group(func(v1 chi.Router) {
				v1.Use(middleware.RoleMiddleware("admin", "asset-manager"))
				v1.Route("/assets", func(v1 chi.Router) {
					v1.Post("/", handler.CreateAsset)
					v1.Get("/", handler.ShowAssets)
					v1.Put("/assign/{id}", handler.AssignedAssets)

					v1.Patch("/unassign/{id}", handler.UnassignedAllAsset)
					v1.Put("/service/{id}", handler.ServiceAssets)
					v1.Put("/update/{id}", handler.UpdateAsset)
					v1.Delete("/delete/{id}", handler.DeleteAsset)
					v1.Patch("/update-asset-status/{id}", handler.UpdateAssetStatus)

					v1.Patch("/service-complete/{id}", handler.ServiceComplete)

				})

				v1.Group(func(v1 chi.Router) {
					v1.Use(middleware.RoleMiddleware("admin"))
					v1.Get("/users-info", handler.GetAllUsers)
					v1.Patch("/assign-role/{id}", handler.AssignedRole)
					v1.Delete("/delete-user", handler.DeleteUser)
				})
			})
		})
	})
	return &Server{
		Router: router,
	}
}
