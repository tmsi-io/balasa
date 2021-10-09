package rtp

import (
	"encoding/binary"
	"errors"
	"fmt"
)

//	/*
//	 *  0                   1                   2                   3
//	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	 * |V=2|P|X|  CC   |M|     PT      |       sequence number         |
//	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	 * |                           timestamp                           |
//	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	 * |           synchronization source (SSRC) identifier            |
//	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
//	 * |            contributing source (CSRC) identifiers             |
//	 * |                             ....                              |
//	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

//	https://tools.ietf.org/html/rfc3550#section-6.4.1
//	 */

func DecodeRTPPacketFromByte(pBuffer []byte) (*RTPPacket, error) {
	PackLen := len(pBuffer)
	if PackLen < RTPPacketLen {
		return nil, errors.New("RTP Packet Parser Failed ! RTP Buffer too Small ! ")
	}
	var pRTPPacket RTPPacket
	pRTPPacket.Version = int(pBuffer[0] >> VersionShift & VersionMask)
	pRTPPacket.Padding = int(pBuffer[0] >> PaddingShift & PaddingMask)
	pRTPPacket.Extension = int(pBuffer[0] >> ExtensionShift & ExtensionMask)
	pRTPPacket.CSRCCount = int(pBuffer[0] & CcMask)
	pRTPPacket.Mark = int(pBuffer[1] >> MarkerShift & MarkerMask)
	pRTPPacket.PayloadType = int(pBuffer[1] & PtMask)
	pRTPPacket.SequenceNumber = binary.BigEndian.Uint16(pBuffer[SeqNumOffset : SeqNumOffset+SeqNumLength])
	pRTPPacket.Timestamp = binary.BigEndian.Uint32(pBuffer[TimestampOffset : TimestampOffset+TimestampLength])
	pRTPPacket.SSRC = int(binary.BigEndian.Uint32(pBuffer[SSRCOffset : SSRCOffset+SSRCLength]))
	currOffset := CSRCOffset
	if PackLen-currOffset > pRTPPacket.CSRCCount*CSRCLength {
		currOffset += pRTPPacket.CSRCCount * CSRCLength
	}
	if PackLen < currOffset {
		return nil, fmt.Errorf("RTP header size insufficient; %d < %d", PackLen, currOffset)
	}
	pRTPPacket.CSRCList = make([]uint32, pRTPPacket.CSRCCount)
	for i := range pRTPPacket.CSRCList {
		offset := CSRCOffset + (i * CSRCLength)
		pRTPPacket.CSRCList[i] = binary.BigEndian.Uint32(pBuffer[offset:])
	}
	if pRTPPacket.Extension > 0 { // is extension enable?
		if PackLen < currOffset+4 {
			return nil, fmt.Errorf("RTP header size too small for extension; %d < %d", PackLen, currOffset)
		}
		pRTPPacket.ExtensionProfile = binary.BigEndian.Uint16(pBuffer[currOffset:])
		currOffset += 2
		extensionLength := int(binary.BigEndian.Uint16(pBuffer[currOffset:])) * 4
		currOffset += 2
		if PackLen < currOffset+extensionLength {
			return nil, fmt.Errorf("RTP header size too small for extension length; %d < %d", PackLen, currOffset+extensionLength)
		}
		pRTPPacket.ExtensionPayload = pBuffer[currOffset : currOffset+extensionLength]
		currOffset += len(pRTPPacket.ExtensionPayload)
	}
	if pRTPPacket.Padding != 0 {
		padLen := int(pBuffer[PackLen-1] & 0xff)
		if padLen < 0 {
			pRTPPacket.PRTPData = pBuffer[RTPPacketLen+4*pRTPPacket.CSRCCount : PackLen-4*pRTPPacket.Padding]
		} else {
			pRTPPacket.PRTPData = pBuffer[RTPPacketLen+4*pRTPPacket.CSRCCount : PackLen-padLen]
		}
	} else {
		pRTPPacket.PRTPData = pBuffer[RTPPacketLen+4*pRTPPacket.CSRCCount : PackLen-4*pRTPPacket.Padding]
	}
	pRTPPacket.PayloadOffset = currOffset
	return &pRTPPacket, nil
}

// DecodeRTPInfoFromByte decode rtp info from byte stream
func DecodeRTPInfoFromByte(pBuffer []byte) (SSRC int, Timestamp uint32, SequenceNumber uint16, mark int, PRTPData []byte, err error) {
	PackLen := len(pBuffer)
	if PackLen < RTPPacketLen {
		return 0, 0, 0, 0, nil, errors.New("RTP Packet Parser Failed ! RTP Buffer too Small ! ")
	}
	Padding := int(pBuffer[0] >> PaddingShift & PaddingMask)
	CSRCCount := int(pBuffer[0] & CcMask)
	mark = int(pBuffer[1] >> 7 & 0x1)
	SequenceNumber = binary.BigEndian.Uint16(pBuffer[SeqNumOffset : SeqNumOffset+SeqNumLength])
	Timestamp = binary.BigEndian.Uint32(pBuffer[TimestampOffset : TimestampOffset+TimestampLength])
	SSRC = int(binary.BigEndian.Uint32(pBuffer[SSRCOffset : SSRCOffset+SSRCLength]))
	if Padding != 0 {
		padLen := int(pBuffer[PackLen-1] & 0xff)
		if padLen < 0 {
			PRTPData = pBuffer[RTPPacketLen+4*CSRCCount : PackLen-4*Padding]
		} else {
			PRTPData = pBuffer[RTPPacketLen+4*CSRCCount : PackLen-padLen]
		}
	} else {
		PRTPData = pBuffer[RTPPacketLen+4*CSRCCount : PackLen-4*Padding]
	}
	return SSRC, Timestamp, SequenceNumber, mark, PRTPData, nil
}
