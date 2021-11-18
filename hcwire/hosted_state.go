package hcwire

import (
	"github.com/lightningnetwork/lnd/lnwire"
)

type HostedState struct {
	ChannelID            lnwire.ChannelID
	NumNextLocalUpdates  uint16
	NextLocalUpdates     []byte // have to define new type here; some other fields still missing here as well
	LastCrossSignedState LastCrossSignedState
}
