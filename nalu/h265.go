package nalu

type FuInfoH265 struct {
	F       int // 1
	Type    int // 6
	LayerId int // 6
	TID     int // 3
	// set to 1, means start
	start int
	// set to 1, means nal stop
	end int
	// FuType MUST be equal to the field Type of the fragmented NAL unit.
	fuType int
}

func (info *FuInfoH265) SetStart() {
	info.start = 1
}

func (info *FuInfoH265) SetEnd() {
	info.end = 1
}

func (info *FuInfoH265) SetFuType(_type int) {
	info.fuType = _type
}

func (info *FuInfoH265) EncodeFuInfo() []byte {
	pData := make([]byte, 3)
	pData[0] = byte(info.Type << 1)
	pData[1] = 1
	pData[2] = byte(info.start<<7 + info.end<<6 + info.fuType)
	return pData
}

func EncodeFuInfoH265(Start, End, FuType int) []byte {
	Type := NALU_TYPE_FU_H265
	pData := make([]byte, 3)
	pData[0] = byte(Type << 1)
	pData[1] = 1
	pData[2] = byte(Start<<7 + End<<6 + FuType)
	return pData
}
