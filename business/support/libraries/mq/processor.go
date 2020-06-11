package mq

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Resp struct {
	Ch  chan bool
	Err error
	Msg Message
}

type ProcessorLog struct {
	Role     string        `json:"role"`
	ReqId    uint64        `json:"req_id"`
	RespId   uint64        `json:"resp_id"`
	ReqTube  string        `json:"req_tube"`
	RespTube string        `json:"resp_tube"`
	HandTime time.Duration `json:"hand_time"`
	Error    string        `json:"error"`
}

func (log *ProcessorLog) Output() {
	log.Role = "processor"
	v, _ := json.Marshal(log)
	loggers.Info.Println(string(v))
}

type Processor struct {
	conn   *Conn
	mqConf *domain.MQConf
	mutex  sync.Mutex
	blocks map[uint64]Resp
}

func NewProcessor(mqConf *domain.MQConf) (*Processor, error) {
	addr := fmt.Sprintf("%s:%d", mqConf.Host, mqConf.Port)
	conn, err := NewConn(mqConf.Network, addr)
	if err != nil {
		return nil, err
	}

	processor := &Processor{
		conn:   conn,
		mqConf: mqConf,
		blocks: make(map[uint64]Resp),
	}
	return processor, nil
}

func (processor *Processor) Register(msg Message, timeout time.Duration) error {
	addr := fmt.Sprintf("%s:%d", processor.mqConf.Host, processor.mqConf.Port)
	conn, err := NewConn(processor.mqConf.Network, addr)
	if err != nil {
		return err
	}
	go processor.catchResp(conn, msg, timeout)
	return nil
}

func (processor *Processor) catchResp(conn *Conn, msg Message, timeout time.Duration) error {
	comsumer := conn.NewComsumer(msg.TubeName())
	loggers.Error.Printf("TubeName:%s cacth resp start", msg.TubeName())
	for {
		resp := msg
		if err := comsumer.Reserve(timeout, resp); err != nil {
			//loggers.Error.Printf("TubeName:%s reserve msg error:%s", resp.TubeName(), err.Error())
			continue
		}
		comsumer.Delete(resp.GetHeader().ID)
		reqID := resp.ReqID()
		processor.mutex.Lock()
		r, ok := processor.blocks[reqID]
		if !ok {
			loggers.Error.Printf("TubeName:%s reserve id:%d no cactch", resp.TubeName(), reqID)
			processor.mutex.Unlock()
			continue
		}
		loggers.Debug.Printf("%#v", r)
		r.Err = nil
		r.Msg = resp
		processor.blocks[reqID] = r
		r.Ch <- true
		processor.mutex.Unlock()
	}
}

func (processor *Processor) Put(req Message, timeout time.Duration) (Message, error) {
	start := time.Now().UTC()
	log := ProcessorLog{
		ReqTube: req.TubeName(),
	}
	defer log.Output()
	loggers.Debug.Printf("%v \n", req)
	loggers.Debug.Println(processor.conn)
	producer := processor.conn.NewProducer(req.TubeName())
	id, err := producer.Put(req)
	log.ReqId = id
	if err != nil {
		log.Error = err.Error()
		return nil, err
	}
	ch := make(chan bool)
	processor.mutex.Lock()
	processor.blocks[id] = Resp{
		Ch: ch,
	}
	processor.mutex.Unlock()

	select {
	case <-time.After(timeout):
		log.Error = "TimeOut"
		return nil, errors.New(log.Error)
	case <-ch:
		processor.mutex.Lock()
		r := processor.blocks[id]
		delete(processor.blocks, id)
		//resp, ok := r.Msg.(*Response)
		//if !ok {
		//	err := errors.New("Response type error")
		//	log.Error = err.Error()
		//	return nil, err
		//}

		processor.mutex.Unlock()
		if r.Err != nil {
			log.Error = r.Err.Error()
			return nil, r.Err
		}
		log.RespId = r.Msg.GetHeader().ID
		log.HandTime = time.Since(start)
		log.RespTube = r.Msg.TubeName()
		//log.Error = resp.Error.ErrCode
		return r.Msg, nil
	}
	return nil, errors.New("UnkownErr")
}
