package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	indexDB *leveldb.DB
)

func initDB(path string) error {
	var err error
	indexDB, err = leveldb.OpenFile(path, nil)
	return err
}

func closeDB() error {
	if indexDB != nil {
		return indexDB.Close()
	}
	return nil
}

func buildKey(prefix string, unitID int64, kw string, id int64) []byte {
	return []byte(fmt.Sprintf("%s:%d:%s:%d", prefix, unitID, kw, id))
}

func decodeStrings(s string) []string {
	var lastIdx int
	var lastW rune
	var res []string
	for i, w := range s {
		if i > 0 {
			var buf = make([]byte, i-lastIdx)
			utf8.EncodeRune(buf, lastW)
			res = append(res, string(buf))
		}
		lastIdx = i
		lastW = w
	}

	var buf = make([]byte, len(s)-lastIdx)
	utf8.EncodeRune(buf, lastW)
	res = append(res, string(buf))

	return res
}

// generateKeywords
func generateKeywords(c string) map[string]struct{} {
	ss := decodeStrings(string(c))
	var kws = make(map[string]struct{}, len(ss)*(len(ss)-1)/2)
	for i := 0; i < len(ss); i++ {
		for j := i + 1; j <= len(ss); j++ {
			kws[strings.Join(ss[i:j], "")] = struct{}{}
		}
	}
	return kws
}

// addIndex
func addIndex(prefix, c string, id, unitID int64) error {
	c = strings.ToUpper(string(c))
	for kw := range generateKeywords(c) {
		err := indexDB.Put(buildKey(prefix, unitID, kw, id), []byte(string(c)), nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// delIndex
func delIndex(prefix, c string, id, unitID int64) error {
	c = strings.ToUpper(string(c))
	for kw := range generateKeywords(c) {
		err := indexDB.Delete(buildKey(prefix, unitID, kw, id), nil)
		if err != nil {
			return err
		}
	}
	return nil
}

type result struct {
	ID   int64
	Item string
}

// search
func search(prefix, kw string, unitID int64, max int) ([]result, error) {
	kw = strings.ToUpper(kw)
	var results = make([]result, 0, max)
	prefix = fmt.Sprintf("%s:%d:%s:", prefix, unitID, kw)

	iter := indexDB.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, []byte(prefix)) {
			break
		}

		id, err := strconv.ParseInt(string(key[len(prefix):]), 10, 64)
		if err != nil {
			return nil, err
		}

		results = append(results, result{
			ID:   id,
			Item: string(iter.Value()),
		})
		max--
		if max <= 0 {
			break
		}
	}

	return results, iter.Error()
}
