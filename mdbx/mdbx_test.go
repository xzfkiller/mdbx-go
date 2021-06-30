package mdbx

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTest1(t *testing.T) {
	env, err := NewEnv()
	if err != nil {
		t.Fatalf("Cannot create environment: %s", err)
	}
	path, err := ioutil.TempDir("", "mdbx_test")
	if err != nil {
		t.Fatalf("Cannot create temporary directory")
	}
	err = os.MkdirAll(path, 0770)
	defer os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Cannot create directory: %s", path)
	}
	err = env.Open(path)
	defer env.Close()
	if err != nil {
		t.Fatalf("Cannot open environment: %s", err)
	}

	var db DBI
	numEntries := 10
	var data = map[string]string{}
	var key string
	var val string
	for i := 0; i < numEntries; i++ {
		key = fmt.Sprintf("Key-%d", i)
		val = fmt.Sprintf("Val-%d", i)
		data[key] = val
	}
	err = env.Update(func(txn *Txn) (err error) {
		db, err = txn.OpenRoot(0)
		if err != nil {
			return err
		}

		for k, v := range data {
			err = txn.Put(db, []byte(k), []byte(v), NoOverwrite)
			if err != nil {
				return fmt.Errorf("put: %v", err)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	err = env.View(func(txn *Txn) error {
		numEntries := 10

		var key string
		var val string
		var err error
		var vRead []byte

		db, err = txn.OpenRoot(0)
		if err != nil {
			return err
		}

		for i := 0; i < numEntries; i++ {
			key = fmt.Sprintf("Key-%d", i)
			val = fmt.Sprintf("Val-%d", i)

			vRead, err = txn.Get(db, []byte(key))
			if IsNotFound(err) {
				return errors.New("MDBX: not found")
			}

			if bytes.Compare(vRead, []byte(val)) != 0 {
				t.Errorf("[BAD]: Value not match, expect: %v, got: %v", val, vRead)
			} else {
				t.Logf("[GOOD]: Found value for key and matched")
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
