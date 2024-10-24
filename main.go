package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"service/cache"
	"service/config"
	"service/db"
	_ "service/docs"
	"service/handlers"
	"service/middleware"
)

// @title Referral System API
// @version 1.0
// @description API для системы реферальных ссылок
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	config.InitConfig()

	cache.InitCache()

	db.InitDB()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Руты
	auth := r.Group(config.C.AuthPath)
	{
		auth.POST(config.C.AuthRegisterPath, handlers.Register)
		auth.POST(config.C.AuthLoginPath, handlers.Login)
	}

	referral := r.Group(config.C.ReferralPath)
	referral.Use(middleware.AuthMiddleware())
	{
		referral.POST(config.C.ReferralCreatePath, handlers.CreateReferralCode)
		referral.DELETE(config.C.ReferralDeletePath, handlers.DeleteReferralCode)
		referral.GET(config.C.ReferralGetPath, handlers.GetReferralCodeByEmail)
		referral.POST(config.C.ReferralRegisterPath, handlers.RegisterWithReferralCode)
		referral.GET(config.C.ReferralReferralsPath, handlers.GetReferralsByReferrerID)
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
