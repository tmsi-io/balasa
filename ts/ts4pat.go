package ts

func (t *Es2Ts) WriteTsPAT(transportStreamId int, list PATPrograms) {
	t.ReTsPacket()
	dataLen := 12 + len(list)*4
	sectionLength := 9 + len(list)*4
	pointer := PacketSize - 5 - dataLen
	offset := 0
	t.WriteTsHead(PATPid, dataLen)
	t.TsPacket[t.CurrLen] = 0x00
	t.TsPacket[t.CurrLen+1] = (byte)(0xb0 | ((sectionLength & 0x0f00) >> 8))
	t.TsPacket[t.CurrLen+2] = (byte)(sectionLength & 0x0ff)
	t.TsPacket[t.CurrLen+3] = (byte)((transportStreamId & 0xff00) >> 8)
	t.TsPacket[t.CurrLen+4] = (byte)(transportStreamId & 0x00ff)
	t.TsPacket[t.CurrLen+5] = 0xc1
	t.TsPacket[t.CurrLen+6] = 0x00
	t.TsPacket[t.CurrLen+7] = 0x00
	t.CurrLen += 8
	offset += 8
	for _, item := range list {
		t.TsPacket[t.CurrLen+0] = byte((item.ProgramNumber & 0xff00) >> 8)
		t.TsPacket[t.CurrLen+1] = byte(item.ProgramNumber & 0x00ff)
		t.TsPacket[t.CurrLen+2] = byte(0xe0 | (item.Pid&0x1f00)>>8)
		t.TsPacket[t.CurrLen+3] = (byte)(item.Pid & 0x00ff)
		t.CurrLen += 4
		offset += 4
	}
	crc32 := crc32Block(0xffffffff, t.TsPacket[5+pointer:], offset)
	t.TsPacket[t.CurrLen] = byte((crc32 & 0xff000000) >> 24)
	t.TsPacket[t.CurrLen+1] = byte((crc32 & 0x00ff0000) >> 16)
	t.TsPacket[t.CurrLen+2] = byte((crc32 & 0x0000ff00) >> 8)
	t.TsPacket[t.CurrLen+3] = byte(crc32 & 0x000000ff)
	t.CurrLen = t.CurrLen + 4
}
