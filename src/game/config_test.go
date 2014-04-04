package game

import (
	"ext"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"runtime"
	"testing"
	"time"
)

type T struct {
	Name   string
	Id     int
	Statue bool
}

func create_tab(t *testing.T) {

	fmt.Println("create_tab")
	session, err := mgo.Dial("t4f-mam-13141.local")
	if err != nil {
		t.Error("cannot Dial")
		panic(err)
	}
	session.DB("test").DropDatabase()
	defer session.Close()

	db := session.DB("test")

	collection := db.C("one")

	err = collection.Insert(&T{"one", 222, false})
}

var na string

type routine struct {
	*Config
}

var da *M

func (self *routine) run(t *testing.T) {
	for {
		runtime.Gosched()
		ext.AssertT(t, self.Config.Get("one")[0]["name"] == "one", "")
	}
}

func (self *routine) setConfig(con *Config) {
	self.Config = con
}

/*func TestCfg(t *testing.T) {

	p1 := new(routine)
	p2 := new(routine)
	na = "one"
	con := NewConfig()
	p1.setConfig(con)
	p2.setConfig(con)

	create_tab(t)

	con.Dial("t4f-mam-13141.local")
	con.UseDB("test")
	con.Load("one")
	con.Close()
	ext.AssertT(t, con.Get("one")[0]["name"] == "one", "")
	//go p1.run(t)
	//go p2.run(t)


	session, err := mgo.Dial("t4f-mam-13141.local")
	if err != nil {
		t.Error("cannot Dial")
		panic(err)
	}

	db := session.DB("test")

	collection := db.C("one")

	idQuerier := bson.M{"name": "one"}
	change := bson.M{"$set": bson.M{"name": "hello"}}
	err = collection.Update(idQuerier, change)
	if err != nil {
		panic(err)
	}
	session.Close()
    ext.AssertT(t, con.Get("one")[0]["name"] == "one", "")
	con.Dial("t4f-mam-13141.local")
	con.UseDB("test")
	con.Load("one")
	con.Close()
	ext.AssertT(t, con.Get("one")[0]["name"] == "hello", "")
}

*/
func modify(t *testing.T) {
	session, err := mgo.Dial("t4f-mam-13141.local")
	if err != nil {
		t.Error("cannot Dial")
		panic(err)
	}
	db := session.DB("test")
	collection := db.C("one")
	idQuerier := bson.M{"name": "one"}
	change := bson.M{"$set": bson.M{"name": "jimyang"}}
	err = collection.Update(idQuerier, change)
	if err != nil {
		panic(err)
	}
	session.Close()
}

func TestCfgMul(t *testing.T) {
	create_tab(t)
	p1 := new(routine)
	p2 := new(routine)
	con := NewConfig()
	p1.setConfig(con)
	p2.setConfig(con)
	con.Dial("t4f-mam-13141.local")
	con.UseDB("test")
	con.Load("one")
	go p1.run(t)
	go p2.run(t)
	modify(t)
	time.Sleep(time.Second * 2)
	con.Load("one")
	con.Close()
	time.Sleep(2000)
}
