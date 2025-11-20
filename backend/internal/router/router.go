package router

import (
	"gym-admin/internal/controller"
	"gym-admin/internal/middleware"
	"gym-admin/internal/scheduler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Initialize scheduler
	cardExpiryScheduler := scheduler.NewCardExpiryScheduler()
	cardExpiryScheduler.Start()
	performanceScheduler := scheduler.NewPerformanceScheduler()
	performanceScheduler.Start()

	// Initialize controllers
	authCtrl := controller.NewAuthController()
	userCtrl := controller.NewUserController()
	coachCtrl := controller.NewCoachController()
	cardCtrl := controller.NewCardController()
	cardTypeCtrl := controller.NewCardTypeController()
	courseCtrl := controller.NewCourseController()
	checkinCtrl := controller.NewCheckInController()
	notificationCtrl := controller.NewNotificationController()
	adminCtrl := controller.NewAdminController(cardExpiryScheduler)
	faceDeviceCtrl := controller.NewFaceDeviceController()
	userFaceCtrl := controller.NewUserFaceController()
	voucherCtrl := controller.NewVoucherController()
	coachCertCtrl := controller.NewCoachCertificationController()
	coachPerfCtrl := controller.NewCoachPerformanceController()

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

				// User status management
				users.POST("/:id/freeze", userCtrl.FreezeUser)
				users.POST("/:id/unfreeze", userCtrl.UnfreezeUser)
				users.POST("/:id/blacklist", userCtrl.AddToBlacklist)
				users.DELETE("/:id/blacklist", userCtrl.RemoveFromBlacklist)
				users.GET("/:id/status-logs", userCtrl.GetStatusLogs)

				// Batch operations
				users.POST("/batch/freeze", userCtrl.BatchFreezeUsers)
				users.POST("/batch/unfreeze", userCtrl.BatchUnfreezeUsers)

				// Status summary
				users.GET("/status/summary", userCtrl.GetUserStatusSummary)
			}

			// Card type routes
			cardTypes := auth.Group("/card-types")
			{
				cardTypes.GET("", cardTypeCtrl.ListCardTypes)
				cardTypes.POST("", cardTypeCtrl.CreateCardType)
				cardTypes.GET("/:id", cardTypeCtrl.GetCardType)
				cardTypes.PUT("/:id", cardTypeCtrl.UpdateCardType)
				cardTypes.DELETE("/:id", cardTypeCtrl.DeleteCardType)
				cardTypes.POST("/:id/enable", cardTypeCtrl.EnableCardType)
				cardTypes.POST("/:id/disable", cardTypeCtrl.DisableCardType)
				cardTypes.POST("/sort-order", cardTypeCtrl.UpdateSortOrder)
			}

			// Membership card routes
			cards := auth.Group("/cards")
			{
				cards.GET("", cardCtrl.ListCards)
				cards.POST("", cardCtrl.CreateCard)
				cards.GET("/:id", cardCtrl.GetCard)
				cards.PUT("/:id", cardCtrl.UpdateCard)
				cards.DELETE("/:id", cardCtrl.DeleteCard)

				// Card operations
				cards.POST("/:id/renew", cardCtrl.RenewCard)
				cards.POST("/:id/freeze", cardCtrl.FreezeCard)
				cards.POST("/:id/unfreeze", cardCtrl.UnfreezeCard)
				cards.POST("/:id/transfer", cardCtrl.TransferCard)
				cards.GET("/:id/operations", cardCtrl.GetCardOperations)
			}

			// Coach routes
			coaches := auth.Group("/coaches")
			{
				coaches.GET("", coachCtrl.ListCoaches)
				coaches.POST("", coachCtrl.CreateCoach)
				coaches.GET("/:id", coachCtrl.GetCoach)
				coaches.PUT("/:id", coachCtrl.UpdateCoach)
				coaches.DELETE("/:id", coachCtrl.DeleteCoach)

				// Coach performance routes
				coaches.GET("/:id/performance", coachPerfCtrl.GetCoachPerformance)
				coaches.POST("/:id/performance/update", coachPerfCtrl.UpdateCoachPerformance)

				// Coach certification routes
				certifications := coaches.Group("/:id/certifications")
				{
					certifications.GET("", coachCertCtrl.GetCoachCertifications)
					certifications.POST("", coachCertCtrl.CreateCertification)
				}
			}

			// Coach rankings route
			auth.GET("/coaches/rankings", coachPerfCtrl.GetCoachRankings)

			// Certification management routes
			certs := auth.Group("/certifications")
			{
				certs.GET("", coachCertCtrl.ListCertifications)
				certs.GET("/:id", coachCertCtrl.GetCertification)
				certs.PUT("/:id", coachCertCtrl.UpdateCertification)
				certs.DELETE("/:id", coachCertCtrl.DeleteCertification)
				certs.POST("/:id/approve", coachCertCtrl.ApproveCertification)
				certs.POST("/:id/reject", coachCertCtrl.RejectCertification)
			}

			// Course routes
			courses := auth.Group("/courses")
			{
				courses.GET("", courseCtrl.ListCourses)
				courses.POST("", courseCtrl.CreateCourse)
				courses.GET("/:id", courseCtrl.GetCourse)
				courses.PUT("/:id", courseCtrl.UpdateCourse)
				courses.DELETE("/:id", courseCtrl.DeleteCourse)
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
				checkins.GET("", checkinCtrl.ListCheckIns)
				checkins.POST("", checkinCtrl.CheckIn)
			}

			// User check-in and statistics routes
			users.GET("/:id/checkin/today", checkinCtrl.GetTodayCheckIn)
			users.GET("/:id/checkin/stats", checkinCtrl.GetUserStats)
			users.GET("/:id/checkin/stats/detailed", checkinCtrl.GetDetailedStats)
			users.GET("/:id/checkin/stats/calendar", checkinCtrl.GetCheckInCalendar)
			users.POST("/:id/checkin/stats/recalculate", checkinCtrl.RecalculateStats)

			// Notification routes
			notifications := auth.Group("/notifications")
			{
				notifications.GET("", notificationCtrl.ListNotifications)
				notifications.GET("/:id", notificationCtrl.GetNotification)
				notifications.POST("/:id/read", notificationCtrl.MarkAsRead)
				notifications.POST("/read-all", notificationCtrl.MarkAllAsRead)
				notifications.DELETE("/:id", notificationCtrl.DeleteNotification)
				notifications.GET("/unread-count", notificationCtrl.GetUnreadCount)
			}

			// Admin routes
			admin := auth.Group("/admin")
			{
				admin.POST("/trigger-expiry-check", adminCtrl.TriggerExpiryCheck)
				admin.GET("/scheduler-status", adminCtrl.GetSchedulerStatus)
				admin.GET("/expiring-cards", adminCtrl.GetExpiringCards)
				admin.GET("/expired-cards", adminCtrl.GetExpiredCards)
			}

			// Voucher routes
			vouchers := auth.Group("/vouchers")
			{
				vouchers.GET("", voucherCtrl.ListVouchers)
				vouchers.POST("/verify", voucherCtrl.VerifyVoucher)
				vouchers.GET("/:id", voucherCtrl.GetVoucher)
			}

			// Face device routes
			faceDevices := auth.Group("/face-devices")
			{
				faceDevices.GET("", faceDeviceCtrl.ListDevices)
				faceDevices.POST("", faceDeviceCtrl.CreateDevice)
				faceDevices.GET("/:id", faceDeviceCtrl.GetDevice)
				faceDevices.PUT("/:id", faceDeviceCtrl.UpdateDevice)
				faceDevices.DELETE("/:id", faceDeviceCtrl.DeleteDevice)

				// Device operations
				faceDevices.POST("/:id/enable", faceDeviceCtrl.EnableDevice)
				faceDevices.POST("/:id/disable", faceDeviceCtrl.DisableDevice)
				faceDevices.POST("/:id/heartbeat", faceDeviceCtrl.UpdateDeviceHeartbeat)

				// Device statistics
				faceDevices.GET("/status/summary", faceDeviceCtrl.GetDeviceStatusSummary)
			}

			// User face routes
			faces := auth.Group("/faces")
			{
				faces.GET("", userFaceCtrl.ListFaces)
				faces.POST("", userFaceCtrl.RegisterFace)
				faces.GET("/:id", userFaceCtrl.GetFace)
				faces.PUT("/:id", userFaceCtrl.UpdateFace)
				faces.DELETE("/:id", userFaceCtrl.DeleteFace)

				// User-specific face operations
				faces.GET("/user/:user_id", userFaceCtrl.GetUserFaces)
				faces.GET("/user/:user_id/main", userFaceCtrl.GetMainFace)
				faces.POST("/user/:user_id/main/:face_id", userFaceCtrl.SetMainFace)
				faces.DELETE("/user/:user_id/all", userFaceCtrl.DeleteUserFaces)

				// Face statistics
				faces.GET("/statistics", userFaceCtrl.GetFaceStatistics)
			}
		}
	}

	return r
}
