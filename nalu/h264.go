package nalu

import "encoding/binary"

//+--FuIndicator--+----FuHeader---+
//|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|F|NRI|  Type   |S|E|R|  Type   |
//+---------------+---------------+

type FuInfoH264 struct {
	F     int // 0 //
	NRI   int // 3
	IType int // &1F ?
	// set to 1, means start
	start int
	// set to 1, means nal stop
	end int
	// force to 0
	r      int
	fuType int
}

func (info *FuInfoH264) SetStart() {
	info.start = 1
}

func (info *FuInfoH264) SetEnd() {
	info.end = 1
}

func (info *FuInfoH264) SetFuType(_type int) {
	info.fuType = _type
}

func (info *FuInfoH264) EncodeFuInfo() []byte {
	pData := make([]byte, 2)
	binary.BigEndian.PutUint16(pData[0:2], uint16(
		info.F<<15+info.NRI<<13+info.IType<<8+info.start<<7+info.end<<6+info.r<<5+info.fuType))
	return pData
}

func EncodeFuInfo(Start, End, NRI, FuType int) []byte {
	var F, R = 0, 0
	IType := Type_FU_A
	pData := make([]byte, 2)
	binary.BigEndian.PutUint16(pData[0:2], uint16(
		F<<15+NRI<<13+IType<<8+Start<<7+End<<6+R<<5+FuType))
	return pData
}

func DecodeFuInfo(data []int) FuInfoH264 {
	return FuInfoH264{
		F:      data[0] >> 7 & 0x1,
		NRI:    data[0] >> 5 & 0x3,
		IType:  data[0] & 0x1f,
		start:  data[1] >> 7 & 0x1,
		end:    data[1] >> 6 & 0x1,
		r:      data[1] >> 5 & 0x1,
		fuType: data[1] & 0x1f,
	}
}
