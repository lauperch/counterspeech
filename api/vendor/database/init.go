/**
This code is only used to show how to interact with the DB.addDatabase
It does not contain actual business logic yet
**/

package database

var session *mgo.Session

type Post struct {
	URL  string
	Text string
}

func Init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	if err != nil {
		return DB, err
	}

	return session, err
}
