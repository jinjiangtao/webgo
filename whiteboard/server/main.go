package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"whiteboard/internal/database"
	"whiteboard/internal/handler"
	"whiteboard/internal/middleware"
	"whiteboard/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	hub := ws.NewHub()
	go hub.Run()

	whiteboardRepo := database.NewWhiteboardRepo()
	operationRepo := database.NewOperationRepo()
	snapshotRepo := database.NewSnapshotRepo()

	whiteboardHandler := handler.NewWhiteboardHandler(whiteboardRepo, operationRepo, snapshotRepo)
	wsHandler := handler.NewWebSocketHandler(hub)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		whiteboards := api.Group("/whiteboards")
		{
			whiteboards.POST("", whiteboardHandler.CreateWhiteboard)
			whiteboards.GET("", whiteboardHandler.ListWhiteboards)
			whiteboards.GET("/:id", whiteboardHandler.GetWhiteboard)
			whiteboards.PUT("/:id", whiteboardHandler.UpdateWhiteboard)
			whiteboards.DELETE("/:id", whiteboardHandler.DeleteWhiteboard)
			whiteboards.POST("/:id/operations", whiteboardHandler.SaveOperations)
			whiteboards.POST("/:id/snapshots", whiteboardHandler.CreateSnapshot)
			whiteboards.GET("/:id/snapshots", whiteboardHandler.ListSnapshots)
			whiteboards.POST("/:id/clear", whiteboardHandler.ClearWhiteboard)
		}
	}

	r.GET("/ws", wsHandler.WebSocketHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Printf("server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}

	log.Println("server exiting")
}
