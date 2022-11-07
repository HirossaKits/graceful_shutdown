package main

import (
	"fmt"
	"log"
	"net/http"

	"entgo.io/ent/examples/fs/ent"
	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"

	"entgo.io/ent/dialect/sql"
)

func main() {

	c := connect()
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		r.Context()
		setUncommitedTranc(c)
		w.Write([]byte("hello world!"))
	})
	const port = 3080
	fmt.Printf("server listening on port: %v ...", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}

func connect() *ent.Client {
	mc := mysql.Config{
		User:                 "mysql",
		Passwd:               "mysql",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	drv, err := sql.Open("mysql", mc.FormatDSN())
	if err != nil {
		log.Fatalf("failed to connect mysql: %v", err)
	}

	return ent.NewClient(ent.Driver(drv))
}

func setUncommitedTranc(c *ent.Client) {
	fmt.Println("ok")
}
