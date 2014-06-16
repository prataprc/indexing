package adminport

import (
	"encoding/json"
	"github.com/couchbase/indexing/secondary/common"
	"reflect"
	"testing"
)

var addr = "http://localhost:9999"

type testMessage struct {
	DefnID          uint64 `json:"defnId"`
	Bucket          string `json:"bucket"`
	IsPrimary       bool   `json:"isPrimary"`
	IName           string `json:"name"`
	Using           string `json:"using"`
	ExprType        string `json:"exprType"`
	PartitionScheme string `json:"partitionType"`
	Expression      string `json:"expression"`
}

func TestLoopback(t *testing.T) {
	//common.SetLogLevel(common.LogLevelTrace)
	q := make(chan bool)

	doServer(addr, t, q)

	client := NewHTTPClient(addr, common.AdminportURLPrefix)
	req := &testMessage{
		DefnID:          uint64(0x1234567812345678),
		Bucket:          "default",
		IsPrimary:       false,
		IName:           "example-index",
		Using:           "forrestdb",
		ExprType:        "n1ql",
		PartitionScheme: "simplekeypartition",
		Expression:      "x+1",
	}
	resp := &testMessage{}
	if err := client.Request(req, resp); err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(req, resp) == false {
		t.Error("unexpected response")
	}
	stats := common.ComponentStat{}
	if err := client.RequestStat("/adminport", &stats); err != nil {
		t.Error(err)
	}
	common.Infof("%v\n", stats)
	if stats["requests"].(float64) != float64(2) {
		t.Error("registered requests", stats["requests"])
	}
}

func BenchmarkClientRequest(b *testing.B) {
	client := NewHTTPClient(addr, common.AdminportURLPrefix)
	req := &testMessage{
		DefnID:          uint64(0x1234567812345678),
		Bucket:          "default",
		IsPrimary:       false,
		IName:           "example-index",
		Using:           "forrestdb",
		ExprType:        "n1ql",
		PartitionScheme: "simplekeypartition",
		Expression:      "x+1",
	}
	resp := &testMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := client.Request(req, resp); err != nil {
			b.Error(err)
		}
	}
}

func doServer(addr string, tb testing.TB, quit chan bool) Server {
	urlPrefix, reqch := common.AdminportURLPrefix, make(chan Request, 10)
	server := NewHTTPServer("test", "localhost:9999", urlPrefix, reqch)
	if err := server.Register(&testMessage{}); err != nil {
		tb.Fatal(err)
	}
	if err := server.Register(&common.ComponentStat{}); err != nil {
		tb.Fatal(err)
	}

	if err := server.Start(); err != nil {
		tb.Fatal(err)
	}

	go func() {
	loop:
		for {
			select {
			case req, ok := <-reqch:
				if ok {
					switch msg := req.GetMessage().(type) {
					case *testMessage:
						if err := req.Send(msg); err != nil {
							tb.Error(err)
						}
					case *common.ComponentStat:
						m := &common.ComponentStat{
							"adminport": server.GetStatistics(),
						}
						if err := req.Send(m); err != nil {
							tb.Error(err)
						}
					}
				} else {
					break loop
				}
			}
		}
		close(quit)
	}()

	return server
}

func (tm *testMessage) Name() string {
	return "testMessage"
}

func (tm *testMessage) Encode() (data []byte, err error) {
	data, err = json.Marshal(tm)
	return
}

func (tm *testMessage) Decode(data []byte) (err error) {
	err = json.Unmarshal(data, tm)
	return
}

func (tm *testMessage) ContentType() string {
	return "application/json"
}
