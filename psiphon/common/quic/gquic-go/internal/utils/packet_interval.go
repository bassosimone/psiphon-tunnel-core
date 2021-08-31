package utils

import "github.com/bassosimone/psiphon-tunnel-core/psiphon/common/quic/gquic-go/internal/protocol"

// PacketInterval is an interval from one PacketNumber to the other
type PacketInterval struct {
	Start protocol.PacketNumber
	End   protocol.PacketNumber
}
