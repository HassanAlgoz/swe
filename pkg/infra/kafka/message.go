package kafka

type Message struct {
	Type   string `json:"type"`
	Header []byte `json:"header"`
	Body   []byte `json:"body"`
}
