package email

import (
	"business/support/libraries/mq"
	"time"
)

type Comsumer struct {
	comsumer *mq.Comsumer
}

func NewComsumer(conn *mq.Conn) *Comsumer {
	comsumer := conn.NewComsumer(TubeName)
	return &Comsumer{comsumer}
}

func (c *Comsumer) Reserve(timeout time.Duration) (*EmailMessage, error) {
	var msg EmailMessage
	if err := c.comsumer.Reserve(timeout, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (c *Comsumer) Delete(msg *EmailMessage) error {
	return c.comsumer.Delete(msg.header.ID)
}

func (c *Comsumer) Bury(msg *EmailMessage) error {
	return c.comsumer.Bury(msg.header.ID, msg.header.Pri)
}
