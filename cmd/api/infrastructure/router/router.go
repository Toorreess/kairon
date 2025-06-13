package router

import (
	"fmt"
	db "kairon/adapters/database"
	"kairon/cmd/api/controllers"
	"kairon/repositories"
	"kairon/usecases"
	"log"
	"net/http"
	"sort"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	s.api.HideBanner = true
	s.api.HidePort = true
	s.api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.api.Use(RequestLogger())
	s.api.Use(AuthLogger(s.authClient))
	v1 := s.api.Group("/api/v1")

	userRepository := repositories.NewUserRepository(s.DBConn)
	userUsecase := usecases.NewUserUsecase(userRepository, s.authClient)
	userHandlers := controllers.NewUserHandler(userUsecase)

	/* Users */
	userRoutes := v1.Group("/users")
	userRoutes.Use(CheckRole([]string{"admin"}))
	{
		userRoutes.GET("/:id", userHandlers.HandleGet)
		userRoutes.POST("", validated(userHandlers.HandlePost))
		userRoutes.PUT("/:id", validatedChanges(userHandlers.HandlePut))
		userRoutes.DELETE("/:id", userHandlers.HandleDelete)
		userRoutes.GET("", userHandlers.HandleList)
	}

	/* Products */
	productRepository := repositories.NewProductRepository(s.DBConn)
	productUsecase := usecases.NewProductUsecase(productRepository)
	productHandlers := controllers.NewProductHandler(productUsecase)

	productRoutes := v1.Group("/products")
	{
		productRoutes.GET("/:id", productHandlers.HandleGet)
		productRoutes.POST("", validated(productHandlers.HandlePost))
		productRoutes.PUT("/:id", validatedChanges(productHandlers.HandlePut))
		productRoutes.DELETE("/:id", productHandlers.HandleDelete)
		productRoutes.GET("", productHandlers.HandleList)
	}

	/* Activities */
	activityRepository := repositories.NewActivityRepository(s.DBConn)
	activityUsecase := usecases.NewActivityUsecase(activityRepository)
	activityHandlers := controllers.NewActivityHandler(activityUsecase)

	activityRoutes := v1.Group("/activities")
	{
		activityRoutes.GET("/:id", activityHandlers.HandleGet)
		activityRoutes.POST("", validated(activityHandlers.HandlePost))
		activityRoutes.PUT("/:id", validatedChanges(activityHandlers.HandlePut))
		activityRoutes.DELETE("/:id", activityHandlers.HandleDelete)
		activityRoutes.GET("", activityHandlers.HandleList)
	}

	/* Members */
	memberRepository := repositories.NewMemberRepository(s.DBConn)
	memberUsecase := usecases.NewMemberUsecase(memberRepository)
	memberHandlers := controllers.NewMemberHandler(memberUsecase)

	memberRoutes := v1.Group("/members")
	{
		memberRoutes.GET("/:id", memberHandlers.HandleGet)
		memberRoutes.POST("", validated(memberHandlers.HandlePost))
		memberRoutes.PUT("/:id", validatedChanges(memberHandlers.HandlePut))
		memberRoutes.DELETE("/:id", memberHandlers.HandleDelete)
		memberRoutes.GET("", memberHandlers.HandleList)
	}

	/* Order */
	orderRepository := repositories.NewOrderRepository(s.DBConn)
	orderUsecase := usecases.NewOrderUsecase(orderRepository, productRepository)
	orderHandlers := controllers.NewOrderHandler(orderUsecase)

	orderRoutes := v1.Group("/orders")
	{
		orderRoutes.GET("/:id", orderHandlers.HandleGet)
		orderRoutes.POST("", validated(orderHandlers.HandlePost))
		orderRoutes.DELETE("/:id", orderHandlers.HandleDelete)
		orderRoutes.GET("", orderHandlers.HandleList)
		orderRoutes.PUT("/:id/pay", orderHandlers.HandlePay)
		orderRoutes.PUT("/:id/cancel", orderHandlers.HandleCancel)
	}

	/* Financial report */
	reportUsecase := usecases.NewReportUsecase(orderRepository)
	reportHandlers := controllers.NewReportHandler(reportUsecase)

	reportRoutes := v1.Group("/reports")
	{
		reportRoutes.GET("/financial", reportHandlers.HandleGetFinancialReport)
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
