package ts

import (
	"io"
)

type Es2Ts struct {
	seq      int
	TsPacket []byte // 188 byte
	CurrLen  int
	First    bool
}

type PATPrograms []PATProgram

type PATProgram struct {
	ProgramNumber int
	Pid           int
}

type PMTStreams []PMTStream

type PMTStream struct {
	streamType    int
	elementaryPID int
	eSInfoLength  int
	esData        []byte
}

func (t *Es2Ts) TransEs2TsData(w io.Writer, isH265 bool, isIFrame bool, nalu []byte, pts uint, withPTS bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if t.First {
		// PAT
		t.WriteTsPAT(1, PATPrograms{PATProgram{
			ProgramNumber: 1,      // PMT ProgramNumber
			Pid:           PMTPid, // PMT pid
		}})
		_, _ = w.Write(t.TsPacket)
		// PMT
		if isH265 {
			t.WriteTsPMT(1, VideoPid, PMTStreams{PMTStream{
				streamType:    H265StreamType,
				elementaryPID: VideoPid,
			}})
		} else {
			t.WriteTsPMT(1, VideoPid, PMTStreams{PMTStream{
				streamType:    H264StreamType,
				elementaryPID: VideoPid,
			}})
		}
		_, _ = w.Write(t.TsPacket)
		t.First = false
	}
	pesHead := EncoderPesHead(PesVideo, len(nalu), pts, withPTS)
	t.WriteEsStream(w, pesHead, true, nalu, len(nalu), isIFrame, pts, 0)
}

func (t *Es2Ts) nextContinuityCount() (result int) {
	result = (t.seq + 1) & 0x0f
	t.seq = result
	return
}

func (t *Es2Ts) ReTsPacket() {
	t.CurrLen = 0
}
