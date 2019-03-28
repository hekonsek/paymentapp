package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hekonsek/paymentapp/payments"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ApiServer struct {
	server *http.Server

	Store payments.PaymentStore

	Port int
}

func (a *ApiServer) Start() error {
	router := gin.Default()
	routes(a, router)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", a.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Problem starting HTTP endpoint: %s\n", err)
		}
	}()
	a.server = server
	return nil
}

func (a *ApiServer) Stop() error {
	return a.server.Shutdown(nil)
}

func (a *ApiServer) WaitForInterruptSignal() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal("Stopping server:", err)
	}
}
