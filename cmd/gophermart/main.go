package main

import (
	"database/sql"
	"flag"
	"github.com/gam6itko/go-musthave-diploma/internal"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var _db *sql.DB
var _jwtIssuer *internal.JWTIssuer

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var dbDsn string
	if envVal, exists := os.LookupEnv("DATABASE_URI"); exists {
		dbDsn = envVal
	}

	var err error
	_db, err = sql.Open("pgx", dbDsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatal(err)
	}

	//jwt
	//jwt token create
	jwtKey, exists := os.LookupEnv("JWT_KEY")
	if !exists {
		log.Fatal("env JWT_KEY not defined")
	}
	_jwtIssuer = internal.NewJWTIssuer([]byte(jwtKey))
}

func main() {
	var bindAddr, dbDsn, accuralAddr string
	var bindAddrDef, dbDsnDef, accuralAddrDef string

	if tmp, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		bindAddrDef = tmp
	}
	if tmp, exists := os.LookupEnv("DATABASE_URI"); exists {
		dbDsnDef = tmp
	}
	if tmp, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); exists {
		accuralAddrDef = tmp
	}

	flag.StringVar(&bindAddr, "a", bindAddrDef, "Net address host:port")
	flag.StringVar(&dbDsn, "d", dbDsnDef, "Database DSN")
	flag.StringVar(&accuralAddr, "r", accuralAddrDef, "accural system address")
	flag.Parse()

	server := &http.Server{
		Addr:    bindAddr,
		Handler: newRouter(),
	}

	log.Printf("Starting server on %s", bindAddr)
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
