package main

import (
	"log"

	faker "github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type floatRecord struct {
	ID           int     `db:"id"`
	NumericParam float64 `db:"numeric_param"`
	RealParam    float64 `db:"real_param"`
}

type decimalRecord struct {
	ID           int             `db:"id"`
	NumericParam decimal.Decimal `db:"numeric_param"`
	RealParam    decimal.Decimal `db:"real_param"`
}

func main() {
	// open connection with postgres
	db, err := sqlx.Open("postgres", "user=postgres password=123456 dbname=postgres host=localhost port=5000 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 10000; i++ {
		var a float32
		if err := faker.FakeData(&a); err != nil {
			log.Println(err)
		}

		if _, err := db.Exec("insert into records (numeric_param,real_param) values ($1, $2)", a, a); err != nil {
			log.Println(err)
		}
	}
}
