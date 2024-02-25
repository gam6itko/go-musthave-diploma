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

	if err := initDatabaseSchema(_db); err != nil {
		log.Fatalf("schema init error. %s", err)
	}

}

func initDatabaseSchema(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	// server_init.sql
	sqlQuery := `CREATE TABLE IF NOT EXISTS public.user
(
    id       BIGSERIAL PRIMARY KEY,
    login varchar,
    password varchar,
    balance_current NUMERIC(7,2) DEFAULT 0 NOT NULL,
    balance_withdraw NUMERIC(7,2) DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS public.order
(
    id  		BIGINT PRIMARY KEY,
    created_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id 	BIGINT,
    status 		SMALLINT,
    sum NUMERIC(7,2),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id)
);

CREATE TABLE IF NOT EXISTS public.withdrawal
(
    id  BIGSERIAL PRIMARY KEY,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    order_id BIGINT,
    user_id BIGINT,
    sum NUMERIC(7,2),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id)
);`
	_, err := db.Exec(sqlQuery)
	return err
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

//func startAccuralPolling() {
//	repo := diploma.NewOrderRepository(_db)
//
//	ticker := time.NewTicker(5 * time.Second)
//	for range ticker.C {
//		orderList, err := repo.FindByStatus(context.TODO(), diploma.StatusProcessing)
//		if err != nil {
//			log.Printf("error: %s", err)
//			continue
//		}
//
//		for _, o := range orderList {
//			acc, err := _accClient.Get(o.ID)
//			if err != nil {
//				log.Printf(err.Error())
//			}
//			accStatus, err := diploma.OrderStatusFromString(acc.Status)
//			if err != nil {
//				log.Printf("error. %s", err)
//			}
//			if o.Status == accStatus {
//				continue
//			}
//
//			err = repo.UpdateStatus(context.TODO(), o.ID, o.Status, o.Accural)
//			if err != nil {
//				log.Printf("update status error. %s", err)
//			}
//		}
//	}
//}
