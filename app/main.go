package main

import (
	"log"

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

	var fRecords []floatRecord
	var dRecords []decimalRecord
	query := "select * from records"

	err = db.Select(&fRecords, query)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Select(&dRecords, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range fRecords {
		log.Println(r.ID, r.NumericParam, r.RealParam)
	}

	for _, r := range dRecords {
		log.Println(r.ID, r.NumericParam, r.RealParam)
	}
}
