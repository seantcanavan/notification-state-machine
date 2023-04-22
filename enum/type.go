package enum

type Type string

const Email = "Email"
const SMS = "SMS"
const Snail = "Snail"

func (t Type) String() string {
	return typeToString[t]
}

func (t Type) Valid() bool {
	return typeToString[t] != ""
}

var typeToString = map[Type]string{
	"Email": Email,
	"SMS":   SMS,
	"Snail": Snail,
}
