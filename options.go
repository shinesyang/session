package session

type Option struct {
	MaxAge    int
	Key       []byte
	SetHeader string
	Header    string
}
