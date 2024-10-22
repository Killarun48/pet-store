package main

import (
	"app/internal/infrastructure/db"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "app/docs"
	"app/internal/infrastructure/responder"
	"app/internal/modules"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var pathDB = "./petstore.db"

type Server struct {
	srv     *http.Server
	users   map[string]string
	sigChan chan os.Signal
}

func NewServer(addr string) *Server {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("get envoironment variables")

	server := &Server{
		sigChan: make(chan os.Signal, 1),
		users:   make(map[string]string),
	}

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	bd, err := db.NewDataBaseSqlite(pathDB)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("initialize database")

	err = bd.Migrate()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("migrate database")

	if pathDB != "./petstore.db" {
		FillFakeData(*bd)
	}

	repositories := modules.NewRepository(bd.DB)
	services := modules.NewService(repositories)
	respond := responder.NewResponder()

	c := modules.NewController(services, respond)

	log.Println("initialize controllers")

	// Инициализируем маршруты
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/v2", func(r chi.Router) {
		r.Mount("/user", c.InitRoutesUser())
		r.Mount("/pet", c.InitRoutesPet())
		r.Mount("/store", c.InitRoutesStore())
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"operationsSorter": "function(a, b) { return ((aOp=JSON.parse(JSON.stringify(a)).operation?.operationId) && (bOp=JSON.parse(JSON.stringify(b)).operation?.operationId)) ? (aOp[0] > bOp[0] ? 1 : -1) : 0; }",
		}),
	))

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	server.srv = srv

	time.Sleep(1 * time.Second)

	return server
}

func (s *Server) Serve() {
	go func() {
		log.Println("Starting server...")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-s.sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}

func (s *Server) Stop() {
	s.sigChan <- syscall.Signal(1)
}

func FillFakeData(db db.DataBaseSqlite) {
	addUsers(db)
}

func addUsers(db db.DataBaseSqlite) {
	type fakeUser struct {
		ID         int    `fake:"{number[2,100]}"`
		UserName   string `fake:"{name}"`
		FirstName  string `fake:"{first_name}"`
		LastName   string `fake:"{last_name}"`
		Email      string `fake:"{email}"`
		Password   string `fake:"{password}"`
		Phone      string `fake:"{phone}"`
		UserStatus int    `fake:"{1,2}"`
	}

	var fakeU fakeUser

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	fakeU.ID = gofakeit.Number(1, 100)
	fakeU.UserName = gofakeit.Username() + strconv.Itoa(2)
	fakeU.FirstName = gofakeit.FirstName()
	fakeU.LastName = gofakeit.LastName()
	fakeU.Email = gofakeit.Email()
	fakeU.Password = gofakeit.Password(true, true, true, true, true, 8)
	fakeU.Phone = gofakeit.Phone()
	fakeU.UserStatus = gofakeit.Number(1, 2)

	insertBuilder := sq.Insert("users").Columns(
		"username",
		"first_name",
		"last_name",
		"email",
		"password",
		"phone",
		"user_status",
	)
	insertBuilder = insertBuilder.Values(
		fakeU.UserName,
		fakeU.FirstName,
		fakeU.LastName,
		fakeU.Email,
		fakeU.Password,
		fakeU.Phone,
		fakeU.UserStatus,
	)
	_, err := insertBuilder.
		RunWith(db.DB).
		Exec()
	if err != nil {
		log.Println(err)
		return
	}
}
