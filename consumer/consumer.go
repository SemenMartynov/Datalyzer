package main

import (
	"flag"
	"fmt"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/crackcomm/nsqueue/consumer"
)

var db *sql.DB
var stmt *sql.Stmt
var err error

var (
	nsqdAddr    = flag.String("nsqd", "172.18.0.3:4150", "nsqd tcp address")
	dbHost      = flag.String("dbhost", "localhost", "Database address")
	maxInFlight = flag.Int("max-in-flight", 30, "Maximum amount of messages in flight to consume")
)

func checkErr(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func handleTest(msg *consumer.Message) {
	t := &time.Time{}
	t.UnmarshalBinary(msg.Body)
	duration := time.Since(*t)
	_, err = stmt.Exec(duration)
	checkErr(err)
	fmt.Printf("Consume latency: %s\n", duration)
	msg.Success()
}

func main() {
	flag.Parse()

	db, err = sql.Open("mysql", "caravel:caravel@tcp(" + *dbHost + ":3306)/testdb?charset=utf8")
	checkErr(err)
	stmt, err = db.Prepare("INSERT INTO kvtable (id, value) VALUES (NULL, ?)")
	checkErr(err)

	consumer.Register("latency-test", "consume", *maxInFlight, handleTest)
	consumer.Connect(*nsqdAddr)
	consumer.Start(true)
}
