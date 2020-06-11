package mq

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"encoding/json"
	"fmt"
	"time"
)

type Accessor struct {
	mqConf *domain.MQConf
}

type AccessorLog struct {
	Role     string        `json:"role"`
	ReqId    uint64        `json:"req_id"`
	RespId   uint64        `json:"resp_id"`
	ReqTube  string        `json:"req_tube"`
	RespTube string        `json:"resp_tube"`
	HandTime time.Duration `json:"hand_time"`
	Error    string        `json:"error"`
}

func (log *AccessorLog) Output() {
	log.Role = "accessor"
	v, _ := json.Marshal(log)
	loggers.Info.Println(string(v))
}

type HandleFunc func(Message) Message

func NewAcessor(mqConf *domain.MQConf, timeout time.Duration) (*Accessor, error) {
	return &Accessor{mqConf}, nil
}

func (accessor *Accessor) Register(msg Message, f HandleFunc, timeout time.Duration) {
	addr := fmt.Sprintf("%s:%d", accessor.mqConf.Host, accessor.mqConf.Port)
	conn, err := NewConn(accessor.mqConf.Network, addr)
	if err != nil {
		panic(err)
	}
	go accessor.serveMsg(conn, msg, f, timeout)
}

func (accessor *Accessor) serveMsg(conn *Conn, msg Message, f HandleFunc, timeout time.Duration) {
	comsumer := conn.NewComsumer(msg.TubeName())
	for {
		req := msg.New()
		if err := comsumer.Reserve(timeout, req); err != nil {
			loggers.Error.Printf("TubeName:%s reserve msg error:%s", msg.TubeName(), err.Error())
			continue
		}
		start := time.Now().UTC()
		reqID := req.GetHeader().ID
		comsumer.Delete(reqID)
		resp := f(req)
		producer := conn.NewProducer(resp.TubeName())
		resp.SetReqID(reqID)
		log := AccessorLog{
			ReqId:    reqID,
			ReqTube:  req.TubeName(),
			RespTube: resp.TubeName(),
			HandTime: time.Since(start),
		}
		id, err := producer.Put(resp)
		if err != nil {
			log.Error = err.Error()
		} else {
			log.RespId = id
		}
		log.Output()
	}
}
