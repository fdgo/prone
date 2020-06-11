package mq

import (
	"github.com/kr/beanstalk"
	"time"
)

type Producer struct {
	tube *beanstalk.Tube
	Conn *Conn
}

func (p *Producer) Put(msg Message) (uint64, error) {
	body, err := msg.Marshal()
	if err != nil {
		return 0, err
	}
	header := msg.GetHeader()
	id, err := p.tube.Put(body, header.Pri, header.Delay, header.TTR)
	if err != nil {
		return 0, err
	}
	header.ID = id
	return id, nil
}

func (p *Producer) Pause(d time.Duration) error {
	return p.tube.Pause(d)
}

func (p *Producer) Kick(bound int) (n int, err error) {
	return p.tube.Kick(bound)
}

func (p *Producer) Stats() (map[string]string, error) {
	return p.tube.Stats()
}
