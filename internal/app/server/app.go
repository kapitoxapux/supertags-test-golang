package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"supertags/internal/app/config"
	"supertags/internal/app/handler"
	"supertags/internal/app/service"
	mysql "supertags/internal/app/storage"

	"github.com/go-chi/chi"
)

type App struct {
	httpServer *http.Server
	service    *service.Service
}

func NewApp() *App {

	db := GetDB()
	service := service.NewService(db)

	return &App{
		service: service,
	}
}

func GetDB() service.Storage {
	config.SetConfig()
	if status, _ := handler.ConnectionDBCheck(); status == http.StatusOK {

		return mysql.NewDB()
	}

	return nil
}

func registerHTTPEndpoints(router *chi.Mux, service service.Service) {

	h := handler.NewHandler(service)

	router.Get("/login", h.GetAuthentication)
	router.Get("/events", h.GetArangoData)

}

func (a *App) Run(ctx context.Context) error {
	route := chi.NewRouter()
	address := config.GetAddress()
	registerHTTPEndpoints(route, *a.service)

	a.httpServer = &http.Server{
		Addr:    address,
		Handler: handler.CustomMiddleware(route),
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}

	}()

	<-ctx.Done()

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	quit := make(chan struct{}, 1)

	go func() {
		time.Sleep(3 * time.Second)
		quit <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("server shutdown: %w", ctx.Err())
	case <-quit:
		log.Println("finished")
	}

	return nil
}
