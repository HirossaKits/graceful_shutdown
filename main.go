package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"graceful-shutdown/ent"

	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
)

func main() {
	const port = 3080

	client := connect()
	server := &http.Server{Addr: fmt.Sprintf("localhost:%v", port), Handler: service(client)}

	serverCtx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		fmt.Println("started server shutdown ...")

		timerCtx, _ := context.WithTimeout(serverCtx, 10*time.Second)

		go func() {
			<-timerCtx.Done()
			if timerCtx.Err() == context.DeadlineExceeded {
				log.Fatal("time exceeded ... force shutdown")
			}
		}()

		err := server.Shutdown(timerCtx)
		if err != nil {
			log.Fatalln(err)
		}
		cancel()
	}()

	fmt.Printf("server listening on port: %v ...\n", port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}

	<-serverCtx.Done()
	fmt.Println("server was shutdown")
}

func connect() *ent.Client {
	mc := mysql.Config{
		User:                 "mysql",
		Passwd:               "mysql",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "test_db",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	client, err := ent.Open("mysql", mc.FormatDSN())
	if err != nil {
		log.Fatalf("failed to connect mysql: %v", err)
	}

	return client
}

func service(client *ent.Client) http.Handler {
	r := chi.NewRouter()
	r.Get("/", rootHandler(client))
	return r
}

func rootHandler(client *ent.Client) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tx, err := client.Tx(ctx)
		if err != nil {
			fmt.Printf("failed to start transaction: %v\n", err)
		}

		user, err := tx.Users.Create().SetName("test0").SetEmail("test0@mail.com").Save(ctx)
		if err != nil {
			fmt.Printf("failed to create a user: %v\n", err)
		}
		fmt.Printf("%+v\n", user)
		// tx.Commit()
	}
	return http.HandlerFunc(fn)
}
