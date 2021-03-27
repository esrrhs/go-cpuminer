package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/wire"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"io"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// ErrStratumStaleWork indicates that the work to send to the pool was stale.
var ErrStratumStaleWork = fmt.Errorf("Stale work, throwing away")
var version = "1.0.0"

// These variables are the chain proof-of-work limit parameters for each default
// network.
var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// mainPowLimit is the highest proof of work value a Decred block can
	// have for the main network.  It is the value 2^224 - 1.
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
)

// Stratum holds all the shared information for a stratum connection.
type Stratum struct {
	// The following variables must only be used atomically.
	ValidShares   uint64
	InvalidShares uint64
	latestJobTime uint32

	cfg       Config
	Conn      net.Conn
	Reader    *bufio.Reader
	ID        uint64
	authID    uint64
	subID     uint64
	submitIDs []uint64
	Diff      float64
	Target    *big.Int
	PoolWork  NotifyWork

	Started uint32
}

// Config holdes the config options that may be used by a stratum pool.
type Config struct {
	Pool string
	User string
	Pass string
}

// NotifyWork holds all the info recieved from a mining.notify message along
// with the Work data generate from it.
type NotifyWork struct {
	Clean             bool
	ExtraNonce1       string
	ExtraNonce2       uint64
	ExtraNonce2Length float64
	Nonce2            uint32
	CB1               string
	CB2               string
	Height            int64
	NtimeDelta        int64
	JobID             string
	Hash              string
	Nbits             string
	Ntime             string
	Version           string
	NewWork           bool
	Work              *Work
}

// StratumMsg is the basic message object from stratum.
type StratumMsg struct {
	Method string `json:"method"`
	// Need to make generic.
	Params []string    `json:"params"`
	ID     interface{} `json:"id"`
}

// StratumRsp is the basic response type from stratum.
type StratumRsp struct {
	Method string `json:"method"`
	// Need to make generic.
	ID     interface{}      `json:"id"`
	Error  StratErr         `json:"error,omitempty"`
	Result *json.RawMessage `json:"result,omitempty"`
}

// StratErr is the basic error type (a number and a string) sent by
// the stratum server.
type StratErr struct {
	ErrNum uint64
	ErrStr string
	Result *json.RawMessage `json:"result,omitempty"`
}

// Basic reply is a reply type for any of the simple messages.
type BasicReply struct {
	ID     interface{} `json:"id"`
	Error  StratErr    `json:"error,omitempty"`
	Result bool        `json:"result"`
}

// SubscribeReply models the server response to a subscribe message.
type SubscribeReply struct {
	SubscribeID       string
	ExtraNonce1       string
	ExtraNonce2Length float64
}

// NotifyRes models the json from a mining.notify message.
type NotifyRes struct {
	JobID          string
	Hash           string
	GenTX1         string
	GenTX2         string
	MerkleBranches []string
	BlockVersion   string
	Nbits          string
	Ntime          string
	CleanJobs      bool
}

// Submit models a submission message.
type Submit struct {
	Params []string    `json:"params"`
	ID     interface{} `json:"id"`
	Method string      `json:"method"`
}

// errJsonType is an error for json that we do not expect.
var errJsonType = errors.New("Unexpected type in json.")

func sliceContains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sliceRemove(s []uint64, e uint64) []uint64 {
	for i, a := range s {
		if a == e {
			return append(s[:i], s[i+1:]...)
		}
	}

	return s
}

// StratumConn starts the initial connection to a stratum pool and sets defaults
// in the pool object.
func StratumConn(pool, user, pass string) (*Stratum, error) {
	var stratum Stratum
	stratum.cfg.User = user
	stratum.cfg.Pass = pass

	loggo.Info("Stratum New start Using user %v pass %v pool %v", user, pass, pool)

	var conn net.Conn
	var err error

	conn, err = net.Dial("tcp", pool)
	if err != nil {
		return nil, err
	}
	stratum.ID = 1
	stratum.Conn = conn
	stratum.cfg.Pool = pool

	loggo.Info("Stratum pool connect ok %v->%v", conn.LocalAddr(), conn.RemoteAddr())

	// We will set it for sure later but this really should be the value and
	// setting it here will prevent so incorrect matches based on the
	// default 0 value.
	stratum.authID = 2

	// Target for share is 1 unless we hear otherwise.
	stratum.Diff = 1
	stratum.Target, err = DiffToTarget(stratum.Diff, mainPowLimit)
	if err != nil {
		return nil, err
	}
	stratum.PoolWork.NewWork = false
	stratum.Reader = bufio.NewReader(stratum.Conn)
	go stratum.Listen()

	err = stratum.Subscribe()
	if err != nil {
		return nil, err
	}
	err = stratum.Auth()
	if err != nil {
		return nil, err
	}

	stratum.Started = uint32(time.Now().Unix())

	loggo.Info("Stratum New ok")

	return &stratum, nil
}

