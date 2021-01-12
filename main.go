package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pallat/todos/auth"
	"github.com/pallat/todos/captcha"
	"github.com/pallat/todos/logger"
	"github.com/pallat/todos/todos"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// use environment variables
	viper.AutomaticEnv()
	// replace "." with "_" for environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// use config file
	viper.SetConfigFile("config.yaml")

	l, _ := zap.NewProduction()
	defer l.Sync()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning:", err)
	}

	viper.SetDefault("app.addr", "0.0.0.0:8888")

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// set connection pool
	sqlDB, sqlErr := db.DB()
	if sqlErr != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// db, err := sql.Open("sqlite", "gorm.db")
	// if err != nil {
	// 	// This will not be a connection error, but a DSN parse error or
	// 	// another initialization error.
	// 	log.Fatal(err)
	// }
	// db.SetConnMaxLifetime(0)
	// db.SetMaxIdleConns(50)
	// db.SetMaxOpenConns(50)

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(logger.Middleware(l))

	router.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	router.GET("/captcha", captchaHandler)
	router.POST("/exchange", exchangeHandler)

	router.GET("/todos", todos.NewListTodoHandler(db))
	router.GET("/todos/:id", todos.NewGetTodoHandler(db))
	router.POST("/todos", todos.NewNewTodoHandler(db))
	router.PUT("/todos/:id", todos.NewUpdateTodoHandler(db))
	router.DELETE("/todos/:id", todos.NewDeleteTodoHandler(db))

	srv := &http.Server{
		Addr:         viper.GetString("app.addr"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("listing at", viper.GetString(("app.addr")))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func captchaHandler(c echo.Context) error {
	key, captcha := captcha.KeyQuestion()
	return c.JSON(http.StatusOK, map[string]string{
		"key":     key,
		"captcha": captcha,
	})
}

func exchangeHandler(c echo.Context) error {
	var ans struct {
		Key    string `json:"key"`
		Answer int    `json:"answer"`
	}

	if err := c.Bind(&ans); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if !captcha.Answer(ans.Key, ans.Answer) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "wrong answer",
		})
	}

	t, err := auth.Token()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
