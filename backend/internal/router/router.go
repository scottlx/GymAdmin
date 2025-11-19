package router

import (
	"gym-admin/internal/controller"
	"gym-admin/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Initialize controllers
	authCtrl := controller.NewAuthController()
	userCtrl := controller.NewUserController()
	coachCtrl := controller.NewCoachController()

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("")
		{
			public.POST("/login", authCtrl.Login)
			public.POST("/register", authCtrl.Register)
		}

		// Protected routes
		auth := v1.Group("")
		auth.Use(middleware.Auth())
		{
			// User routes
			users := auth.Group("/users")
			{
				users.GET("", userCtrl.ListUsers)
				users.POST("", userCtrl.CreateUser)
				users.GET("/:id", userCtrl.GetUser)
				users.PUT("/:id", userCtrl.UpdateUser)
				users.DELETE("/:id", userCtrl.DeleteUser)
				users.GET("/:id/stats", userCtrl.GetUserStats)
			}

			// Membership card routes
			cards := auth.Group("/cards")
			{
				cards.GET("", nil)           // TODO: implement
				cards.POST("", nil)          // TODO: implement
				cards.GET("/:id", nil)       // TODO: implement
				cards.PUT("/:id", nil)       // TODO: implement
			}

			// Coach routes
			coaches := auth.Group("/coaches")
			{
				coaches.GET("", coachCtrl.ListCoaches)
				coaches.POST("", coachCtrl.CreateCoach)
				coaches.GET("/:id", coachCtrl.GetCoach)
				coaches.PUT("/:id", coachCtrl.UpdateCoach)
				coaches.DELETE("/:id", coachCtrl.DeleteCoach)
			}

			// Course routes
			courses := auth.Group("/courses")
			{
				courses.GET("", nil)         // TODO: implement
				courses.POST("", nil)        // TODO: implement
				courses.GET("/:id", nil)     // TODO: implement
				courses.PUT("/:id", nil)     // TODO: implement
				courses.DELETE("/:id", nil)  // TODO: implement
			}

			// Booking routes
			bookings := auth.Group("/bookings")
			{
				bookings.GET("", nil)        // TODO: implement
				bookings.POST("", nil)       // TODO: implement
				bookings.DELETE("/:id", nil) // TODO: implement
			}

			// Check-in routes
			checkins := auth.Group("/checkins")
			{
				checkins.GET("", nil)        // TODO: implement
				checkins.POST("", nil)       // TODO: implement
			}

			// Voucher routes
			vouchers := auth.Group("/vouchers")
			{
				vouchers.GET("", nil)        // TODO: implement
				vouchers.POST("/verify", nil) // TODO: implement
			}
		}
	}

	return r
}
