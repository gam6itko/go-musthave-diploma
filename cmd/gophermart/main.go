package main

import (
	"database/sql"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
)

var _db *sql.DB
var _jwtIssuer *diploma.JWTIssuer
var _accClient *diploma.AccrualClient

func init() {

	var err error
	_db, err = sql.Open("pgx", _appConfig.dbDsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatal(err)
	}

	_jwtIssuer = diploma.NewJWTIssuer(_appConfig.jwtKey)

	httpClient := &http.Client{}
	_accClient = diploma.NewAccrualClient(
		httpClient,
		_appConfig.accrualAddr,
	)

	if err := initDatabaseSchema(_db); err != nil {
		log.Fatalf("schema init error. %s", err)
	}

}

func main() {
	server := &http.Server{
		Addr:    _appConfig.listenAdd,
		Handler: newRouter(),
	}

	log.Printf("Starting server on %s", _appConfig.listenAdd)
	if err := server.ListenAndServe(); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		log.Printf("http server error. %s", err)
	}
}

func newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL)
			h.ServeHTTP(w, r)
		})
	})

	r.Use(compressMiddleware)

	// ниже попытка поиграться в DI и тестирование
	r.Post("/api/user/register", func(w http.ResponseWriter, r *http.Request) {
		postUserRegister(w, r, _db, _jwtIssuer)
	})
	r.Post("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		postUserLogin(w, r, _db, _jwtIssuer)
	})
	r.Post("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		postUserOrders(w, r, _db, _accClient)
	})
	r.Get("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		getUserOrders(w, r, _db)
	})
	r.Get("/api/user/balance", func(w http.ResponseWriter, r *http.Request) {
		getUserBalance(w, r, _db)
	})
	r.Post("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
		postUserBalanceWithdraw(w, r, _db)
	})
	r.Get("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {
		getUserWithdrawals(w, r, _db)
	})

	return r
}
