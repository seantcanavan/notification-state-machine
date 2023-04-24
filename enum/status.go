package enum

type Status string

const (
	Created    = "Created"
	Error      = "Error"
	NA         = "NA"
	Processing = "Processing"
	Queued     = "Queued"
	Sent       = "Sent"
)

var statusToString = map[Status]string{
	"Created":    Created,
	"Error":      Error,
	"NA":         NA,
	"Processing": Processing,
	"Queued":     Queued,
	"Sent":       Sent,
}

func (s Status) String() string {
	return statusToString[s]
}

func (s Status) Valid() bool {
	return statusToString[s] != ""
}
