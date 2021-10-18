package ts

func (t *Es2Ts) WriteTsHead(pid int, dataLen int) {
	t.TsPacket[t.CurrLen] = 0x47
	t.TsPacket[t.CurrLen+1] = (byte)(0x40 | ((pid & 0x1f00) >> 8))
	t.TsPacket[t.CurrLen+2] = (byte)(pid & 0xff)
	controls := 0x10
	t.TsPacket[t.CurrLen+3] = (byte)(controls | t.nextContinuityCount())
	pointer := PACKET_SIZE - 5 - dataLen
	t.TsPacket[t.CurrLen+4] = byte(pointer)
	t.CurrLen += 5
	for i := 0; i < pointer; i++ {
		t.TsPacket[t.CurrLen] = 0xff
		t.CurrLen++
	}
}
