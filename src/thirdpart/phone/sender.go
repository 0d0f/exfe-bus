package phone

type Sender interface {
	Codes() []string
	Send(phone string, contents string) (string, error)
}
