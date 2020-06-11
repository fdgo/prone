package mq

import (
	"encoding/json"
	"time"

	"github.com/kr/beanstalk"
)

type Comsumer struct {
	ts   *beanstalk.TubeSet
	conn *Conn
}

func (c *Comsumer) Reserve(timeout time.Duration, msg Message) error {
	id, body, err := c.ts.Reserve(timeout)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}

	header, err := c.conn.getHeader(id)
	if err != nil {
		return err
	}
	msg.SetHeader(header)
	return nil
}

func (c *Comsumer) Delete(id uint64) error {
	return c.conn.Delete(id)
}

func (c *Comsumer) Release(id uint64, pri uint32, delay time.Duration) error {
	return c.conn.Release(id, pri, delay)
}

func (c *Comsumer) Touch(id uint64) error {
	return c.conn.Touch(id)
}

func (c *Comsumer) Bury(id uint64, pri uint32) error {
	return c.conn.Bury(id, pri)
}
