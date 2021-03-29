package main

import "encoding/json"

type JSONRpcReq struct {
	Id      int         `json:"id"`
	Method  string      `json:"method"`
	JsonRPC string      `json:"jsonrpc"`
	Params  interface{} `json:"params"`
}

type LoginParam struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
	Agent string `json:"agent"`
	Rigid string `json:"rigid"`
}

///////////////////////////////////////////////////

type JSONRpcRsp struct {
	Id      int              `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  *JobReply        `json:"result"`
	Error   *ErrorReply      `json:"error"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
}

type ErrorReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JobReply struct {
	Id         string        `json:"id"`
	Job        *JobReplyData `json:"job"`
	Extensions []string      `json:"extensions"`
	Status     string        `json:"status"`
}

type JobReplyData struct {
	Blob     string `json:"blob"`
	JobId    string `json:"job_id"`
	Target   string `json:"target"`
	Algo     string `json:"algo"`
	Height   uint64 `json:"height"`
	SeedHash string `json:"seed_hash"`
}
