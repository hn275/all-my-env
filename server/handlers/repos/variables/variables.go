package variables

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"sync"
	"time"

	"github.com/hn275/envhub/server/crypto"
	"github.com/hn275/envhub/server/db"
	"gorm.io/gorm"
)

type variableHandler struct {
	*gorm.DB
}

type EnvVariable struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	RepoURL string `json:"repo_url"`
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

func (v *EnvVariable) Cipher(repoID uint32) (*db.Variable, error) {
	id := genVariableID(repoID)
	idStr := base64.StdEncoding.EncodeToString(id)
	newValue, err := crypto.Encrypt(v.Value, id)
	if err != nil {
		return nil, err
	}

	s := db.Variable{
		ID:           idStr,
		Key:          v.Key,
		Value:        base64.StdEncoding.EncodeToString(newValue),
		RepositoryID: repoID,
	}

	return &s, nil
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
// Where `counter var` is reset to 0 every second.
func genVariableID(repoID uint32) []byte {
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

	return buf
}
