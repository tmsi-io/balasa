package ts

import (
	"io"
)

func (t *Es2Ts) WriteEsStream(w io.Writer, pesHead []byte, start bool, es []byte, dataLen int, gotPcr bool, PCRBase, PCRExtn uint) {
	pesHeadLen := len(pesHead)
	t.ReTsPacket()
	var gotAdaptationField bool
	t.TsPacket[t.CurrLen] = 0x47
	if !start {
		t.TsPacket[t.CurrLen+1] = (byte)(0x00 | ((VideoPid & 0x1f00) >> 8))
	} else {
		t.TsPacket[t.CurrLen+1] = (byte)(0x40 | ((VideoPid & 0x1f00) >> 8))
	}
	t.TsPacket[2] = (byte)(VideoPid & 0xff)
	if start && gotPcr {
		controls := 0x30
		t.TsPacket[3] = byte(controls | (t.nextContinuityCount()))
		t.TsPacket[t.CurrLen+4] = 7
		t.TsPacket[t.CurrLen+5] = 0x10
		t.TsPacket[t.CurrLen+6] = byte(PCRBase >> 25)
		t.TsPacket[t.CurrLen+7] = byte((PCRBase >> 17) & 0xFF)
		t.TsPacket[t.CurrLen+8] = byte((PCRBase >> 9) & 0xFF)
		t.TsPacket[t.CurrLen+9] = byte((PCRBase >> 1) & 0xFF)
		t.TsPacket[t.CurrLen+10] = byte(((PCRBase & 0x1) << 7) | 0x7E | (PCRExtn >> 8))
		t.TsPacket[t.CurrLen+11] = byte(PCRExtn >> 1)
		t.CurrLen += 12
		gotAdaptationField = true
	} else if dataLen+pesHeadLen < MaxPayloadSize {
		controls := 0x30
		t.TsPacket[t.CurrLen+3] = (byte)(controls | t.nextContinuityCount())
		if dataLen+pesHeadLen == MaxPayloadSize-1 {
			t.TsPacket[t.CurrLen+4] = 0
			t.CurrLen += 5
		} else {
			t.TsPacket[t.CurrLen+4] = 1
			t.TsPacket[t.CurrLen+5] = 0x00
			t.CurrLen += 6
		}
		gotAdaptationField = true
	} else {
		controls := 0x10
		t.TsPacket[t.CurrLen+3] = (byte)(controls | t.nextContinuityCount())
		t.CurrLen += 4
	}
	if gotAdaptationField {
		if dataLen+pesHeadLen < PacketSize-t.CurrLen {
			padLen := PacketSize - t.CurrLen - dataLen - pesHeadLen
			for i := 0; i < padLen; i++ {
				t.TsPacket[t.CurrLen] = 0xff
				t.CurrLen++
			}
			t.TsPacket[4] += byte(padLen)
		}
	}
	if start {
		n := copy(t.TsPacket[t.CurrLen:], pesHead)
		t.CurrLen += n
	}
	if PacketSize-t.CurrLen == dataLen {
		copy(t.TsPacket[t.CurrLen:], es)
		_, _ = w.Write(t.TsPacket)
	} else {
		n := copy(t.TsPacket[t.CurrLen:], es)
		_, _ = w.Write(t.TsPacket)
		t.WriteEsStream(w, nil, false, es[n:], dataLen-n, false, 0, 0)
	}
}
