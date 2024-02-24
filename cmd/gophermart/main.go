package main

import (
	"database/sql"
	"github.com/gam6itko/go-musthave-diploma/internal"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var _db *sql.DB
var _jwtIssuer *internal.JWTIssuer

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	_jwtIssuer = internal.NewJWTIssuer(_appConfig.jwtKey)
}

func main() {
	server := &http.Server{
		Addr:    _appConfig.listenAdd,
		Handler: newRouter(),
	}

	log.Printf("Starting server on %s", _appConfig.listenAdd)
	if err := server.ListenAndServe(); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		log.Printf(err.Error())
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

	r.Post("/api/user/register", postUserRegister)
	r.Post("/api/user/login", postUserLogin)
	r.Post("/api/user/orders", postUserOrders)
	r.Get("/api/user/orders", getUserOrders)
	r.Get("/api/user/balance", getUserBalance)
	r.Post("/api/user/balance/withdraw", postUserBalanceWithdraw)
	r.Get("/api/user/withdrawals", getUserWithdrawals)

	return r
}
