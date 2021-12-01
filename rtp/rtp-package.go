package rtp

//used  rfc3550-RTP

//The RTP header has the following format:
//              0               1               2               3
//1 2 3 4 5 6 7 8 1 2 3 4 5 6 7 8 1 2 3 4 5 6 7 8 1 2 3 4 5 6 7 8
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|V=2|P|X|CC =4|M| PT =7       | sequence number =16           |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|                     timestamp = 32                          |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|           synchronization source (SSRC) identifier          |
//+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
//|           contributing source (CSRC) identifiers            |
//|                               ....                          |
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

import (
	"encoding/binary"
	"fmt"
	"github.com/tmsi-io/balasa/nalu"
)

type RTPPacket struct {
	Version          int
	Padding          int
	Extension        int
	CSRCCount        int
	Mark             int
	PayloadType      int
	SequenceNumber   uint16
	Timestamp        uint32
	SSRC             int
	CSRCList         []uint32
	PayloadOffset    int
	PRTPData         []byte
	ExtensionProfile uint16
	ExtensionPayload []byte

	NaluType      int             //
	IsFuA         bool            //
	FuAHeaderH264 nalu.FuInfoH264 //
	FuAHeaderH265 nalu.FuInfoH265 //
	NaluHeader    []byte          //
	NaluLen       int             //
	NaluRaw       []byte          //
}

func CreatRTPPacket() *RTPPacket {
	return new(RTPPacket)
}

func (pThis *RTPPacket) String() string {
	return fmt.Sprintf("Version: %d, Padding: %d, Extension: %d, CSRCCount: %d, PayloadType: %d, SequenceNumber: %d, Timestamp: %d, SSRC:%d, PayloadOffset:%d",
		pThis.Version, pThis.Padding, pThis.Extension, pThis.CSRCCount, pThis.PayloadType, pThis.SequenceNumber, pThis.Timestamp, pThis.SSRC, pThis.PayloadOffset)
}

func (pThis *RTPPacket) EncodeRTPHeader(ssrc int) []byte {
	//1 2 3 4 5 6 7 8 1 2 3 4 5 6 7 8   16
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|V=2|P|X|CC =4|M| PT =7       |
	pData := make([]byte, RTPPacketLen)
	binary.BigEndian.PutUint16(pData[0:2], uint16(
		pThis.Version<<14+pThis.Padding<<13+pThis.Extension<<12+pThis.CSRCCount<<8+pThis.Mark<<7+pThis.PayloadType))
	binary.BigEndian.PutUint16(pData[2:4], uint16(pThis.SequenceNumber))
	binary.BigEndian.PutUint32(pData[4:8], uint32(pThis.Timestamp))
	binary.BigEndian.PutUint32(pData[8:12], uint32(ssrc))
	if pThis.CSRCCount > 0 {
		for i := 0; i < pThis.CSRCCount; i++ {
			binary.BigEndian.PutUint32(pData[(12+4*i):(12+i*4+4)], uint32(pThis.CSRCList[i]))
		}
	}
	return pData
}

func (pThis *RTPPacket) EncodeNaluData() {
	if pThis.PayloadType == PayloadH264 {
		pThis.NaluHeader = pThis.FuAHeaderH264.EncodeFuInfo()
	} else if pThis.PayloadType == PayloadH265 {
		pThis.NaluHeader = pThis.FuAHeaderH265.EncodeFuInfo()
	}
}

func EncodeRTPHeader(Padding, Extension, CSRCCount, Mark, PayloadType, ssrc int, SequenceNumber uint16, Timestamp uint32) []byte {
	//1 2 3 4 5 6 7 8 1 2 3 4 5 6 7 8   16
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|V=2|P|X|CC =4|M| PT =7       |
	pData := make([]byte, RTPPacketLen)
	binary.BigEndian.PutUint16(pData[0:2], uint16(
		DefaultRTPVer<<14+Padding<<13+Extension<<12+CSRCCount<<8+Mark<<7+PayloadType))
	binary.BigEndian.PutUint16(pData[2:4], SequenceNumber)
	binary.BigEndian.PutUint32(pData[4:8], Timestamp)
	binary.BigEndian.PutUint32(pData[8:12], uint32(ssrc))
	return pData
}

func GetTotalBufferLen(IsFuA bool, PayloadType, CSRCCount, RawLen int) int {
	if !IsFuA {
		return RTPPacketLen + 4*CSRCCount + RawLen
	} else {
		if PayloadType == PayloadH265 {
			return RTPPacketLen + 4*CSRCCount + RawLen + FuHeaderLenH265
		} else if PayloadType == PayloadH264 {
			return RTPPacketLen + 4*CSRCCount + RawLen + FuHeaderLenH264
		} else if PayloadType == PayloadTS {
			return RTPPacketLen + 4*CSRCCount + RawLen
		} else {
			return 0
		}
	}
}

func (pThis *RTPPacket) GetTotalBufferLen() int {
	if !pThis.IsFuA {
		return RTPPacketLen + 4*pThis.CSRCCount + len(pThis.NaluRaw)
	} else {
		if pThis.PayloadType == PayloadH265 {
			return RTPPacketLen + 4*pThis.CSRCCount + len(pThis.NaluRaw) + FuHeaderLenH265
		} else if pThis.PayloadType == PayloadH264 {
			return RTPPacketLen + 4*pThis.CSRCCount + len(pThis.NaluRaw) + FuHeaderLenH264
		} else if pThis.PayloadType == PayloadTS {
			return RTPPacketLen + 4*pThis.CSRCCount + len(pThis.NaluRaw)
		} else {
			return 0
		}
	}
}
