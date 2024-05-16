package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Blank import to register the sqlite3 driver
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"google.golang.org/protobuf/proto"
)

type dbInfo struct {
	id       string
	db       *sql.DB
	file     string
	stmt     *sql.Stmt
	inserted bool
}

type dbMap map[string]*dbInfo

var databases dbMap

func newDBInfo(file string) {
	id := uuid.New().String()

	db, err := sql.Open("sqlite3", file+"?_journal=WAL")
	if err != nil {
		panic(err.Error())
	}

	stmt, err := db.Prepare("INSERT INTO message(pubsubTopic, contentTopic, payload, version, timestamp, id, messageHash, storedAt) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	databases[id] = &dbInfo{
		id:   id,
		file: file,
		db:   db,
		stmt: stmt,
	}
}

func (d dbMap) values() []*dbInfo {
	var result []*dbInfo
	for _, x := range d {
		result = append(result, x)
	}
	return result
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Use: populatedb [path_to_db1 path_to_db2 ...]")
		return
	}

	databases = make(dbMap)
	for _, dbPath := range os.Args[1:] {
		newDBInfo(dbPath)
	}

	pubsubTopics := []string{
		"/waku/2/rs/1/1",
		"/waku/2/rs/1/2",
		"/waku/2/rs/1/3",
	}

	// Wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	timeInterval := 5 * time.Second
	t := time.NewTicker(timeInterval)

	minuteTimer := time.NewTimer(0)

	fmt.Print("                                                                  \t\t\t\t\t\t")
	i := 0
	for range databases {
		i++
		fmt.Printf("db_%d\t", i)
	}
	fmt.Println()

	for {
		select {
		case currTime := <-minuteTimer.C:
			fmt.Println(currTime)
			minuteTimer.Reset(1 * time.Minute)
		case <-t.C:
			insertCnt := rand.Intn(len(databases)) + 1 // Insert in at least one store node

			dbArr := databases.values()
			rand.Shuffle(len(dbArr), func(i, j int) { dbArr[i], dbArr[j] = dbArr[j], dbArr[i] })

			for _, x := range dbArr {
				x.inserted = false
			}

			now := time.Now().UnixNano()
			msg := &pb.WakuMessage{
				Timestamp:    proto.Int64(now),
				Payload:      []byte{1, 2, 3, 4, 5},
				ContentTopic: "test",
			}
			envelope := protocol.NewEnvelope(msg, now, pubsubTopics[rand.Intn(len(pubsubTopics))])
			hash := envelope.Hash()
			for i := 0; i < insertCnt; i++ {
				dbArr[i].inserted = true
				_, err := dbArr[i].stmt.Exec([]byte(envelope.PubsubTopic()), []byte(msg.ContentTopic), msg.Payload, msg.GetVersion(), msg.GetTimestamp(), envelope.Index().Digest, hash.Bytes(), now)
				if err != nil {
					panic(err.Error())
				}
			}

			fmt.Printf("%s\t%s\t%v\t", hash.String(), envelope.PubsubTopic(), msg.GetTimestamp())
			for _, x := range databases {
				fmt.Printf("%v\t", x.inserted)
			}
			fmt.Println()

		case <-ch:
			return
		}
	}
}
