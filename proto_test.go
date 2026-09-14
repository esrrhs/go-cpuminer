package main

import (
	"encoding/json"
	"testing"
)

func TestJSONRpcReqMarshal(t *testing.T) {
	req := JSONRpcReq{
		Id:      1,
		Method:  "login",
		JsonRPC: "2.0",
		Params: LoginParam{
			Login: "wallet",
			Pass:  "x",
			Agent: "go-cpuminer",
		},
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var decoded JSONRpcReq
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Method != "login" || decoded.Id != 1 {
		t.Fatalf("decoded = %+v", decoded)
	}
}

func TestJobReplyUnmarshal(t *testing.T) {
	raw := []byte(`{"id":1,"jsonrpc":"2.0","result":{"id":"worker","job":{"blob":"aa","job_id":"1","target":"ffff0000","algo":"cn-heavy/xhv","height":10},"status":"OK"}}`)
	var rsp JSONRpcRsp
	if err := json.Unmarshal(raw, &rsp); err != nil {
		t.Fatal(err)
	}
	if rsp.Result == nil || rsp.Result.Job == nil {
		t.Fatal("missing result job")
	}
	if rsp.Result.Job.Algo != "cn-heavy/xhv" || rsp.Result.Job.Height != 10 {
		t.Fatalf("job = %+v", rsp.Result.Job)
	}
}