// Reconnect reconnects to a stratum server if the connection has been lost.
func (s *Stratum) Reconnect() error {
	var conn net.Conn
	var err error

	conn, err = net.Dial("tcp", s.cfg.Pool)
	if err != nil {
		return err
	}
	s.Conn = conn
	s.Reader = bufio.NewReader(s.Conn)
	err = s.Subscribe()
	if err != nil {
		return nil
	}
	// Should NOT need this.
	time.Sleep(5 * time.Second)
	err = s.Auth()
	if err != nil {
		return nil
	}

	// If we were able to reconnect, restart counter
	s.Started = uint32(time.Now().Unix())

	return nil
}

// Listen is the listener for the incoming messages from the stratum pool.
func (s *Stratum) Listen() {
	defer common.CrashLog()

	loggo.Info("Stratum Starting Listener")

	for {
		result, err := s.Reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				loggo.Error("Stratum Connection lost!  Reconnecting.")
				err = s.Reconnect()
				if err != nil {
					loggo.Error("Stratum Reconnect failed. %v", err)
					time.Sleep(5 * time.Second)
				}
			} else {
				loggo.Error("Stratum ReadString fail %v", err)
			}
			continue
		}

		loggo.Info("Stratum recv %v", strings.TrimSuffix(result, "\n"))
		resp, err := s.Unmarshal([]byte(result))
		if err != nil {
			loggo.Error("Stratum Unmarshal fail %v", err)
			continue
		}
		loggo.Info("Stratum recv \n", common.StructToTable(&resp))

		switch resp.(type) {
		case *BasicReply:
			s.handleBasicReply(resp)
		case StratumMsg:
			s.handleStratumMsg(resp)
		case NotifyRes:
			s.handleNotifyRes(resp)
		case *SubscribeReply:
			s.handleSubscribeReply(resp)
		default:
			loggo.Info("Stratum Unhandled message: ", result)
		}
	}
}

func (s *Stratum) handleBasicReply(resp interface{}) {
	aResp := resp.(*BasicReply)

	if int(aResp.ID.(uint64)) == int(s.authID) {
		if aResp.Result {
			loggo.Info("Stratum handleBasicReply Logged in")
		} else {
			loggo.Error("Stratum handleBasicReply Auth failure.")
		}
	}
	if sliceContains(s.submitIDs, aResp.ID.(uint64)) {
		if aResp.Result {
			atomic.AddUint64(&s.ValidShares, 1)
			loggo.Info("Stratum handleBasicReply Share accepted")
		} else {
			atomic.AddUint64(&s.InvalidShares, 1)
			loggo.Error("Stratum handleBasicReply Share rejected ", aResp.Error.ErrStr)
		}
		s.submitIDs = sliceRemove(s.submitIDs, aResp.ID.(uint64))
	}
}

func (s *Stratum) handleStratumMsg(resp interface{}) {
	nResp := resp.(StratumMsg)
	loggo.Info("Stratum handleStratumMsg %v %v", nResp.Method, nResp.ID)
	// Too much is still handled in unmarshaler.  Need to
	// move stuff other than unmarshalling here.
	switch nResp.Method {
	case "client.show_message":
		loggo.Info("Stratum handleStratumMsg show_message: %v", nResp.Params)

	case "client.reconnect":
		loggo.Debug("Reconnect requested")
		wait, err := strconv.Atoi(nResp.Params[2])
		if err != nil {
			loggo.Error("Stratum handleStratumMsg reconnect: %v", err)
			return
		}
		time.Sleep(time.Duration(wait) * time.Second)
		pool := nResp.Params[0] + ":" + nResp.Params[1]
		s.cfg.Pool = pool
		err = s.Reconnect()
		if err != nil {
			loggo.Error("Stratum handleStratumMsg reconnect: %v", err)
			return
		}

	case "client.get_version":
		loggo.Debug("get_version request received.")
		msg := StratumMsg{
			Method: nResp.Method,
			ID:     nResp.ID,
			Params: []string{"go-cpuminer/" + version},
		}
		m, err := json.Marshal(msg)
		if err != nil {
			loggo.Error("Stratum handleStratumMsg get_version: %v", err)
			return
		}
		_, err = s.Conn.Write(m)
		if err != nil {
			loggo.Error("Stratum handleStratumMsg get_version: %v", err)
			return
		}
		_, err = s.Conn.Write([]byte("\n"))
		if err != nil {
			loggo.Error("Stratum handleStratumMsg get_version: %v", err)
			return
		}
	}
}

