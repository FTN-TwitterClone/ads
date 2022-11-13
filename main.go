package main

import (
	"ads/controller"
	"ads/controller/jwt"
	"ads/repository/cassandra"
	"ads/service"
	"ads/tls"
	"ads/tracing"
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	ctx := context.Background()
	exp, tracingErr := tracing.NewExporter()
	if tracingErr != nil {
		log.Fatalf("failed to initialize exporter: %v", tracingErr)
	}
	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := tracing.NewTraceProvider(exp)
	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	// Finally, set the tracer that can be used for this package.
	tracer := tp.Tracer("tweet")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	adsRepository, err := cassandra.NewCassandraAdsRepository(tracer)
	if err != nil {
		log.Fatal(err)
	}

	adsService := service.NewAdsService(adsRepository, tracer)

	adsController := controller.NewAdsController(adsService, tracer)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(
		tracing.ExtractTraceInfoMiddleware,
		jwt.ExtractJWTUserMiddleware(tracer),
	)

	router.HandleFunc("/first", adsController.DeleteLike).Methods("GET")
	router.HandleFunc("/{id}/next", adsController.DeleteLike).Methods("GET")

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	// start server
	srv := &http.Server{
		Addr:      "0.0.0.0:8000",
		Handler:   handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(router),
		TLSConfig: tls.GetHTTPServerTLSConfig(),
	}

	go func() {
		log.Println("server starting")

		certFile := os.Getenv("CERT")
		keyFile := os.Getenv("KEY")

		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
