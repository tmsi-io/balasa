package ts

var tableMade = false

func crc32_block(crc uint32, data []byte, blkLen int) uint32 {

	var i, j int

	if !tableMade {
		makeCrcTable()
	}
	for j = 0; j < blkLen; j++ {
		i = int((uint8(crc>>24) ^ data[j]) & 0xff)
		crc = (crc << 8) ^ crcTable[i]
	}
	return crc
}

var crcTable []uint32 = make([]uint32, 256)

// makeCrcTable Populate the (internal) CRC table. safely be called more than once.
func makeCrcTable() {
	var i, j int
	var crc uint32

	for i = 0; i < 256; i++ {
		crc = uint32(i) << 24
		for j = 0; j < 8; j++ {
			if (crc & 0x80000000) != 0 {
				crc = (crc << 1) ^ CRC32_POLY
			} else {
				crc = crc << 1
			}
		}
		crcTable[i] = crc
	}
	tableMade = true
}
