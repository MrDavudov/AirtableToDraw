package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MrDavudov/AirtableToDraw/pkg/service"
	"github.com/joho/godotenv"
	"github.com/mehanizm/airtable"
	"github.com/spf13/viper"
	"github.com/studio-b12/gowebdav"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(ctx context.Context) error {
	// Инициализация config.yaml
	if err := initConfig(); err != nil {
		log.Fatalf("error Initializing configs: %s", err)
	}

	// Подключения .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err)
	}

	// Подключение к Airtable
	clientAT := airtable.NewClient(os.Getenv("API"))

	// Подключение к NextCloud
	clientNC := gowebdav.NewClient(os.Getenv("root"), 
									os.Getenv("username"), 
									os.Getenv("password"))

	// Взаимосвязи
	services := service.New(clientAT, clientNC)
	handlers := New(services)


	s.httpServer = &http.Server{
		Addr:           viper.GetString("port"),
		Handler:        handlers,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Запуск сервара
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()
	log.Println("Start server")
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err) 
	}

	longShutdown := make(chan struct{}, 1)

	go func() {
		time.Sleep(3 * time.Second)
		longShutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("Server shutdown: %w", ctx.Err())
	case <-longShutdown:
		log.Println("Stop server")
	}

	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}