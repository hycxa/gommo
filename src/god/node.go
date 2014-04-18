package god

import (
	"encoding/json"
	"ext"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var (
	workTab [WORKER_NUM_LIMIT]*Worker
)

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

func NewNode(name, network, address string) {
}

func GetNodeInfo() NodeInfo {
	return nil
}

func (self *Node) ConnOtherSvr() error {
	client := &http.Client{}
	reqPost, err := http.NewRequest("POST", "http://127.0.0.1:20000/locatePost", nil)
	if err != nil {
		return ext.LogError(err)
	}
	reqPost.Header.Set("Node-Addr", self.String)
	reqPost.Header.Set("service", self.NodeType)
	postRep, err := client.Do(reqPost)
	defer postRep.Body.Close()
	if err != nil {
		return ext.LogError(err)
	}

	reqGet, err := http.NewRequest("GET", "http://127.0.0.1:20000/locateGet", nil)
	if err != nil {
		return ext.LogError(err)
	}
	getSvrTab := make([]string, 2)
	getSvrTab[0] = NODE_GS_TYPE
	if self.NodeType == NODE_GS_TYPE {
		getSvrTab[0] = NODE_GS_TYPE
		getSvrTab[1] = NODE_GATE_TYPE
	}
	for i := 0; i < len(getSvrTab); i++ {
		if getSvrTab[i] != "" {
			reqGet.Header.Set("service", getSvrTab[i])
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
					err := self.Dial("tcp", svrList[index])
					if err != nil {
						ext.LogError(err)
					}
				}
			}

		}
	}
	return err
}
