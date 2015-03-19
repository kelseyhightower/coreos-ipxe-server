package boltdb

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var buckets = []string{"machines", "profiles"}

type Store struct {
	db *bolt.DB
}

func New(path string) (*Store, error) {
	dbOptions := bolt.Options{
		Timeout: time.Duration(time.Second * 10),
	}
	db, err := bolt.Open(path, 0644, &dbOptions)
	if err != nil {
		return nil, err
	}
	for _, bucket := range buckets {
		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return &Store{db: db}, nil
}

func (s *Store) Get(bucket, id string) ([]byte, error) {
	var result []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		result = tx.Bucket([]byte(bucket)).Get([]byte(id))
		return nil
	})
	return result, err
}

func (s *Store) Save(bucket, id string, value []byte) error {
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(id), value)
		return err
	})
	return nil
}

func (s *Store) List(bucket string) (map[string][]byte, error) {
	m := make(map[string][]byte)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			m[string(k)] = v
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
