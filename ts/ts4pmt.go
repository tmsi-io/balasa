package ts

func (t *Es2Ts) WriteTsPMT(programNumber int, pcrId int, list PMTStreams) {
	t.ReTsPacket()
	dataLen := 16
	sectionLength := 13
	offset := 0

	for ii := 0; ii < len(list); ii++ {
		sectionLength += 5 + list[ii].eSInfoLength
		dataLen += 5 + list[ii].eSInfoLength
	}
	pointer := PacketSize - 5 - dataLen
	t.WriteTsHead(PMTPid, dataLen)
	t.TsPacket[t.CurrLen+0] = 0x02
	t.TsPacket[t.CurrLen+1] = byte(0xb0 | ((sectionLength & 0x0f00) >> 8))
	t.TsPacket[t.CurrLen+2] = byte(sectionLength & 0x0ff)
	t.TsPacket[t.CurrLen+3] = byte((programNumber & 0xff00) >> 8)
	t.TsPacket[t.CurrLen+4] = byte(programNumber & 0x00ff)
	t.TsPacket[t.CurrLen+5] = 0xc1
	t.TsPacket[t.CurrLen+6] = 0x00
	t.TsPacket[t.CurrLen+7] = 0x00
	t.TsPacket[t.CurrLen+8] = byte(0xe0 | ((pcrId & 0x1f00) >> 8))
	t.TsPacket[t.CurrLen+9] = byte(pcrId & 0x00ff)
	t.TsPacket[t.CurrLen+10] = 0xf0
	t.TsPacket[t.CurrLen+11] = 0x00
	t.CurrLen += 12
	offset += 12
	for _, item := range list {
		pid := item.elementaryPID
		esLen := item.eSInfoLength
		t.TsPacket[t.CurrLen+0] = byte(item.streamType)
		t.TsPacket[t.CurrLen+1] = (byte)(0xE0 | ((pid & 0x1F00) >> 8))
		t.TsPacket[t.CurrLen+2] = (byte)(pid & 0x00FF)
		t.TsPacket[t.CurrLen+3] = byte(((esLen & 0xff00) >> 8) | 0xf0)
		t.TsPacket[t.CurrLen+4] = byte(esLen & 0x00ff)
		copy(t.TsPacket[t.CurrLen:], item.esData[:item.eSInfoLength])
		t.CurrLen += 5 + esLen
		offset += 5 + esLen
	}
	crc32 := crc32Block(0xffffffff, t.TsPacket[5+pointer:], offset)
	t.TsPacket[t.CurrLen+0] = byte((crc32 & 0xff000000) >> 24)
	t.TsPacket[t.CurrLen+1] = byte((crc32 & 0x00ff0000) >> 16)
	t.TsPacket[t.CurrLen+2] = byte((crc32 & 0x0000ff00) >> 8)
	t.TsPacket[t.CurrLen+3] = byte(crc32 & 0x000000ff)
	t.CurrLen = t.CurrLen + 4
}
