package queue

type MockQueue struct {
	messages map[string][]*Message
	queueURL string
}

func (mockQueue *MockQueue) CreateSession() {
	mockQueue.messages = make(map[string][]*Message)
}

func (mockQueue *MockQueue) SendMessage(message *Message) error {
	if mockQueue.messages[message.ID] == nil {
		mockQueue.messages[message.ID] = make([]*Message, 0)
	}

	mockQueue.messages[message.ID] = append(mockQueue.messages[message.ID], message)
	return nil
}

func (mockQueue *MockQueue) GetMessages() map[string][]*Message {
	return mockQueue.messages
}

func (mockQueue *MockQueue) DeleteMessage(receiptHandle string) error {
	return nil
}
