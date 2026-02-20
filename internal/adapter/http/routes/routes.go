package routes

import (
	"log"
	_ "mecanica_xpto/docs" // This will be auto-generated
	"mecanica_xpto/internal/adapter/http/handlers"
	"mecanica_xpto/internal/adapter/http/middleware"
	repository2 "mecanica_xpto/internal/adapter/persistence/repository"
	"mecanica_xpto/internal/infrastructure/database"
	"mecanica_xpto/internal/usecase"
	"mecanica_xpto/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router = gin.Default()

const PORT = 8080

// Run will start the server
func Run() {
	setMiddlewares()

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	getRoutes()

	err := router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("Failed to startup the application: %v", err.Error())
	}
}

func getRoutes() {
	// Config JWT
	jwtCfg := utils.LoadJWTConfig()
	jwtService := utils.NewJWTService(jwtCfg)

	db := database.ConnectDatabase()
	userRepository := repository2.NewUserRepository(db)

	// Handler de autenticaÃƒÂ§ÃƒÂ£o
	authHandler := handlers.NewAuthHandler(
		usecase.NewAuthUseCase(jwtService, userRepository),
	)

	// Rotas pÃƒÂºblicas
	v1 := router.Group("/v1")
	v1.POST("/login", authHandler.Login)

	partsSupplyRepository := repository2.NewPartsSupplyRepository(db)
	partsSupplyUseCase := usecase.NewPartsSupplyUseCase(partsSupplyRepository)
	partsSupplyHandler := handlers.NewPartsSupplyHandler(partsSupplyUseCase)

	serviceRepository := repository2.NewServiceRepository(db)
	serviceUseCase := usecase.NewServiceUseCase(serviceRepository)
	serviceHandler := handlers.NewServiceHandler(serviceUseCase)

	vehiclesRepository := repository2.NewVehicleRepository(db)
	vehiclesUseCase := usecase.NewVehicleService(vehiclesRepository)
	vehicleHandler := handlers.NewVehicleHandler(vehiclesUseCase)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handlers.NewUserHandler(userUseCase)
	customerRepository := repository2.NewCustomerRepository(db)
	customerUseCase := usecase.NewCustomerUseCase(customerRepository, userRepository)
	customerHandler := handlers.NewCustomerHandler(customerUseCase)

	// Rotas protegidas
	authGroup := v1.Group("/")
	authGroup.Use(middleware.AuthMiddleware(jwtService))
	addPingRoutes(authGroup)
	addUserRoutes(authGroup, userHandler)
	addPartsSupplyRoutes(authGroup, partsSupplyHandler)
	addVehicleRoutes(authGroup, vehicleHandler)
	addServiceRoutes(authGroup, serviceHandler)
	addCustomerRoutes(authGroup, customerHandler)
}

func setMiddlewares() {

	middleware.SetTrustedProxies(router)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Recovered from panic: %v", recovered)
		c.AbortWithStatus(500)
	}))
}
