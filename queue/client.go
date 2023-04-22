package queue

var Client clientI

type clientI interface {
	CreateSession()
	SendMessage(message *Message) error
	GetMessages() map[string][]*Message
	DeleteMessage(receiptHandle string) error
}

func NewClient(live bool) {
	if live {
		sqsQueue := SQSQueue{}
		Client = &sqsQueue
	} else {
		mockQueue := MockQueue{}
		Client = &mockQueue
	}

	Client.CreateSession()
}

type Message struct {
	ID       string      `json:"value,omitempty"`
	Resource string      `json:"target,omitempty"`
	Type     MessageType `json:"type,omitempty"`
}
