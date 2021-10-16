package ts

const (
	PmtProgramNumber = 0x01
)

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
	PACKET_SIZE           = 188
	MAX_PAYLOAD_SIZE      = 184
	PES_MAXPAYLOAD        = 0xffff
	CRC32_POLY            = 0x04c11db7
	MAX_PTS_VALUE    uint = 0x1FFFFFFFF
)
