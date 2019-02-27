/**
This code is only used to show how to interact with the DB.addDatabase
It does not contain actual business logic yet
**/

package database

import (
	"time"

	"gopkg.in/mgo.v2"
)

var Session *mgo.Session
var err error

type Post struct {
	URL  string
	Text string
}

func Init() (*mgo.Session, error) {
	for i := 0; i < 5; i++ {
		Session, err = mgo.Dial("mongo:27017")
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic(err)
	}

	return Session, err
}
