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
	if payloadLen > PES_MAXPAYLOAD {
		header[4] = 0
		header[5] = 0
	} else {
		header[4] = byte((payloadLen & 0xff00) >> 8)
		header[5] = byte(payloadLen & 0x00ff)
	}

	return header
}

/*
		 PTS/DTS = （PTS1 & 0x0e) << 29 + (PTS2 & 0xfffe) << 14 + (PTS3 & 0xfffe ) >> 1
		+-------+-------+---------------+---------------+---------------+-------------+-+
	  	|7|6|5|4|3|2|1|0|7|6|5|4|3|2|1|0|7|6|5|4|3|2|1|0|7|6|5|4|3|2|1|0|7|6|5|4|3|2|1|0|
		+---+---+-----+-+---------------+-------------+-+---------------+-------------+-+
		|0 0|1 0| PTS1|1|			PTS2 29..15		  |1|			PTS3 14..00		  |1|
		+---+---+-----+-+---------------+-------------+-+---------------+-------------+-+
		|0 0|1 1| DTS1|1|			DTS2 29..15		  |1|			DTS3 14..00		  |1|
		+---+---+-----+-+---------------+-------------+-+---------------+-------------+-+
*/

func encodePTS(data []byte, pts uint) {
	if pts > MAX_PTS_VALUE {
		var temp uint = pts
		for temp > MAX_PTS_VALUE {
			temp -= MAX_PTS_VALUE
		}
		pts = temp
	}
	pts1 := int((pts >> 30) & 0x07)
	pts2 := int((pts >> 15) & 0x7FFF)
	pts3 := int(pts & 0x7FFF)

	data[0] = byte((2 << 4) | (pts1 << 1) | 0x01)
	data[1] = byte((pts2 & 0x7F80) >> 7)
	data[2] = byte(((pts2 & 0x007F) << 1) | 0x01)
	data[3] = byte((pts3 & 0x7F80) >> 7)
	data[4] = byte(((pts3 & 0x007F) << 1) | 0x01)
}