func (s *Stratum) handleNotifyRes(resp interface{}) {
	nResp := resp.(NotifyRes)
	s.PoolWork.JobID = nResp.JobID
	s.PoolWork.CB1 = nResp.GenTX1
	heightHex := nResp.GenTX1[186:188] + nResp.GenTX1[184:186]
	height, err := strconv.ParseInt(heightHex, 16, 32)
	if err != nil {
		loggo.Error("Stratum handleNotifyRes failed to parse height %v", err)
		height = 0
		return
	}

	s.PoolWork.Height = height
	s.PoolWork.CB2 = nResp.GenTX2
	s.PoolWork.Hash = nResp.Hash
	s.PoolWork.Nbits = nResp.Nbits
	s.PoolWork.Version = nResp.BlockVersion
	parsedNtime, err := strconv.ParseInt(nResp.Ntime, 16, 64)
	if err != nil {
		loggo.Error("Stratum handleNotifyRes failed to parse height %v", err)
		return
	}

	s.PoolWork.Ntime = nResp.Ntime
	s.PoolWork.NtimeDelta = parsedNtime - time.Now().Unix()
	s.PoolWork.Clean = nResp.CleanJobs
	s.PoolWork.NewWork = true
}

func (s *Stratum) handleSubscribeReply(resp interface{}) {
	nResp := resp.(*SubscribeReply)
	s.PoolWork.ExtraNonce1 = nResp.ExtraNonce1
	s.PoolWork.ExtraNonce2Length = nResp.ExtraNonce2Length
	loggo.Debug("Stratum handleSubscribeReply Subscribe reply received")
}

// Auth sends a message to the pool to authorize a worker.
func (s *Stratum) Auth() error {
	msg := StratumMsg{
		Method: "mining.authorize",
		ID:     s.ID,
		Params: []string{s.cfg.User, s.cfg.Pass},
	}
	// Auth reply has no method so need a way to identify it.
	// Ugly, but not much choice.
	id, ok := msg.ID.(uint64)
	if !ok {
		return errJsonType
	}
	s.authID = id
	s.ID++
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = s.Conn.Write(m)
	if err != nil {
		return err
	}
	_, err = s.Conn.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}

