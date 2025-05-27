package sender

type SenderInstance interface {
	Push(data []byte)
	Start()
	Stop()
	GetName() string
}
