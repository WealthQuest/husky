package husky

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
)

type UUID interface {
	Id() int64
	Uuid() string
	TimeID() string
}

type _UUID struct {
	ins *snowflake.Node
}

func (u *_UUID) Id() int64 {
	return u.ins.Generate().Int64()
}

func (u *_UUID) Uuid() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func (u *_UUID) TimeID() string {
	now := time.Now().Format("20060102150405")
	n := rand.Int63n(1000000000000000000)
	return fmt.Sprintf("%s%018d", now, n)
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
