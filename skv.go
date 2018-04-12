package main

import "github.com/boltdb/bolt"

type Skv struct {
	db *bolt.DB
}

func (v *Skv) Init(dbfile string) {
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		Error.Fatal("skv init fail: %s", err)
		return
	}
	v.db = db
	Info.Println("skv init suc", dbfile)
}

func (v *Skv) Close() {
	v.db.Close()
}
