package router

import (
	db "kairon/adapters/database"
	"kairon/cmd/api/controllers"
	"kairon/repositories"
	"kairon/usecases"
	"fmt"
	"log"
	"net/http"
	"sort"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

type Server struct {
	api        *echo.Echo
	DBConn     *db.Connection
	authClient *auth.Client
}

func NewServer(db *db.Connection, authClient *auth.Client) Server {
	return Server{
		api:        echo.New(),
		DBConn:     db,
		authClient: authClient,
	}
}

func (s *Server) Run(port int) {
	s.api.Use(RequestLogger())
	s.api.Use(AuthLogger(s.authClient))
	v1 := s.api.Group("/api/v1")

	userRepository := repositories.NewUserRepository(s.DBConn)
	userUsecase := usecases.NewUserUsecase(userRepository, s.authClient)
	userHandlers := controllers.NewUserHandler(userUsecase)

	userRoutes := v1.Group("/users")
	userRoutes.Use(CheckRole([]string{"admin"}))
	{
		userRoutes.GET("/:id", userHandlers.HandleGet)
		userRoutes.POST("", validated(userHandlers.HandlePost))
		userRoutes.PUT("/:id", validatedChanges(userHandlers.HandlePut))
		userRoutes.DELETE("/:id", userHandlers.HandleDelete)
		userRoutes.GET("", userHandlers.HandleList)
	}

	printRoutes(s.api.Routes())

	v1.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm alive")
	})

	if err := s.api.Start(fmt.Sprintf(":%d", port)); err != nil {
		log.Printf("Error starting server: %s", err.Error())
	}
}

func printRoutes(routes []*echo.Route) {
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})
	for _, route := range routes {
		log.Printf("Method: %-7s Path: %-30s Name: %s\n", route.Method, route.Path, route.Name)
	}
}
