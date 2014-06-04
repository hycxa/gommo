package god

import (
	"encoding/json"
	"ext"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Node interface {
	Connect(name string) (ok bool)
	Nodes() []Node
	//Processes() []Process
}

type node struct {
}

var (
	workTab  [WORKER_NUM_LIMIT]*Worker
	nodeInfo NodeInfo
)

func NewNode(name string) Node {
	return &node{}
}

func (n *node) Connect(name string) (ok bool) {
	return true
}

func (n *node) Nodes() []Node {
	return nil
}

func init() {
	rand.Seed(time.Now().Unix())
}

func GetOneWorker() *Worker {
	index := rand.Intn(WORKER_NUM_LIMIT)
	if workTab[index] == nil {
		workTab[index] = NewWorker()
	}
	return workTab[index]
}

func NodeInit(name, network, addr string, port int) {
	nodeInfo.Name = name
	nodeInfo.Network = network
	nodeInfo.String = addr + ":" + strconv.Itoa(port)
	otherAddr := addr + ":" + strconv.Itoa(port+1)
	mes := NewMessenger()
	NewAcceptor(mes, REMOTE_NODE_TYPE, network, nodeInfo.String)
	NewAcceptor(mes, CLIENT_TYPE, network, otherAddr)
	ConnOtherSvr(mes)
}

func GetNodeInfo() *NodeInfo {
	return &nodeInfo
}

func ConnOtherSvr(mes Messenger) error {
	client := &http.Client{}
	reqPost, err := http.NewRequest("POST", "http://127.0.0.1:20000/locatePost", nil)
	if err != nil {
		return ext.LogError(err)
	}
	reqPost.Header.Set("Node-Addr", nodeInfo.String)
	postRep, err := client.Do(reqPost)
	defer postRep.Body.Close()
	if err != nil {
		return ext.LogError(err)
	}

	reqGet, err := http.NewRequest("GET", "http://127.0.0.1:20000/locateGet", nil)
	if err != nil {
		return ext.LogError(err)
	}

	resp, err := client.Do(reqGet)
	defer resp.Body.Close()
	if err != nil {
		return ext.LogError(err)
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ext.LogError(err)
		}
		var svrList []string
		err = json.Unmarshal(body, &svrList)
		if err != nil {
			return ext.LogError(err)
		}
		for index := 0; index < len(svrList); index++ {
			conn, err := net.Dial("tcp", svrList[index])
			if err != nil {
				ext.LogError(err)
			} else {
				NewRemote(mes, conn)
			}
		}
	}

	return err
}
