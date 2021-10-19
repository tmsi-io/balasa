package ts

func EncoderPesHead(streamId int, dataLen int, pts uint, withPTS bool) []byte {
	var header []byte
	var payloadLen int
	if !withPTS {
		payloadLen = dataLen + 3
		header = make([]byte, 9)
	} else {
		payloadLen = dataLen + 8
		header = make([]byte, 14)
	}
	header[0] = 0x00
	header[1] = 0x00
	header[2] = 0x01
	header[3] = byte(streamId)
	header[6] = 0x80
	if !withPTS {
		header[7] = 0x00
		header[8] = 0x00
	} else {
		header[7] = 0x80
		header[8] = 0x05
		encodePTS(header[9:14], pts)
	}
	if payloadLen > PesMaxPayload {
		header[4] = 0
		header[5] = 0
	} else {
		header[4] = byte((payloadLen & 0xff00) >> 8)
		header[5] = byte(payloadLen & 0x00ff)
	}

	return header
}
