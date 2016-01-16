package main

import "github.com/boltdb/bolt"

type change struct {
	Bucket string
	Key    string
	Data   []byte
}

var db *bolt.DB
var changeChannel chan change

func pushChange(bucket string, key string, data []byte) {
	changeChannel <- change{bucket, key, data}
}

func boltWorker() {
	for {
		c := <-changeChannel
		db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(c.Bucket))
			b.Put([]byte(c.Key), c.Data)
			return nil
		})
	}
}
