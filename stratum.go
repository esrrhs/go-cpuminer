package main

import (
	"bufio"
	"encoding/json"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"net"
	"strings"
	"time"
)

type Stratum struct {
	pool  string
	alg   *Algorithm
	user  string
	pass  string
	agent string
	rigid string
	rpcid string

	ext_algo      bool
	ext_nicehash  bool
	ext_connect   bool
	ext_keepalive bool

	conn   net.Conn
	reader *bufio.Reader
	job    *Job
}

func NewStratum(pool string, alg *Algorithm, user string, pass string) (*Stratum, error) {
	var s Stratum
	s.user = user
	s.pass = pass
	s.alg = alg
	s.pool = pool

	err := s.Reconnect()
	if err != nil {
		loggo.Error("Stratum New fail %v %v", s.pool, err)
		return nil, nil
	}

	go s.listen()

	loggo.Info("Stratum New ok")

	return &s, nil
}

func (s *Stratum) Reconnect() error {

	loggo.Info("Stratum New start Using user %v pass %v pool %v", s.user, s.pass, s.pool)

	if s.conn != nil {
		s.conn.Close()
	}

	conn, err := net.Dial("tcp", s.pool)
	if err != nil {
		loggo.Error("Stratum Dial fail %v %v", s.pool, err)
		return err
	}
	s.conn = conn

	loggo.Info("Stratum pool connect ok %v->%v", conn.LocalAddr(), conn.RemoteAddr())

	s.reader = bufio.NewReader(s.conn)

	err = s.login()
	if err != nil {
		loggo.Error("Stratum login fail %v", err)
		return err
	}
	return nil
}

func (s *Stratum) listen() {
	defer common.CrashLog()

	loggo.Info("Stratum Starting Listener")

	for {
		result, err := s.reader.ReadString('\n')
		if err != nil {
			loggo.Error("Stratum Connection lost %v", err)
			time.Sleep(time.Second)
			s.Reconnect()
			continue
		}

		loggo.Debug("Stratum recv %v", strings.TrimSuffix(result, "\n"))
		var rsp JSONRpcRsp
		err = json.Unmarshal([]byte(result), &rsp)
		if err != nil {
			loggo.Error("Stratum Unmarshal fail %v %v", result, err)
			continue
		}

		if !s.handleRsp(rsp) {
			loggo.Error("Stratum handleRsp fail %v", result)
			continue
		}
	}
}

func (s *Stratum) handleRsp(rsp JSONRpcRsp) bool {
	loggo.Debug("Stratum handleRsp %v", rsp.Id)
	err := rsp.Error
	if err != nil {
		loggo.Error("Stratum handleRsp error %v", err)
		return false
	}
	id := rsp.Id
	if id != 0 {
		return s.handleResponse(id, rsp)
	}

	return s.handleNotify(rsp)
}

func (s *Stratum) handleNotify(rsp JSONRpcRsp) bool {
	loggo.Debug("Stratum handleNotify %v", rsp.Method)

	if rsp.Method == "" {
		loggo.Error("Stratum handleNotify no method")
		return false
	}

	m := rsp.Method
	if m == "job" {
		return s.handleNotifyJob(rsp)
	}

	return true
}

func (s *Stratum) handleNotifyJob(rsp JSONRpcRsp) bool {
	loggo.Debug("Stratum handleNotifyJob %v", rsp.Method)

	if rsp.Params == nil {
		loggo.Error("Stratum handleNotifyJob no Params")
		return false
	}

	var job JobReplyData
	err := json.Unmarshal(*rsp.Params, &job)
	if err != nil {
		loggo.Error("Stratum handleNotifyJob Unmarshal fail %v", err)
		return false
	}

	return s.parseJob(&job)
}

func (s *Stratum) handleResponse(id int, rsp JSONRpcRsp) bool {
	loggo.Debug("Stratum handleResponse %v", id)
	if id == 1 {
		return s.handleLogin(rsp)
	}

	return s.handleSubmitResponse(id, rsp)
}

