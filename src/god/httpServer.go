package god

import (
	"encoding/json"
	"net/http"
	"testing"
)

type svrinfo struct {
	Addr string
}

var serverList map[string]svrinfo

func init() {
	serverList = make(map[string]svrinfo)
}

func retSvrList(w http.ResponseWriter, req *http.Request) {
	var retList []string
	for _, v := range serverList {
		retList = append(retList, v.Addr)
	}
	b, _ := json.Marshal(retList)
	w.Write(b)
}

func updateSvrList(w http.ResponseWriter, req *http.Request) {
	nodeAddr := req.Header.Get("Node-Addr")
	if nodeAddr != "" {
		serverList[nodeAddr] = svrinfo{Addr: nodeAddr}
	}
}

func TestServer(t *testing.T) {
	http.HandleFunc("/locateGet", retSvrList)
	http.HandleFunc("/locatePost", updateSvrList)
	go http.ListenAndServe(":20000", nil)
}
