package congestion

import "github.com/bassosimone/psiphon-tunnel-core/psiphon/common/quic/gquic-go/internal/protocol"

type connectionStats struct {
	slowstartPacketsLost protocol.PacketNumber
	slowstartBytesLost   protocol.ByteCount
}
