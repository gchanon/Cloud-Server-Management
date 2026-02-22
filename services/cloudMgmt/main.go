package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golf/cloudmgmt/appUtility/config"
	"github.com/golf/cloudmgmt/services/cloudMgmt/behavior"
	"github.com/golf/cloudmgmt/services/cloudMgmt/handler"
	"github.com/golf/cloudmgmt/services/cloudMgmt/handler/middleware"
)

func main() {

	appConfig := config.LoadConfig()

	userBehavior := behavior.NewUserBehavior()
	auditBehavior := behavior.NewAuditTrailBehavior()
	serverBehavior := behavior.NewServerBehavior()

	userHandler := handler.NewUserHandler(userBehavior, appConfig.JWTSecret, appConfig.JWTExpireHour, appConfig.CookieDomain)

	serverHandler := handler.NewServerHandler(serverBehavior)

	if errGenSeed := userBehavior.GenSeedUser(); errGenSeed != nil {
		log.Fatal("Failed to generate seed data:", errGenSeed)
	}

	appHandler := fiber.New(fiber.Config{
		AppName: "Cloud Management Service",
	})

	appHandler.Use(recover.New())
	appHandler.Use(logger.New())
	appHandler.Use(cors.New(cors.Config{
		AllowOrigins:     appConfig.AllowedOrigin, // allow only browser from port 3000 responding to task 2 req.
		AllowCredentials: true,                    // allow cookies responding to task 2 req.
	}))

	// login api responding to task 2 req.
	authGroup := appHandler.Group("/auth")
	authGroup.Post("login", userHandler.Login)

	// server mgmt api responding to task 2 req.
	serverGroup := appHandler.Group("/servers", middleware.AuthMiddleware(appConfig.JWTSecret), middleware.AuditMiddleware(auditBehavior, serverBehavior))

	serverGroup.Get("", serverHandler.GetAllServer(appConfig))
	serverGroup.Post("", serverHandler.AddServer(appConfig))
	serverGroup.Post("/:serverId/power", serverHandler.PowerControlServer(appConfig))

	log.Printf("User Service starting on port %s", appConfig.ServicePort)
	log.Fatal(appHandler.Listen(":" + appConfig.ServicePort))

}
