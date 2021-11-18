package hcwire

type StateUpdate struct {
	BlockDay         uint32
	LocalUpdates     uint32
	RemoteUpdates    uint32
	LocalSigOfRemote [64]byte
}
