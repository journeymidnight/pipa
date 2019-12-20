package pipa

type Operation interface {
	GetType() string
	GetOption(key string) string
	DoProcess(data []byte) (result []byte, err error)
	Close() error
}
