package main

import (
	"context"
	"database/sql"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"time"
)

var _db *sql.DB
var _jwtIssuer *diploma.JWTIssuer
var _accClient *diploma.AccuralClient

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
	_accClient = diploma.NewAccuralClient(
		httpClient,
		_appConfig.accrualAddr,
	)
}

func main() {
	//todo startAccuralPolling()

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

func startAccuralPolling() {

	repo := diploma.NewOrderRepository(_db)

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		orderList, err := repo.FindByStatus(context.TODO(), diploma.StatusProcessing)
		if err != nil {
			log.Printf("error: %s", err)
			continue
		}

		for _, o := range orderList {
			acc, err := _accClient.Get(o.Id)
			if err != nil {
				log.Printf(err.Error())
			}
			accStatus, err := diploma.OrderStatusFromString(acc.Status)
			if err != nil {
				log.Printf(err.Error())
			}
			if o.Status == accStatus {
				continue
			}

			err = repo.UpdateStatus(context.TODO(), o.Id, o.Status, o.Accural)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	}
}
