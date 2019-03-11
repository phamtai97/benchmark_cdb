package executor

import (
	"benchmark_cockroachdb/configvar"
	"benchmark_cockroachdb/storage"
	"database/sql"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

//Message struct
type Message struct {
	Id      string
	Balance int64
}

const balance = 1000000

// var msgQueue *lane.Queue
var iMsg = int64(0)

//randomInt
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

//GenerateID generate id
func GenerateID(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

//CreateMsg create mess
func CreateMsg(id string, balance int64) *Message {
	return &Message{
		Id:      id,
		Balance: balance,
	}
}

// //LoadMsg load msg
// func LoadMsg(wg *sync.WaitGroup, numMsg int64) {
// 	defer wg.Done()

// 	i := atomic.AddInt64(&iMsg, 1)
// 	for i <= numMsg {
// 		msgReq := CreateMsg()
// 		msgQueue.Enqueue(msgReq)
// 		i = atomic.AddInt64(&iMsg, 1)
// 	}
// }

// //PrepareData prepare data
// func PrepareData(numGoroutine int, numMsg int64) {
// 	var wg sync.WaitGroup
// 	for i := 0; i < numGoroutine; i++ {
// 		wg.Add(1)
// 		go LoadMsg(&wg, numMsg)
// 	}
// 	wg.Wait()
// }

//InsertMsg insert message
// func InsertMsg(wg *sync.WaitGroup, db *sql.DB, numberMsg int64) {
// 	defer wg.Done()
// 	i := atomic.AddInt64(&iMsg, 1)

// 	for i <= numberMsg {
// 		id := strconv.FormatInt(i, 10)
// 		msg := CreateMsg(id, balance)
// 		if _, err := db.Exec(
// 			"INSERT INTO benchmark.accounts"+
// 				" (id, balance) VALUES ($1, $2)", msg.Id, msg.Balance); err != nil {
// 			log.Fatal("Insert msg error: ", err)
// 		}
// 		i = atomic.AddInt64(&iMsg, 1)
// 	}
// }

func InsertMsg(wg *sync.WaitGroup, db *sql.DB, numberMsg int64) {
	defer wg.Done()
	i := atomic.AddInt64(&iMsg, 1)

	for i <= numberMsg {
		id := strconv.FormatInt(i, 10)
		msg := CreateMsg(id, balance)
		if _, err := db.Exec(
			"INSERT INTO accounts"+
				" (id, balance) VALUES ($1, $2)", msg.Id, msg.Balance); err != nil {
			log.Fatal("Insert msg error: ", err)
		}
		i = atomic.AddInt64(&iMsg, 1)
	}
}

//QueryAccount query account
// func QueryAccount(wg *sync.WaitGroup, db *sql.DB, numberMsg int64) {
// 	defer wg.Done()
// 	i := atomic.AddInt64(&iMsg, 1)
// 	var idRes string
// 	var balanceRes int64

// 	for i <= numberMsg {
// 		// id := strconv.FormatInt(1, 10)
// 		id := "1"
// 		if err := db.QueryRow("SELECT * FROM benchmark.accounts WHERE id = $1", id).Scan(&idRes, &balanceRes); err != nil {
// 			log.Fatal("Query fail: ", err)
// 		}
// 		i = atomic.AddInt64(&iMsg, 1)
// 	}
// }

func QueryAccount(wg *sync.WaitGroup, db *sql.DB, numberMsg int64) {
	defer wg.Done()
	i := atomic.AddInt64(&iMsg, 1)
	var idRes string
	var balanceRes int64

	for i <= numberMsg {
		// id := strconv.FormatInt(1, 10)
		id := "1"
		if err := db.QueryRow("SELECT * FROM accounts WHERE id = $1", id).Scan(&idRes, &balanceRes); err != nil {
			log.Fatal("Query fail: ", err)
		}
		i = atomic.AddInt64(&iMsg, 1)
	}
}

//RunInsert benchmark
func RunInsert(numberGoroutine int, db *sql.DB, numberMsg int64) {
	var wg sync.WaitGroup
	for i := 0; i < numberGoroutine; i++ {
		wg.Add(1)
		go InsertMsg(&wg, db, numberMsg)
	}
	wg.Wait()
}

//RunQuery run query
func RunQuery(numberGoroutine int, db *sql.DB, numberMsg int64) {
	var wg sync.WaitGroup
	for i := 0; i < numberGoroutine; i++ {
		wg.Add(1)
		go QueryAccount(&wg, db, numberMsg)
	}
	wg.Wait()
}

//HandlerExec handle exec
func HandlerExec(configVar *configvar.CliConfVar, db *sql.DB) {
	switch configVar.Type {
	case 0:
		RunInsert(configVar.NumGoroutine, db, configVar.NumMsg)
		break
	case 1:
		RunQuery(configVar.NumGoroutine, db, configVar.NumMsg)
	}
}

//Execute execute
func Execute(configVar *configvar.CliConfVar) {
	user := "postgres"
	password := "taiptht"
	nameDB := "benchmark"
	hostDB := "localhost"
	portDB := 5432
	ssmode := true

	storage, err := storage.InitStoragePostgre(user, password, hostDB, portDB, nameDB, ssmode)
	if err != nil {
		log.Fatal("[Server] Create Storage: %s", err.Error())
	}

	log.Info("---- Start benchmark ----")

	start := time.Now()
	HandlerExec(configVar, storage.GetDB())

	t := time.Now()
	elapsed := t.Sub(start)

	seconds := elapsed.Seconds()

	cal := math.RoundToEven(float64(configVar.NumMsg) / seconds)

	log.WithFields(log.Fields{"Average Msg/s ": cal}).Info("---- Benchmark DONE ----")
}
