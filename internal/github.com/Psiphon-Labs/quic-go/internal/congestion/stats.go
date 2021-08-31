package congestion

import "github.com/bassosimone/psiphon-tunnel-core/internal/github.com/Psiphon-Labs/quic-go/internal/protocol"

type connectionStats struct {
	slowstartPacketsLost protocol.PacketNumber
	slowstartBytesLost   protocol.ByteCount
}
