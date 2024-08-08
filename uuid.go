package husky

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/rand"
)

type UUID interface {
	Gen() int64
}

type _UUID struct {
	ins *snowflake.Node
}

func (u *_UUID) Gen() int64 {
	return u.ins.Generate().Int64()
}

var uuidIns map[string]UUID

func init() {
	uuidIns = make(map[string]UUID)
}

func InitUUID(nodeSeq int64, key ...string) {
	if nodeSeq >= 1024 {
		panic("`nodeSeq` must be less than 1024")
	}
	node, err := snowflake.NewNode(nodeSeq)
	if err != nil {
		panic(err)
	}
	_ins := &_UUID{node}
	if len(key) == 0 {
		uuidIns[""] = _ins
	} else {
		uuidIns[key[0]] = _ins
	}
}

func Uuid(key ...string) UUID {
	if len(key) == 0 {
		return uuidIns[""]
	} else {
		return uuidIns[key[0]]
	}
}

func UuidStr() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func UuidTime() string {
	now := time.Now().Format("20060102150405")
	n := rand.Int63n(1000000000000000000)
	return fmt.Sprintf("%s%018d", now, n)
}
