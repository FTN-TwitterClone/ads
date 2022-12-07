package main

import (
	"context"
	"github.com/FTN-TwitterClone/ads/controller"
	"github.com/FTN-TwitterClone/ads/controller/jwt"
	"github.com/FTN-TwitterClone/ads/repository/cassandra"
	"github.com/FTN-TwitterClone/ads/repository/mongo"
	"github.com/FTN-TwitterClone/ads/service"
	"github.com/FTN-TwitterClone/ads/tls"
	"github.com/FTN-TwitterClone/ads/tracing"
	"github.com/FTN-TwitterClone/grpc-stubs/proto/ads"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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
	tracer := tp.Tracer("ads")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	eventsRepository, err := cassandra.NewCassandraEventsRepository(tracer)
	if err != nil {
		log.Fatal(err)
	}

	reportsRepository, err := mongo.NewMongoReportsRepository(tracer)
	if err != nil {
		log.Fatal(err)
	}

	adsService := service.NewAdsService(eventsRepository, reportsRepository, tracer)

	adsController := controller.NewAdsController(adsService, tracer)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(
		tracing.ExtractTraceInfoMiddleware,
		jwt.ExtractJWTUserMiddleware(tracer),
	)

	router.HandleFunc("/users/{tweetId}/{username}/visit/", adsController.AddProfileVisitedEvent).Methods("POST")
	router.HandleFunc("/ads/{id}/reports/{year}/{month}/", adsController.GetMonthlyReport).Methods("GET")
	router.HandleFunc("/ads/{id}/reports/{year}/{month}/{day}/", adsController.GetDailyReport).Methods("GET")

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

	lis, err := net.Listen("tcp", "0.0.0.0:9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds := credentials.NewTLS(tls.GetgRPCServerTLSConfig())

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	ads.RegisterAdsServiceServer(grpcServer, service.NewgRPCAdsService(tracer, eventsRepository, reportsRepository))
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		return
	}

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
