package main

import (
	"auth/internal/config"
	mylog "auth/internal/lib/log"
	"auth/internal/service"
	"auth/internal/service/handlers"
	"auth/internal/storage/postgres"
	"auth/internal/storage/redis"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/julienschmidt/httprouter"
)

//NOTE: CONFIG_PATH=./config/local_config.yaml go run ./cmd/main.go
func main() {
	//TODO: Инициализировать конфиг
	cnf := config.LoadConfig();
	//TODO: Инициализировать логгер
	log := mylog.LoadLogger(cnf.Env);

	log.Debug("Init logger", slog.String("env", cnf.Env));

	//TODO: Инициализировать БД POSTGRES
	storage, err := postgres.New(cnf.Postgres);
	if (err != nil) {
		log.Error("Error connect POSTGRES", slog.String("error", err.Error()));
		os.Exit(1);
	}
	log.Debug("Init storage");

	//TODO: Инициализировать БД REDIS
	cache, err := redis.New(cnf.Redis);
	if (err != nil) {
		log.Error("Error connect Redis", slog.String("error", err.Error()));
		os.Exit(1);
	}
	log.Debug("Init cache");

	log.Debug("TTL", 
		slog.String("time1", cnf.Server.TokenTTL.String()),
		slog.String("time2", cnf.Server.Timeout.String()),
		slog.String("time3", cnf.Server.IdleTimeout.String()),
		slog.String("time4", cnf.Postgres.MaxIdleTime.String()));

	//TODO: Инициализировать приложение (будет содержать все предыдущее)
	app := service.New(cnf, storage, storage, cache, log);

	//TODO: Инициализировать router
	router := httprouter.New();

	//TODO: Настроить маршруты
	router.HandlerFunc(http.MethodPost, "/Signup", handlers.SignupUserHandler(app));
	router.HandlerFunc(http.MethodPost, "/Register", handlers.VerifyOTPHandler(app));
	router.HandlerFunc(http.MethodPost, "/Login", handlers.LoginHandlers(app));
	router.HandlerFunc(http.MethodPost, "/Verify", handlers.VerifyTokenHandlers(app));

	//TODO: Запустить сервер

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port),
		Handler: handlers.RecoverPanic(app, router),
		IdleTimeout: cnf.Server.IdleTimeout,
		ReadTimeout: cnf.Server.Timeout,
		WriteTimeout: cnf.Server.Timeout,
	};

	// для сигналов остановки сервера
	done := make(chan os.Signal, 1);
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM);

	go func(){
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", 
					slog.String("error:", err.Error()));
		}
	} ()

	log.Info("server started");

	//TODO: Закрыть сервер при сигналах
	<-done //примем значение, если будет сигнал
	log.Info("stopping server");

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second);
	defer cancel();

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server");
		return;
	}
}