package hcwire

type StateOverride struct {
	Blockday         uint32
	LocalBalanceMSAT uint64
	LocalUpdates     uint32
	RemoteUpdates    uint32
	LocalSigOfRemote [64]byte
}
