package rdb

import (
	"encoding/json"
	"fmt"
	"github.com/8090Lambert/go-redis-parser/protocol"
)

// Some of HashEntry manager.
type HashMap struct {
	Field KeyObject
	Len   uint64
	Entry []HashEntry
}

// HashTable entry.
type HashEntry struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (r *ParseRdb) readHashMap(key KeyObject) error {
	length, _, err := r.loadLen()
	if err != nil {
		return err
	}
	hashTable := HashMap{Field: key, Len: length, Entry: make([]HashEntry, 0, length)}
	for i := uint64(0); i < length; i++ {
		field, err := r.loadString()
		if err != nil {
			return err
		}
		value, err := r.loadString()
		if err != nil {
			return err
		}
		hashTable.Entry = append(hashTable.Entry, HashEntry{Field: ToString(field), Value: ToString(value)})
	}
	//r.d1 = append(r.d1, hashTable.String())
	r.d1 = append(r.d1, hashTable)

	return nil
}

func (r *ParseRdb) readHashMapWithZipmap(key KeyObject) error {
	zipmap, err := r.loadString()
	if err != nil {
		return err
	}
	buf := newInput(zipmap)
	blen, err := buf.ReadByte()
	if err != nil {
		return err
	}

	length := int(blen)
	if blen > 254 {
		length, err = countZipmapItems(buf)
		if err != nil {
			return err
		}
		length /= 2
	}

	hashTable := HashMap{Field: key, Len: uint64(length), Entry: make([]HashEntry, 0, length)}
	for i := 0; i < length; i++ {
		field, err := loadZipmapItem(buf, false)
		if err != nil {
			return err
		}
		value, err := loadZipmapItem(buf, true)
		if err != nil {
			return err
		}
		hashTable.Entry = append(hashTable.Entry, HashEntry{Field: ToString(field), Value: ToString(value)})
	}
	//r.d1 = append(r.d1, hashTable.String())
	r.d1 = append(r.d1, hashTable)

	return nil
}

func (r *ParseRdb) readHashMapZiplist(key KeyObject) error {
	b, err := r.loadString()
	if err != nil {
		return err
	}
	buf := newInput(b)
	length, err := loadZiplistLength(buf)
	if err != nil {
		return err
	}
	length /= 2

	hashTable := HashMap{Field: key, Len: uint64(length), Entry: make([]HashEntry, 0, length)}
	for i := int64(0); i < length; i++ {
		field, err := loadZiplistEntry(buf)
		if err != nil {
			return err
		}
		value, err := loadZiplistEntry(buf)
		if err != nil {
			return err
		}
		hashTable.Entry = append(hashTable.Entry, HashEntry{Field: ToString(field), Value: ToString(value)})
	}
	//r.d1 = append(r.d1, hashTable.String())
	r.d1 = append(r.d1, hashTable)

	return nil
}

func (hm HashMap) Type() string {
	return protocol.Hash
}

func (hm HashMap) String() string {
	return fmt.Sprintf("{HashMap: {Key: %s, Len: %d, Entries: %s}}", hm.Key(), hm.Len, hm.Value())
}

func (hm HashMap) Key() string {
	return ToString(hm.Field)
}

func (hm HashMap) Value() string {
	itemStr, _ := json.Marshal(hm.Entry)
	return ToString(itemStr)
}

// 计算 hash 结构 field + value 的大小
func (hm HashMap) ConcreteSize() uint64 {
	var kv string
	for _, val := range hm.Entry {
		kv += ToString(val.Field) + ToString(val.Value)
	}
	return uint64(len([]byte(kv)))
}
