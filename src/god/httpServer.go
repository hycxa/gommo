package god

import (
	"encoding/json"
	"net/http"
	"testing"
)

type svrinfo struct {
	Addr    string
	SvrType string
}

var serverList map[string]svrinfo

func init() {
	serverList = make(map[string]svrinfo)
}

func retSvrList(w http.ResponseWriter, req *http.Request) {
	svrType := req.Header.Get("service")
	var retList []string
	for _, v := range serverList {
		if v.SvrType == svrType {
			retList = append(retList, v.Addr)
		}
	}
	b, _ := json.Marshal(retList)
	w.Write(b)
}

func updateSvrList(w http.ResponseWriter, req *http.Request) {
	nodeAddr := req.Header.Get("Node-Addr")
	nodeType := req.Header.Get("service")
	if nodeAddr != "" && nodeType != "" {
		serverList[nodeAddr] = svrinfo{Addr: nodeAddr, SvrType: nodeType}
	}
}

func TestServer(t * testing.T) {
	http.HandleFunc("/locateGet", retSvrList)
	http.HandleFunc("/locatePost", updateSvrList)
	go http.ListenAndServe(":20000", nil)
}
