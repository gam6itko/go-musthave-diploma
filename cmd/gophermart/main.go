package main

import (
	"database/sql"
	"github.com/gam6itko/go-musthave-diploma/internal/accrual"
	"github.com/gam6itko/go-musthave-diploma/internal/controller"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	repository2 "github.com/gam6itko/go-musthave-diploma/internal/repository"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
)

var _db *sql.DB
var _jwtIssuer *diploma.JWTIssuer
var _accClient *accrual.Client

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
	_accClient = accrual.NewAccrualClient(
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

	userRepo := repository2.NewUserRepository(_db)
	orderRepo := repository2.NewOrderRepository(_db)
	wRepoRepo := repository2.NewWithdrawalRepository(_db)
	// controllers
	anonController := controller.NewAnonController(_jwtIssuer, userRepo)
	userController := controller.NewUserController(userRepo)
	orderController := controller.NewOrderController(_accClient, orderRepo)
	wController := controller.NewWithdrawalController(wRepoRepo, userRepo)
	// ниже попытка поиграться в DI и тестирование
	r.Post("/api/user/register", anonController.PostUserRegister)
	r.Post("/api/user/login", anonController.PostUserLogin)
	r.Post("/api/user/orders", orderController.PostUserOrders)
	r.Get("/api/user/orders", orderController.GetUserOrders)
	r.Get("/api/user/balance", userController.GetUserBalance)
	r.Post("/api/user/balance/withdraw", wController.PostUserBalanceWithdraw)
	r.Get("/api/user/withdrawals", wController.GetUserWithdrawals)

	return r
}
