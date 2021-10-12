package rtp

import (
	"encoding/binary"
	"errors"
)

const RTSPRTPTCPHeaderLen = 4
const RTSPRTPTCPHeaderTag = 0x24

type RTSPRTPTCPHeader struct {
	Magic        int //8 b
	Channel      int //8 b
	RtpPacketLen int //16 b
}

// ParserRTSPRTPTCPHeader Parse RTP-TCP header for TCP transport mode
func ParserRTSPRTPTCPHeader(pBuffer []byte) (*RTSPRTPTCPHeader, error) {
	if len(pBuffer) < RTSPRTPTCPHeaderLen {
		return nil, errors.New("ParserRTSPRTPTCPHeader Error ! Reason Buffer too small ! ")
	}
	pHeader := new(RTSPRTPTCPHeader)
	pHeader.Magic = int(pBuffer[0])
	pHeader.Channel = int(pBuffer[1])
	pHeader.RtpPacketLen = int((pBuffer[2] << 8) | pBuffer[3])
	return pHeader, nil
}

// BuildRTSPRTPTCPHeader add RTP-TCP header for TCP transport mode.
func BuildRTSPRTPTCPHeader(nChannel, nDataLen int, pBuffer []byte) error {
	if len(pBuffer) < RTSPRTPTCPHeaderLen {
		return errors.New("ParserRTSPRTPTCPHeader Error ! Reason Buffer too small ! ")
	}
	pBuffer[0] = RTSPRTPTCPHeaderTag
	pBuffer[1] = byte(nChannel)
	binary.BigEndian.PutUint16(pBuffer[2:4], uint16(nDataLen))
	return nil
}
