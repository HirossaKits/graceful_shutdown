package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
)

type seeder struct {
	db *sql.DB
}

func NewSeeder() *seeder {
	return &seeder{
		db: connect(),
	}
}

func main() {
	seeder := NewSeeder()
	seeder.initDB()
	seeder.seedDB()
	defer seeder.db.Close()
}

func connect() *sql.DB {
	mc := mysql.Config{
		User:                 "mysql",
		Passwd:               "mysql",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "test_db",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", mc.FormatDSN())
	if err != nil {
		log.Fatalf("failed to open mysql: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect mysql: %v", err)
	}
	return db
}

func (s *seeder) initDB() {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "seed", "init.sql")
	c, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to open init file: %v\n", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("failed to begin Tx: %v\n", err)
	}

	_, err = s.db.Exec(string(c))
	if err != nil {
		fmt.Printf("failed to exec query: %v\n", err)
		if err := tx.Rollback(); err != nil {
			log.Fatalf("failed to rollback DB: %v\n", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit DB: %v\n", err)
	}
}

func (s *seeder) seedDB() {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "seed", "seed.sql")
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open seed file: %v\n", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("failed to begin Tx: %v\n", err)
	}

	for scanner.Scan() {
		buf := make([]byte, 0, 256)
		bs := scanner.Bytes()

		if i := indexOf(bs, ';'); -1 < i {
			buf = append(buf, bs[:i]...)

			ch := make(chan error)
			go func() {
				_, err := s.db.Exec(string(buf))
				ch <- err
				if err != nil {
					fmt.Printf("failed to exec query: %v\n", err)
					if err := tx.Rollback(); err != nil {
						log.Fatalf("failed to rollback DB: %v\n", err)
					}
				}
			}()

			<-ch

			buf = bs[i+1:]
		} else {
			buf = append(buf, bs...)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit DB: %v\n", err)
	}
}

func indexOf(target []byte, search byte) int {
	index := -1
	for i, elem := range target {
		if elem == search {
			index = i
			break
		}
	}
	return index
}
