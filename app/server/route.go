package server

import (
	"github.com/quanht2k/golang_basic_training/app/api/users"
)

func SetupRoutes(server *Server) {
	api := server.Fiber.Group("/api/v1")

	userApi := users.InitUserApi(server.DB)

	user := api.Group("/users")
	user.Get("/", userApi.GetListUser())
	user.Get("/:id", userApi.GetUserDetail())
	user.Put("/:id", userApi.UpdateUser())
	user.Post("/create-user", userApi.SignUp())
	user.Post("/sign-in", userApi.SignIn())
	user.Delete("/:id", userApi.DeleteOneUser())
}