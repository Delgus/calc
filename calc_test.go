package main

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type testRow struct {
	NumericParam decimal.Decimal `db:"numeric_param"`
	RealParam    decimal.Decimal `db:"real_param"`
}

type testCase struct {
	x, y, wantFloat                   float64
	name, wantString, wantStringFixed string
}

func newTestCase(name string, x, y, wantFloat float64, wantString, wantStringFixed string) *testCase {
	return &testCase{
		name:            name,
		x:               x,
		y:               y,
		wantString:      wantString,
		wantStringFixed: wantStringFixed,
		wantFloat:       wantFloat,
	}
}

func (tc *testCase) save(db *sqlx.DB) error {
	query := "insert into records (name, numeric_param, real_param) values ($1, $2, $3),($1, $4, $5)"
	_, err := db.Exec(query, tc.name, tc.x, tc.x, tc.y, tc.y)
	return err
}

func (tc *testCase) check(db *sqlx.DB, t *testing.T) {
	var rows []testRow
	query := "select numeric_param, real_param::numeric from records where name = $1"
	if err := db.Select(&rows, query, tc.name); err != nil {
		t.Error(err)
	}

	var numeric, real decimal.Decimal
	for i := range rows {
		numeric = numeric.Add(rows[i].NumericParam)
		real = real.Add(rows[i].RealParam)
	}

	if f, _ := numeric.Round(2).Float64(); f != tc.wantFloat {
		t.Errorf(`%s: unexpected float64 from numeric - %v, expected %v`, tc.name, f, tc.wantFloat)
	}
	if s := numeric.Round(2).String(); s != tc.wantString {
		t.Errorf(`%s: unexpected string from numeric - %v, expected %v`, tc.name, s, tc.wantString)
	}
	if sf := numeric.Round(2).StringFixed(2); sf != tc.wantStringFixed {
		t.Errorf(`%s: unexpected string fixed from numeric - %v, expected %v`, tc.name, sf, tc.wantStringFixed)
	}

	if f, _ := real.Round(2).Float64(); f != tc.wantFloat {
		t.Errorf(`%s: unexpected float64 from real - %v, expected %v`, tc.name, f, tc.wantFloat)
	}
	if s := real.Round(2).String(); s != tc.wantString {
		t.Errorf(`%s: unexpected string from real - %v, expected %v`, tc.name, s, tc.wantString)
	}
	if sf := real.Round(2).StringFixed(2); sf != tc.wantStringFixed {
		t.Errorf(`%s: unexpected string fixed from real - %v, expected %v`, tc.name, sf, tc.wantStringFixed)
	}
}

func TestAll(t *testing.T) {
	// open connection with postgres
	db, err := sqlx.Open("postgres", "user=postgres password=123456 dbname=postgres host=localhost port=5000 sslmode=disable")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	defer db.Exec("delete from records")

	testCases := []*testCase{
		newTestCase("c1", 0.001, 0.003, 0.00, "0", "0.00"),
		newTestCase("c2", 0.002, 0.003, 0.01, "0.01", "0.01"),
		newTestCase("c3", 0.00099, 0.00401, 0.01, "0.01", "0.01"),
		newTestCase("c4", 0.00099, 0.00400, 0.00, "0", "0.00"),
		newTestCase("c5", 0.99, 0.01, 1.00, "1", "1.00"),
		newTestCase("c6", 0.99, 0.001, 0.99, "0.99", "0.99"),
		newTestCase("c7", 0.004, 0.99, 0.99, "0.99", "0.99"),
		newTestCase("c8", 0.005, 0.99, 1.00, "1", "1.00"),
		newTestCase("c9", -0.001, -0.003, 0.00, "0", "0.00"),
		newTestCase("c10", -0.002, -0.003, -0.01, "-0.01", "-0.01"),
		newTestCase("c11", -0.00099, -0.00401, -0.01, "-0.01", "-0.01"),
		newTestCase("c12", -0.00099, -0.00400, 0.00, "0", "0.00"),
		newTestCase("c13", -0.99, -0.01, -1.00, "-1", "-1.00"),
		newTestCase("c14", -0.99, -0.001, -0.99, "-0.99", "-0.99"),
		newTestCase("c15", -0.004, -0.99, -0.99, "-0.99", "-0.99"),
		newTestCase("c16", -0.005, -0.99, -1.00, "-1", "-1.00"),
		newTestCase("c17", 9999.9999, 9999.9999, 20000.00, "20000", "20000.00"),
		newTestCase("c18", -9999.9999, -9999.9999, -20000.00, "-20000", "-20000.00"),
	}

	for i := range testCases {
		if err := testCases[i].save(db); err != nil {
			t.Error(err)
			return
		}
	}

	for i := range testCases {
		testCases[i].check(db, t)
	}
}
