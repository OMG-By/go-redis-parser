package rdb

import (
	"fmt"
	"github.com/8090Lambert/go-redis-parser/protocol"
)

type SelectionDB struct {
	Index uint64
}

type ResizeDB struct {
	DBSize     uint64
	ExpireSize uint64
}

func Resize(dbSize, expireSize uint64) ResizeDB {
	return ResizeDB{DBSize: dbSize, ExpireSize: expireSize}
}

func (r ResizeDB) String() string {
	return fmt.Sprintf("{ResizeDB: %s}", r.Value())
}

func (r ResizeDB) Key() string {
	return "resize db"
}

func (r ResizeDB) Value() string {
	return fmt.Sprintf("{dbsize: %d, expireSize: %d}", r.DBSize, r.ExpireSize)
}

func (r ResizeDB) ConcreteSize() uint64 {
	return 0
}

func (r ResizeDB) Type() string {
	return protocol.SelectDB
}

func Selection(index uint64) SelectionDB {
	return SelectionDB{Index: index}
}

func (s SelectionDB) Type() string {
	return protocol.SelectDB
}

func (s SelectionDB) String() string {
	return fmt.Sprintf("{Select: %d}", s.Index)
}

func (s SelectionDB) Key() string {
	return "select"
}

func (s SelectionDB) Value() string {
	return ToString(s.Index)
}

func (s SelectionDB) ConcreteSize() uint64 {
	return 0
}
