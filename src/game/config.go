package game

import (
	"fmt"
	"god"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"proto"
)

type M []map[string]interface{}

type Config struct {
	*god.Process
	Session *mgo.Session
	Db      *mgo.Database
	Mlist   map[string]M
}

func NewConfig() *Config {

	cfg := new(Config)
	cfg.Process = god.NewProcess(cfg)
	cfg.init()
	return cfg
}

func (self *Config) Handle(pID proto.PacketID, data god.Marshaler) (retID proto.PacketID, ret god.Marshaler, err error) {

	info := data.(proto.CfgFlush)
	result := proto.CfgRsp{false}
	self.Dial(info.Url)
	self.UseDB(info.Db)
	result.State = self.Load(info.Modem)

	ret = result

	return proto.CFG_FLUSH_RSP, ret, nil
}
func (self *Config) init() {
	self.Mlist = make(map[string]M)
}

func (self *Config) Get(name string) M {

	return self.Mlist[name]
}
func (self *Config) Dial(url string) {

	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Errorf("URL %s cannot connect!", url)
		return
	}
	self.Session = session
}
func (self *Config) UseDB(name string) {

	db := self.Session.DB(name)
	if db == nil {
		fmt.Errorf("DB %s cannot find!", name)
		return
	}
	self.Db = db
}
func (self *Config) Load(name string) bool {
	var data M
	collection := self.Db.C(name)
	err := collection.Find(bson.M{}).All(&data)

	if err != nil {
		fmt.Errorf(" %v is not found!", data)
		return false
	}
	self.Mlist[name] = data

	return true
}
func (self *Config) Close() {

	self.Session.Close()
}
