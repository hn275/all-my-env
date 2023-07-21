package variables

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"sync"
	"time"

	"github.com/hn275/envhub/server/db"
	"gorm.io/gorm"
)

type base64VariableID = string

type variableHandler struct {
	*gorm.DB
}

var (
	Handlers *variableHandler

	// mapping repo id to it's counter, is refreshed every second
	counterMap map[uint32]uint16
	m          sync.Mutex
)

func init() {
	Handlers = &variableHandler{db.New()}
	m = sync.Mutex{}
	counterMap = make(map[uint32]uint16)
}

func RefreshVariableCounter() {
	for {
		m.Lock()
		counterMap = make(map[uint32]uint16)
		m.Unlock()
		time.Sleep(time.Second)
	}
}

// schema to generate id:
// `[repository id, time utc, process id, counter var]`.
// Where `counter var` is reset to 0 every second. The `id` is the
// base64-encoding of the byte array.
func genVariableID(repoID uint32) base64VariableID {
	bufSize := 14
	buf := make([]byte, bufSize)

	binary.BigEndian.PutUint32(buf[:4], repoID)

	t := time.Now().UTC().Unix()
	binary.BigEndian.PutUint32(buf[4:8], uint32(t))

	pid := os.Getpid()
	binary.BigEndian.PutUint32(buf[8:12], uint32(pid))

	counter := counterMap[repoID]
	m.Lock()
	counterMap[repoID] = counter + 1
	m.Unlock()
	binary.BigEndian.PutUint16(buf[12:], counter)

	return base64.StdEncoding.EncodeToString(buf)
}
