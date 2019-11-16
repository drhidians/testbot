package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drhidians/testbot/server"

	"github.com/go-kit/kit/log"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	_userRepo "github.com/drhidians/testbot/user/repository"
	_userUcase "github.com/drhidians/testbot/user/usecase"
)

func init() {
	viper.AutomaticEnv()

	pflag.String("db", viper.GetString("TESTBOT_db"), "строка для подключения к базе данных (postgres://user:pass@host:port/dbname)")
	pflag.Int("db-max-open-conns", viper.GetInt("TESTBOT_db-max-open-conns"), "максимальный размер пула подключений к БД")
	pflag.Int("db-max-idle-conns", viper.GetInt("TESTBOT_db-max-idle-conns"), "максимальное кол-во простаювающих соеденений к БД")
	pflag.String("secret", viper.GetString("TESTBOT_secret"), "секрет для подписи JWT-токенов")
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

	db, err := sql.Open("postgres", viper.GetString("db"))
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

	userRepo := _userRepo.NewPostgresUserRepository(db)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	uu := _userUcase.NewUserUseCase(userRepo, botRepo, timeoutContext)

	//TO DO . Use of context for test purposes. DELETE AFTER

	/*authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)


	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDeliver.NewArticleHandler(e, au)
	*/
	srv := server.New(bs, ts, hs, log.With(logger, "component", "http"))
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
