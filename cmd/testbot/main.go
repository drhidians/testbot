package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drhidians/testbot/middleware/jwtauth"
	"github.com/drhidians/testbot/server"

	"github.com/go-kit/kit/log"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	_userRepo "github.com/drhidians/testbot/user/repository"
	_userUcase "github.com/drhidians/testbot/user/usecase"

	_botRepo "github.com/drhidians/testbot/bot/repository"
	_botUcase "github.com/drhidians/testbot/bot/usecase"

	tg "github.com/drhidians/testbot"
)

func init() {
	viper.AutomaticEnv()
	pflag.String("db", viper.GetString("TESTBOT_db"), "строка для подключения к базе данных (postgres://user:pass@host:port/dbname)")
	pflag.Int("db-max-open-conns", viper.GetInt("TESTBOT_db-max-open-conns"), "максимальный размер пула подключений к БД")
	pflag.Int("db-max-idle-conns", viper.GetInt("TESTBOT_db-max-idle-conns"), "максимальное кол-во простаювающих соеденений к БД")
	pflag.String("secret", viper.GetString("TESTBOT_secret"), "секрет для подписи JWT-токенов")
	pflag.String("bot-token", viper.GetString("TESTBOT_bot-token"), "токен бота полученный у @BotFather")
	pflag.Int("bot-webhook-max-conns", viper.GetInt("TESTBOT_bot-webhook-max-conns"), "максимальное количество параллельных HTTP-соединении от Telegram-сервера для бота")
	pflag.String("addr", viper.GetString("TESTBOT_addr"), "адресс на котром будет запущен сервер (:8000, localhost:8000, ...)")
	pflag.String("domain", viper.GetString("TESTBOT_domain"), "домен на который будет установлен вебхук и который нужно использовать при построений абсолютных URL.")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetDefault("context.timeout", 15)
}

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	db, err := sql.Open("postgres", viper.GetString("db")+"?sslmode=disable")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(viper.GetInt("db-max-idle-conns"))
	db.SetMaxOpenConns(viper.GetInt("db-max-open-conns"))

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	b, err := tg.NewBot(ctx, viper.GetString("bot-token"))
	if err != nil {
		panic(err)
	}

	w, err := b.GetWebhookInfo(ctx)

	if err != nil {
		panic(err)
	}

	domainURL := "https://" + viper.GetString("domain") + "/"
	webhookURL := domainURL + "bot/webhook"
	if w.URL != webhookURL {
		webhookConfig := b.NewWebhook(webhookURL, viper.GetInt("bot-webhook-max-conns"))
		err = b.SetWebhook(ctx, webhookConfig)
		if err != nil {
			panic(err)
		}
	}

	tokenAuth := jwtauth.New("HS256", []byte(viper.GetString("secret")), nil)

	userRepo := _userRepo.NewPostgresUserRepository(db)
	botRepo := _botRepo.NewBotRepository(b, tokenAuth, domainURL)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	bs := _botUcase.NewBotService(userRepo, botRepo, timeoutContext)
	bs = _botUcase.NewLoggingService(logger, bs)

	us := _userUcase.NewUserService(userRepo, botRepo, timeoutContext)
	us = _userUcase.NewLoggingService(logger, us)

	srv := server.New(bs, us, log.With(logger, "component", "http"), tokenAuth)

	errs := make(chan error, 2)

	go func() {
		logger.Log("transport", "http", "address", viper.GetString("addr"), "msg", "listening")
		errs <- http.ListenAndServe(viper.GetString("addr"), srv)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
