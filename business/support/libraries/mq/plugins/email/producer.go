package email

import "business/support/libraries/mq"

type Producer struct {
	producer *mq.Producer
}

func NewProducer(conn *mq.Conn) *Producer {
	producer := conn.NewProducer(TubeName)
	return &Producer{producer}
}

func (p *Producer) Put(email *EmailMessage) error {
	if email.EmailType == 0 {
		if email.TemplateName != "" {
			email.EmailType = EmailTypeTemplate
		} else {
			email.EmailType = EmailTypeHTML
		}
	}
	return p.producer.Put(email)
}
