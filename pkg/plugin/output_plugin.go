package plugin

type OutputPlugin interface {
	Init(map[string]interface{}) error
	Send([]byte) error
}
