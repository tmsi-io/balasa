package rtp

import (
	"errors"
	"fmt"
	"github.com/tmsi-io/balasa/nalu"
)

// DecodeH26XFrameToRTPs2 Decode original frame data to rtp package.
func DecodeH26XFrameToRTPs2(sType int, FrameBuff []byte, ts uint32) (FinallySend []RTPPacket, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("%v", err1)
		}
	}()
	DataLen := len(FrameBuff)
	if DataLen == 0 {
		return nil, errors.New("FrameBuff Too short. ")
	}
	if headerLen, errSearch := ParserH26XHeaderLen(FrameBuff, DataLen); errSearch != nil {
		return nil, errSearch
	} else {
		if len(FrameBuff) < headerLen+1 {
			return nil, errors.New("FrameBuff Too short. ")
		}
		NeedSendRaw := FrameBuff[headerLen:]
		NeedSendLen := DataLen - headerLen
		if NeedSendLen < 1 {
			return nil, errors.New("Frame buff too short to decode. ")
		} else {
			var isStart bool = true
			var fuIndex uint16
			var Nalu_H264, Nalu_H265, LayerIdH265, TIdH265 int
			if NeedSendLen > Default_MTU {
				if sType == PayloadH265 {
					Nalu_H265 = int((NeedSendRaw[0] >> 1) & 0x3F)
					if 0 <= Nalu_H265 && Nalu_H265 <= 49 {
						LayerIdH265 = int(((NeedSendRaw[0] & 0x01) << 5) | ((NeedSendRaw[1] >> 3) & 0x1f))
						TIdH265 = int(NeedSendRaw[1] & 0x07)
						NeedSendRaw = NeedSendRaw[2:]
						NeedSendLen -= 2
					} else {
						return FinallySend, fmt.Errorf("NaluH265 Error: %d", Nalu_H265)
					}

				} else if sType == PayloadH264 {
					Nalu_H264 = int(NeedSendRaw[0] & 0x1F)
					NeedSendRaw = NeedSendRaw[1:]
					NeedSendLen -= 1
				}
			}
			for NeedSendLen > 0 {
				fuIndex++
				var Pack RTPPacket
				Pack.Version = DefaultRTPVer
				Pack.Timestamp = ts
				Pack.PayloadType = sType
				Pack.FuAHeaderH264 = nalu.FuInfoH264{}
				Pack.FuAHeaderH265 = nalu.FuInfoH265{}
				if sType == PayloadH264 {
					Pack.FuAHeaderH264.SetFuType(Nalu_H264)
					Pack.FuAHeaderH264.NRI = nalu.DefaultNRI
					Pack.FuAHeaderH264.IType = nalu.Type_FU_A
				} else if sType == PayloadH265 {
					Pack.FuAHeaderH265.SetFuType(Nalu_H265)
					Pack.FuAHeaderH265.LayerId = LayerIdH265
					Pack.FuAHeaderH265.TID = TIdH265
					Pack.FuAHeaderH265.Type = nalu.NALU_TYPE_FU_H265
				}
				if NeedSendLen <= Default_MTU && isStart {
					Pack.NaluRaw = make([]byte, NeedSendLen)
					copy(Pack.NaluRaw, NeedSendRaw)
					Pack.NaluLen = NeedSendLen
					Pack.Mark = 1
					FinallySend = append(FinallySend, Pack)
					break
				} else {
					Pack.IsFuA = true
					if isStart {
						Pack.FuAHeaderH264.SetStart()
						Pack.FuAHeaderH265.SetStart()
					} else {
						Pack.FuAHeaderH264.SetStart()
						Pack.FuAHeaderH265.SetStart()
					}
					if NeedSendLen <= Default_MTU {
						Pack.FuAHeaderH264.SetEnd()
						Pack.FuAHeaderH265.SetEnd()
						Pack.NaluRaw = make([]byte, NeedSendLen)
						copy(Pack.NaluRaw, NeedSendRaw)
						Pack.NaluLen = NeedSendLen
						Pack.Mark = 1
						FinallySend = append(FinallySend, Pack)
						break
					} else {
						Pack.FuAHeaderH264.SetEnd()
						Pack.FuAHeaderH265.SetEnd()
						Pack.NaluLen = Default_MTU
						Pack.NaluRaw = make([]byte, Default_MTU)
						copy(Pack.NaluRaw, NeedSendRaw[0:Default_MTU])
						NeedSendRaw = NeedSendRaw[Default_MTU:]
						NeedSendLen -= Default_MTU
						FinallySend = append(FinallySend, Pack)
						isStart = false
					}
				}
			}
		}
		return FinallySend, nil
	}
}

func DecodeAudioFrameToRTPs(sType int, FrameBuff []byte, ts uint32) (RTPPacket, error) {
	var Pack RTPPacket
	DataLen := len(FrameBuff)
	if DataLen == 0 {
		return Pack, errors.New("FrameBuff Too short. ")
	}
	Pack.Version = DefaultRTPVer
	Pack.Timestamp = ts
	Pack.PayloadType = sType
	Pack.NaluLen = len(FrameBuff)
	Pack.NaluRaw = FrameBuff
	return Pack, nil
}

func ParserH26XHeaderLen(FrameBuff []byte, dLen int) (index int, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Can't Find HeaderLen in H26x: %v. ", err2)
		}
	}()
	Pos := 0
	IsFindPos := false
	for Pos < dLen-0 {
		if FrameBuff[Pos] == 0 && FrameBuff[Pos+1] == 1 {
			IsFindPos = true
			break
		}
		Pos++
	}
	if IsFindPos {
		return Pos + 2, nil
	} else {
		return 0, fmt.Errorf("Can't Find HeaderLen in H264Raw ")
	}
}
