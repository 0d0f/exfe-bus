package imsg

type LoadType int

const (
	Ping LoadType = iota
	Pong
	Send
	Respond
)

type Load struct {
	Type    LoadType `json:"type"`
	To      string   `json:"to"`
	Content string   `json:"content"`
}
