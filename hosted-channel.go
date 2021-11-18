package main

import "github.com/raphjaph/go-hosted-channels/hcwire"

var defaultHCInit hcwire.InitHostedChannel = hcwire.InitHostedChannel{
	MaxHTLCValueInFlightMSAT:           100000000,
	HTLCMinimumMSAT:                    1000,
	MaxAcceptedHTLCs:                   30,
	ChannelCapacityMSAT:                1000000000,
	LiabilityDeadlineBlockdays:         360,
	MinimalOnChainRefundAmountSatoshis: 100000,
	InitialClientBalanceMSAT:           0,
	Features:                           []byte{},
}