func (s *Stratum) handleSubmitResponse(id int, rsp JSONRpcRsp) bool {
	loggo.Info("Stratum handleSubmitResponse %v", id)
	// TODO
	return true
}

func (s *Stratum) handleLogin(rsp JSONRpcRsp) bool {
	result := rsp.Result
	if result == nil {
		loggo.Error("Stratum handleLogin no result")
		return false
	}

	loggo.Debug("Stratum handleLogin rsp")

	if result.Id == "" {
		loggo.Error("Stratum handleLogin no Id")
		return false
	}

	s.rpcid = result.Id

	if !s.parseExtensions(result) {
		loggo.Error("Stratum parseExtensions fail")
		return false
	}

	if !s.parseJob(result.Job) {
		loggo.Error("Stratum parseJob fail")
		return false
	}

	loggo.Info("Stratum handleLogin ok")

	return true
}

func (s *Stratum) parseJob(job *JobReplyData) bool {
	j := &Job{
		algorithm: s.alg,
		nicehash:  s.ext_nicehash,
		clientId:  s.rpcid,
	}

	if job.JobId == "" {
		loggo.Error("Stratum parseJob no JobId")
		return false
	}
	j.id = job.JobId

	if job.Algo != "" {
		j.algorithm = NewAlgorithm(job.Algo)
		if j.algorithm == nil {
			loggo.Error("Stratum parseJob fail Algorithm %v", job.Algo)
			return false
		}
	} else {
		if j.algorithm == nil {
			loggo.Error("Stratum no default Algorithm")
			return false
		}
	}

	if !j.setBlob(job.Blob) {
		loggo.Error("Stratum parseJob fail Blob %v", job.Blob)
		return false
	}

	if !j.setTarget(job.Target) {
		loggo.Error("Stratum parseJob fail Target %v", job.Target)
		return false
	}

	j.height = job.Height

	if j.algorithm.family() == RANDOM_X {
		if !j.setSeedHash(job.SeedHash) {
			loggo.Error("Stratum parseJob fail SeedHash %v", job.SeedHash)
			return false
		}
	}

	s.job = j

	loggo.Info("Stratum parseJob ok id=%v algo=%v height=%v target=%v diff=%v", j.id, j.algorithm.name(), j.height, j.target, j.diff)

	return true
}

func (s *Stratum) parseExtensions(result *JobReply) bool {
	for _, name := range result.Extensions {
		if name == "algo" {
			s.ext_algo = true
		} else if name == "nicehash" {
			s.ext_nicehash = true
		} else if name == "connect" {
			s.ext_connect = true
		} else if name == "keepalive" {
			s.ext_keepalive = true
		} else {
			loggo.Info("Stratum parseExtensions unknow %v", name)
		}
	}
	return true
}

func (s *Stratum) send(id int, method string, p interface{}) error {

	m, err := json.Marshal(p)
	if err != nil {
		loggo.Error("Stratum send Marshal fail %v", err)
		return err
	}

	req := JSONRpcReq{
		Id:      id,
		Method:  method,
		JsonRPC: "2.0",
		Params:  (*json.RawMessage)(&m),
	}

	reqm, err := json.Marshal(&req)
	if err != nil {
		loggo.Error("Stratum send Marshal fail %v", err)
		return err
	}

	_, err = s.conn.Write(reqm)
	if err != nil {
		loggo.Error("Stratum send Write fail %v", err)
		return err
	}

	_, err = s.conn.Write([]byte("\n"))
	if err != nil {
		loggo.Error("Stratum send Write fail %v", err)
		return err
	}

	return nil
}

func (s *Stratum) login() error {
	msg := LoginParam{
		Login: s.user,
		Pass:  s.pass,
		Agent: s.agent,
		Rigid: s.rigid,
	}

	loggo.Info("Stratum start login...")

	return s.send(1, "login", &msg)
}
