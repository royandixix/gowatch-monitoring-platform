package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/royandi/gowatch/backend/internal/config"
	"github.com/royandi/gowatch/backend/internal/handler"
	"github.com/royandi/gowatch/backend/internal/middleware"
	"github.com/royandi/gowatch/backend/internal/repository"
	"github.com/royandi/gowatch/backend/internal/service"
	"github.com/royandi/gowatch/backend/internal/worker"
)

func main() {
	db, err :=
		config.ConnectDatabase()

	if err != nil {
		log.Fatal(
			"❌ Database error: ",
			err,
		)
	}

	defer db.Close()

	fmt.Println(
		"✅ PostgreSQL berhasil terhubung",
	)

	jwtSecret :=
		os.Getenv(
			"JWT_SECRET",
		)

	if jwtSecret == "" {
		log.Fatal(
			"❌ JWT_SECRET belum dikonfigurasi",
		)
	}

	frontendURL :=
		os.Getenv(
			"FRONTEND_URL",
		)

	if frontendURL == "" {
		frontendURL =
			"http://localhost:5173"
	}

	userRepository :=
		repository.NewUserRepository(
			db,
		)

	monitorRepository :=
		repository.NewMonitorRepository(
			db,
		)

	monitorResultRepository :=
		repository.NewMonitorResultRepository(
			db,
		)

	dashboardRepository :=
		repository.NewDashboardRepository(
			db,
		)

	authService :=
		service.NewAuthService(
			userRepository,
			jwtSecret,
		)

	monitorService :=
		service.NewMonitorService(
			monitorRepository,
		)

	profileService :=
		service.NewProfileService(
			userRepository,
		)

	authHandler :=
		handler.NewAuthHandler(
			authService,
		)

	monitorHandler :=
		handler.NewMonitorHandler(
			monitorService,
		)

	monitorResultHandler :=
		handler.NewMonitorResultHandler(
			monitorResultRepository,
		)

	profileHandler :=
		handler.NewProfileHandler(
			profileService,
		)

	dashboardHandler :=
		handler.NewDashboardHandler(
			dashboardRepository,
		)

	monitorWorker :=
		worker.NewMonitorWorker(
			monitorRepository,
			monitorResultRepository,
		)

	appContext, stop :=
		signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)

	defer stop()

	workerDone :=
		make(
			chan struct{},
		)

	go func() {
		defer close(
			workerDone,
		)

		monitorWorker.Start(
			appContext,
		)
	}()

	router :=
		gin.Default()

	// Jangan percaya X-Forwarded-For dari sembarang client.
	// Saat deploy dengan reverse proxy nanti,
	// proxy spesifik akan kita konfigurasi.
	err =
		router.SetTrustedProxies(
			nil,
		)

	if err != nil {
		log.Fatal(
			"❌ Gagal mengatur trusted proxy: ",
			err,
		)
	}

	router.Use(
		middleware.CORSMiddleware(
			frontendURL,
		),
	)

	router.GET(
		"/",
		handler.Home,
	)

	api :=
		router.Group(
			"/api/v1",
		)

	{
		api.GET(
			"/health",
			handler.Health,
		)

		authRoutes :=
			api.Group(
				"/auth",
			)

		{
			authRoutes.POST(
				"/register",
				middleware.RateLimit(
					5,
					10*time.Minute,
				),
				authHandler.Register,
			)

			authRoutes.POST(
				"/login",
				middleware.RateLimit(
					10,
					time.Minute,
				),
				authHandler.Login,
			)

			authRoutes.GET(
				"/me",
				middleware.AuthMiddleware(
					jwtSecret,
				),
				authHandler.Me,
			)
		}

		protected :=
			api.Group("")

		protected.Use(
			middleware.AuthMiddleware(
				jwtSecret,
			),
		)

		{
			protected.GET(
				"/dashboard",
				dashboardHandler.Show,
			)

			protected.GET(
				"/check",
				middleware.RateLimit(
					30,
					time.Minute,
				),
				handler.Check,
			)

			protected.GET(
				"/check-multiple",
				middleware.RateLimit(
					10,
					time.Minute,
				),
				handler.CheckMultiple,
			)

			protected.POST(
				"/monitors",
				monitorHandler.Create,
			)

			protected.GET(
				"/monitors",
				monitorHandler.List,
			)

			protected.GET(
				"/monitors/:id",
				monitorHandler.GetByID,
			)

			protected.PUT(
				"/monitors/:id",
				monitorHandler.Update,
			)

			protected.DELETE(
				"/monitors/:id",
				monitorHandler.Delete,
			)

			protected.GET(
				"/monitors/:id/results",
				monitorResultHandler.ByMonitor,
			)

			protected.GET(
				"/monitors/:id/stats",
				monitorResultHandler.Stats,
			)

			protected.GET(
				"/history",
				monitorResultHandler.History,
			)

			protected.PUT(
				"/profile",
				profileHandler.Update,
			)

			protected.PUT(
				"/profile/password",
				profileHandler.ChangePassword,
			)
		}
	}

	port :=
		os.Getenv(
			"APP_PORT",
		)

	if port == "" {
		port =
			"8081"
	}

	server :=
		&http.Server{
			Addr: ":" + port,

			Handler: router,

			ReadHeaderTimeout: 5 * time.Second,

			ReadTimeout: 15 * time.Second,

			WriteTimeout: 30 * time.Second,

			IdleTimeout: 60 * time.Second,
		}

	fmt.Println(
		"======================================",
	)

	fmt.Println(
		"🚀 GoWatch API - Gin",
	)

	fmt.Println(
		"🐘 PostgreSQL connected",
	)

	fmt.Println(
		"🔐 JWT Authentication enabled",
	)

	fmt.Println(
		"🛡️ SSRF Protection enabled",
	)

	fmt.Println(
		"⏱️ HTTP Timeout enabled",
	)

	fmt.Println(
		"🚦 Rate Limiting enabled",
	)

	fmt.Println(
		"📡 Monitor CRUD enabled",
	)

	fmt.Println(
		"📈 Monitor Statistics enabled",
	)

	fmt.Println(
		"📊 Dashboard Analytics enabled",
	)

	fmt.Println(
		"📚 Monitoring History enabled",
	)

	fmt.Println(
		"👤 Profile Management enabled",
	)

	fmt.Println(
		"⚙️ Background Monitoring enabled",
	)

	fmt.Println(
		"🌐 Backend: http://localhost:" +
			port,
	)

	fmt.Println(
		"🖥️ Frontend: " +
			frontendURL,
	)

	fmt.Println(
		"======================================",
	)

	serverError :=
		make(
			chan error,
			1,
		)

	go func() {
		serverError <- server.ListenAndServe()
	}()

	select {
	case <-appContext.Done():
		log.Println(
			"🛑 Shutdown signal diterima",
		)

	case err :=
		<-serverError:

		if err != nil &&
			!errors.Is(
				err,
				http.ErrServerClosed,
			) {
			log.Printf(
				"❌ HTTP server berhenti: %v",
				err,
			)
		}

		stop()
	}

	shutdownContext,
		shutdownCancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

	defer shutdownCancel()

	log.Println(
		"⏳ Menghentikan HTTP server...",
	)

	if err :=
		server.Shutdown(
			shutdownContext,
		); err != nil {
		log.Printf(
			"❌ Graceful shutdown gagal: %v",
			err,
		)
	} else {
		log.Println(
			"✅ HTTP server berhenti dengan aman",
		)
	}

	select {
	case <-workerDone:
		log.Println(
			"✅ Background worker selesai",
		)

	case <-time.After(
		5 * time.Second,
	):
		log.Println(
			"⚠️ Timeout menunggu background worker",
		)
	}

	log.Println(
		"👋 GoWatch backend selesai",
	)
}
