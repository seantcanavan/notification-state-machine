package queue

type MessageType string

const (
	CreateReq = "CreateReq"
)

var messageTypeToString = map[MessageType]string{
	CreateReq: "CreateReq",
}

func (r MessageType) String() string {
	return messageTypeToString[r]
}

func (r MessageType) Valid() bool {
	return messageTypeToString[r] != ""
}