// Subscribe sends the subscribe message to get mining info for a worker.
func (s *Stratum) Subscribe() error {
	msg := StratumMsg{
		Method: "mining.subscribe",
		ID:     s.ID,
		Params: []string{"go-cpuminer/" + version},
	}
	s.subID = msg.ID.(uint64)
	s.ID++
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = s.Conn.Write(m)
	if err != nil {
		return err
	}
	_, err = s.Conn.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal provides a json unmarshaler for the commands.
// I'm sure a lot of this can be generalized but the json we deal with
// is pretty yucky.
func (s *Stratum) Unmarshal(blob []byte) (interface{}, error) {
	var (
		objmap map[string]json.RawMessage
		method string
		id     uint64
	)

	err := json.Unmarshal(blob, &objmap)
	if err != nil {
		return nil, err
	}
	// decode command
	// Not everyone has a method.
	err = json.Unmarshal(objmap["method"], &method)
	if err != nil {
		method = ""
	}
	err = json.Unmarshal(objmap["id"], &id)
	if err != nil {
		return nil, err
	}
	if id == s.authID {
		var (
			objmap      map[string]json.RawMessage
			id          uint64
			result      bool
			errorHolder []interface{}
		)
		err := json.Unmarshal(blob, &objmap)
		if err != nil {
			return nil, err
		}
		resp := &BasicReply{}

		err = json.Unmarshal(objmap["id"], &id)
		if err != nil {
			return nil, err
		}
		resp.ID = id

		err = json.Unmarshal(objmap["result"], &result)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(objmap["error"], &errorHolder)
		if err != nil {
			return nil, err
		}
		resp.Result = result

		if errorHolder != nil {
			errN, ok := errorHolder[0].(float64)
			if !ok {
				return nil, errJsonType
			}
			errS, ok := errorHolder[1].(string)
			if !ok {
				return nil, errJsonType
			}
			resp.Error.ErrNum = uint64(errN)
			resp.Error.ErrStr = errS
		}

		return resp, nil

	}
	if id == s.subID {
		var resi []interface{}
		err := json.Unmarshal(objmap["result"], &resi)
		if err != nil {
			return nil, err
		}
		resp := &SubscribeReply{}

		var objmap2 map[string]json.RawMessage
		err = json.Unmarshal(blob, &objmap2)
		if err != nil {
			return nil, err
		}

		var resJS []json.RawMessage
		err = json.Unmarshal(objmap["result"], &resJS)
		if err != nil {
			return nil, err
		}

		if len(resJS) == 0 {
			return nil, errJsonType
		}

		var msgPeak []interface{}
		err = json.Unmarshal(resJS[0], &msgPeak)
		if err != nil {
			return nil, err
		}

		if msgPeak[0] != nil {
			// The pools do not all agree on what this message looks like
			// so we need to actually look at it before unmarshalling for
			// real so we can use the right form.  Yuck.
			if msgPeak[0] == "mining.notify" {
				var innerMsg []string
				err = json.Unmarshal(resJS[0], &innerMsg)
				if err != nil {
					return nil, err
				}
				resp.SubscribeID = innerMsg[1]
			} else {
				var innerMsg [][]string
				err = json.Unmarshal(resJS[0], &innerMsg)
				if err != nil {
					return nil, err
				}

				for i := 0; i < len(innerMsg); i++ {
					if innerMsg[i][0] == "mining.notify" {
						resp.SubscribeID = innerMsg[i][1]
					}
					if innerMsg[i][0] == "mining.set_difficulty" {
						// Not all pools correctly put something
						// in here so we will ignore it (we
						// already have the default value of 1
						// anyway and pool can send a new one.
						// dcr.coinmine.pl puts something that
						// is not a difficulty here which is why
						// we ignore.
					}
				}
			}
		}

		resp.ExtraNonce1 = resi[1].(string)
		resp.ExtraNonce2Length = resi[2].(float64)
		return resp, nil
	}
	if sliceContains(s.submitIDs, id) {
		var (
			objmap      map[string]json.RawMessage
			id          uint64
			result      bool
			errorHolder []interface{}
		)
		err := json.Unmarshal(blob, &objmap)
		if err != nil {
			return nil, err
		}
		resp := &BasicReply{}

		err = json.Unmarshal(objmap["id"], &id)
		if err != nil {
			return nil, err
		}
		resp.ID = id

		err = json.Unmarshal(objmap["result"], &result)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(objmap["error"], &errorHolder)
		if err != nil {
			return nil, err
		}
		resp.Result = result

		if errorHolder != nil {
			errN, ok := errorHolder[0].(float64)
			if !ok {
				return nil, errJsonType
			}
			errS, ok := errorHolder[1].(string)
			if !ok {
				return nil, errJsonType
			}
			resp.Error.ErrNum = uint64(errN)
			resp.Error.ErrStr = errS
		}

		return resp, nil
	}
	switch method {
	case "mining.notify":
		var resi []interface{}
		err := json.Unmarshal(objmap["params"], &resi)
		if err != nil {
			return nil, err
		}
		var nres = NotifyRes{}
		jobID, ok := resi[0].(string)
		if !ok {
			return nil, errJsonType
		}
		nres.JobID = jobID
		hash, ok := resi[1].(string)
		if !ok {
			return nil, errJsonType
		}
		nres.Hash = hash
		genTX1, ok := resi[2].(string)
		if !ok {
			return nil, errJsonType
		}
		nres.GenTX1 = genTX1
		genTX2, ok := resi[3].(string)
		if !ok {
			return nil, errJsonType
		}
		nres.GenTX2 = genTX2
		//ccminer code also confirms this
		//nres.MerkleBranches = resi[4].([]string)
		blockVersion, ok := resi[5].(string)
		if !ok {
			return nil, errJsonType
		}
		nres.BlockVersion = blockVersion
		nbits, ok := resi[6].(string)
		if !ok {
			return nil, errJsonType
		}
		nres.Nbits = nbits
		ntime, ok := resi[7].(string)
		if !ok {
			return nil, errJsonType
		}
		nres.Ntime = ntime
		cleanJobs, ok := resi[8].(bool)
		if !ok {
			return nil, errJsonType
		}
		nres.CleanJobs = cleanJobs
		return nres, nil

	case "mining.set_difficulty":
		var resi []interface{}
		err := json.Unmarshal(objmap["params"], &resi)
		if err != nil {
			return nil, err
		}

		difficulty, ok := resi[0].(float64)
		if !ok {
			return nil, errJsonType
		}
		s.Target, err = DiffToTarget(difficulty, mainPowLimit)
		if err != nil {
			return nil, err
		}
		s.Diff = difficulty
		var nres = StratumMsg{}
		nres.Method = method
		diffStr := strconv.FormatFloat(difficulty, 'E', -1, 32)
		var params []string
		params = append(params, diffStr)
		nres.Params = params
		loggo.Info("Stratum difficulty set to %v", difficulty)
		return nres, nil

	case "client.show_message":
		var resi []interface{}
		err := json.Unmarshal(objmap["result"], &resi)
		if err != nil {
			return nil, err
		}
		msg, ok := resi[0].(string)
		if !ok {
			return nil, errJsonType
		}
		var nres = StratumMsg{}
		nres.Method = method
		var params []string
		params = append(params, msg)
		nres.Params = params
		return nres, nil

	case "client.get_version":
		var nres = StratumMsg{}
		var id uint64
		err = json.Unmarshal(objmap["id"], &id)
		if err != nil {
			return nil, err
		}
		nres.Method = method
		nres.ID = id
		return nres, nil

	case "client.reconnect":
		var nres = StratumMsg{}
		var id uint64
		err = json.Unmarshal(objmap["id"], &id)
		if err != nil {
			return nil, err
		}
		nres.Method = method
		nres.ID = id

		var resi []interface{}
		err := json.Unmarshal(objmap["params"], &resi)
		if err != nil {
			return nil, err
		}

		if len(resi) < 3 {
			return nil, errJsonType
		}
		hostname, ok := resi[0].(string)
		if !ok {
			return nil, errJsonType
		}
		p, ok := resi[1].(float64)
		if !ok {
			return nil, errJsonType
		}
		port := strconv.Itoa(int(p))
		w, ok := resi[2].(float64)
		if !ok {
			return nil, errJsonType
		}
		wait := strconv.Itoa(int(w))

		nres.Params = []string{hostname, port, wait}

		return nres, nil

	default:
		resp := &StratumRsp{}
		err := json.Unmarshal(blob, &resp)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// PrepWork converts the stratum notify to getwork style data for mining.
func (s *Stratum) PrepWork() error {
	// Build final extranonce, which is basically the pool user and worker
	// ID.
	en1, err := hex.DecodeString(s.PoolWork.ExtraNonce1)
	if err != nil {
		loggo.Error("Stratum Error decoding ExtraNonce1 %v", err)
		return err
	}

	// Work out padding.
	tmp := []string{"%0", strconv.Itoa(int(s.PoolWork.ExtraNonce2Length) * 2), "x"}
	fmtString := strings.Join(tmp, "")
	en2, err := hex.DecodeString(fmt.Sprintf(fmtString, s.PoolWork.ExtraNonce2))
	if err != nil {
		loggo.Error("Stratum Error decoding ExtraNonce2 %v", err)
		return err
	}
	extraNonce := append(en1[:], en2[:]...)

	// Put coinbase transaction together.
	cb1, err := hex.DecodeString(s.PoolWork.CB1)
	if err != nil {
		loggo.Error("Stratum Error decoding Coinbase pt 1 %v", err)
		return err
	}
	cb2, err := hex.DecodeString(s.PoolWork.CB2)
	if err != nil {
		loggo.Error("Stratum Error decoding Coinbase pt 2 %v", err)
		return err
	}

	// Serialize header.
	bh := wire.BlockHeader{}
	v, err := ReverseToInt(s.PoolWork.Version)
	if err != nil {
		return err
	}
	bh.Version = v

	nbits, err := hex.DecodeString(s.PoolWork.Nbits)
	if err != nil {
		loggo.Error("Stratum Error decoding nbits %v", err)
		return err
	}

	b, _ := binary.Uvarint(nbits)
	bh.Bits = uint32(b)
	t := time.Now().Unix() + s.PoolWork.NtimeDelta
	bh.Timestamp = time.Unix(t, 0)
	bh.Nonce = 0

	// Serialized version.
	blockHeader, err := bh.Bytes()
	if err != nil {
		return err
	}

	data := blockHeader
	copy(data[31:139], cb1[0:108])

	var workdata [180]byte
	workPosition := 0

	version := new(bytes.Buffer)
	err = binary.Write(version, binary.LittleEndian, v)
	if err != nil {
		return err
	}
	copy(workdata[workPosition:], version.Bytes())

	prevHash := RevHash(s.PoolWork.Hash)
	p, err := hex.DecodeString(prevHash)
	if err != nil {
		loggo.Error("Stratum Error encoding previous hash %v", err)
		return err
	}

	workPosition += 4
	copy(workdata[workPosition:], p)
	workPosition += 32
	copy(workdata[workPosition:], cb1[0:108])
	workPosition += 108
	copy(workdata[workPosition:], extraNonce)
	workPosition = 176
	copy(workdata[workPosition:], cb2)

	var randomBytes = make([]byte, 4)
	_, err = rand.Read(randomBytes)
	if err != nil {
		loggo.Error("Stratum Unable to generate random bytes %v", err)
		return err
	}

	var workData [192]byte
	copy(workData[:], workdata[:])
	givenTs := binary.LittleEndian.Uint32(
		workData[128+4*TimestampWord : 132+4*TimestampWord])
	atomic.StoreUint32(&s.latestJobTime, givenTs)

	if s.Target == nil {
		loggo.Error("Stratum No target set! Reconnecting to pool.")
		err = s.Reconnect()
		if err != nil {
			loggo.Error("Stratum Reconnect failed. %v", err)
			return err
		}
		return nil
	}

	w := NewWork(workData, s.Target, givenTs, uint32(time.Now().Unix()), false)

	loggo.Info("Stratum prepated work data %v, target %032x", hex.EncodeToString(w.Data[:]), w.Target.Bytes())
	s.PoolWork.Work = w

	return nil
}

// PrepSubmit formats a mining.sumbit message from the solved work.
func (s *Stratum) PrepSubmit(data []byte) (Submit, error) {
	loggo.Debug("Stratum got valid work to submit %x", data)
	loggo.Debug("Stratum got valid work hash %v", chainhash.HashH(data[0:180]))
	data2 := make([]byte, 180)
	copy(data2, data[0:180])

	sub := Submit{}
	sub.Method = "mining.submit"

	// Format data to send off.
	hexData := hex.EncodeToString(data)
	decodedData, err := hex.DecodeString(hexData)
	if err != nil {
		loggo.Error("Stratum Error decoding data %v", err)
		return sub, err
	}

	var submittedHeader wire.BlockHeader
	bhBuf := bytes.NewReader(decodedData[0:wire.MaxBlockHeaderPayload])
	err = submittedHeader.Deserialize(bhBuf)
	if err != nil {
		loggo.Error("Error generating header %v", err)
		return sub, err
	}

	latestWorkTs := atomic.LoadUint32(&s.latestJobTime)
	if uint32(submittedHeader.Timestamp.Unix()) != latestWorkTs {
		return sub, ErrStratumStaleWork
	}

	s.ID++
	sub.ID = s.ID
	s.submitIDs = append(s.submitIDs, s.ID)

	// The timestamp string should be:
	//
	//   timestampStr := fmt.Sprintf("%08x",
	//     uint32(submittedHeader.Timestamp.Unix()))
	//
	// but the "stratum" protocol appears to only use this value
	// to check if the miner is in sync with the latest announcement
	// of work from the pool. If this value is anything other than
	// the timestamp of the latest pool work timestamp, work gets
	// rejected from the current implementation.
	timestampStr := fmt.Sprintf("%08x", latestWorkTs)
	nonceStr := fmt.Sprintf("%08x", submittedHeader.Nonce)
	xnonceStr := hex.EncodeToString(data[144:156])

	sub.Params = []string{s.cfg.User, s.PoolWork.JobID, xnonceStr, timestampStr, nonceStr}

	return sub, nil
}
