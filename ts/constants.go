package ts

const (
	PATPid   = 0x00
	PMTPid   = 0x66
	VideoPid = 0x68
)

const (
	H264StreamType = 0x1b
	H265StreamType = 0x24
)

const (
	PesVideo = 0xe0
)

const (
	PacketSize          = 188
	MaxPayloadSize      = 184
	PesMaxPayload       = 0xffff
	Crc32Poly           = 0x04c11db7
	MaxPtsValue    uint = 0x1FFFFFFFF
)
