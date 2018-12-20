package protocol

// Packet 消息包
type Packet struct {
	Size uint16
	ID   uint16
	Body interface{}
}
