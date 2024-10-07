package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "funnypicture"
	pathkey := CASPathTransformFunc(key)

	expectedOriginalKey := "21d7a8342d65b75149512f65125e3519e2e63fd8"
	expectedPathname := "21d7a/8342d/65b75/14951/2f651/25e35/19e2e/63fd8"

	pathname := pathkey.PathName
	original := pathkey.FileName
	if pathname != expectedPathname {
		t.Errorf("have %s want %s", pathname, expectedPathname)
	}
	if original != expectedOriginalKey {
		t.Errorf("have %s want %s", original, expectedOriginalKey)
	}
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "funnydata"
	data := []byte("lol lol funny data lmaooooooo")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	count := 50
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("funnydata_%d", i)
		data := []byte("lol lol funny data lmaooooooo")

		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key: %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to NOT have key: %s", key)
		}
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
