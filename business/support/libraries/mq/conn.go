package mq

import (
	"encoding/json"
	"github.com/kr/beanstalk"
	"strconv"
	"time"
)

type Conn struct {
	*beanstalk.Conn
}

func NewConn(network, addr string) (*Conn, error) {
	c, err := beanstalk.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Conn{c}, nil
}

func (c *Conn) NewProducer(tubeName string) *Producer {
	return &Producer{&beanstalk.Tube{c.Conn, tubeName}, c}
}

func (c *Conn) NewComsumer(name string) *Comsumer {
	ts := beanstalk.NewTubeSet(c.Conn, name)
	return &Comsumer{ts, c}
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}

func (c *Conn) Pause(name string, d time.Duration) error {
	t := &beanstalk.Tube{Conn: c.Conn, Name: name}
	return t.Pause(d)
}

func (c *Conn) PeekBuried(name string, msg Message) error {
	tube := &beanstalk.Tube{Conn: c.Conn, Name: name}
	id, body, err := tube.PeekBuried()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}
	header, err := c.getHeader(id)
	if err != nil {
		return err
	}
	msg.SetHeader(header)
	return nil
}

func (c *Conn) PeekReady(name string, msg Message) error {
	tube := &beanstalk.Tube{Conn: c.Conn, Name: name}
	id, body, err := tube.PeekReady()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}

	header, err := c.getHeader(id)
	if err != nil {
		return err
	}
	msg.SetHeader(header)
	return nil
}

func (c *Conn) PeekDelayed(name string, msg Message) error {
	tube := &beanstalk.Tube{Conn: c.Conn, Name: name}
	id, body, err := tube.PeekDelayed()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}

	header, err := c.getHeader(id)
	if err != nil {
		return err
	}
	msg.SetHeader(header)
	return nil
}

func (c *Conn) getHeader(id uint64) (*Header, error) {
	var info map[string]string
	var header Header
	header.ID = id
	info, err := c.Conn.StatsJob(header.ID)
	if err != nil {
		return nil, err
	}
	if tube, ok := info["tube"]; ok {
		header.Tube = tube
	}
	if state, ok := info["state"]; ok {
		header.State = state
	}
	if v, ok := info["pri"]; ok {
		if pri, err := strconv.Atoi(v); err == nil {
			header.Pri = uint32(pri)
		}
	}
	if v, ok := info["age"]; ok {
		if age, err := strconv.Atoi(v); err == nil {
			header.Age = time.Second * time.Duration(age)
		}
	}
	if v, ok := info["delay"]; ok {
		if delay, err := strconv.Atoi(v); err == nil {
			header.Delay = time.Second * time.Duration(delay)
		}
	}
	if v, ok := info["ttr"]; ok {
		if ttr, err := strconv.Atoi(v); err == nil {
			header.TTR = time.Second * time.Duration(ttr)
		}
	}
	if v, ok := info["time-left"]; ok {
		if timeLeft, err := strconv.Atoi(v); err == nil {
			header.TimeLeft = time.Second * time.Duration(timeLeft)
		}
	}
	if v, ok := info["file"]; ok {
		if file, err := strconv.Atoi(v); err == nil {
			header.File = uint32(file)
		}
	}
	if v, ok := info["reserves"]; ok {
		if reserves, err := strconv.Atoi(v); err == nil {
			header.Rservers = uint32(reserves)
		}
	}
	if v, ok := info["timeouts"]; ok {
		if timeouts, err := strconv.Atoi(v); err == nil {
			header.Timeouts = uint32(timeouts)
		}
	}
	if v, ok := info["releases"]; ok {
		if releases, err := strconv.Atoi(v); err == nil {
			header.Releases = uint32(releases)
		}
	}
	if v, ok := info["buries"]; ok {
		if buries, err := strconv.Atoi(v); err == nil {
			header.Buries = uint32(buries)
		}
	}
	if v, ok := info["kicks"]; ok {
		if kicks, err := strconv.Atoi(v); err == nil {
			header.Kicks = uint32(kicks)
		}
	}
	return &header, nil
}
