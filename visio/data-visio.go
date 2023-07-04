package visio

import "time"

type vdata struct {
	tcycle int
	ton    [16]int
	toff   [16]int
	yellow [16]bool
	prev   [16]bool
}

func (v *vdata) init() {
	v.tcycle = 0
	for i := 0; i < 16; i++ {
		v.toff[i] = 0
		v.ton[i] = 0
		v.yellow[i] = false
		v.prev[i] = false
	}
}
func (v *vdata) set(tcycle int, naps [16]bool) {
	v.tcycle = tcycle
	for i := 0; i < 16; i++ {
		v.ton[i] = 0
		v.toff[i] = 0
		v.yellow[i] = false

	}
	for i := 0; i < 16; i++ {
		if naps[i] {
			//Направление включается
			v.toff[i] = tcycle
		} else {
			v.ton[i] = tcycle
		}
		if naps[i] != v.prev[i] {
			v.yellow[i] = true
		}
	}
	v.prev = naps
}
func (v vdata) makeSpecial(c int, state int) [29]byte {
	var buffer [29]byte
	buffer[4] = byte(c)
	buffer[5] = 6 //АПП
	buffer[6] = byte(state)
	buffer[28] = crc(buffer)
	return buffer
}
func (v vdata) makeBuffer(c int) [29]byte {
	var buffer [29]byte
	buffer[0] = 0
	buffer[1] = 14
	buffer[2] = 10
	buffer[3] = 25
	buffer[4] = byte(c)
	buffer[5] = 6 //АПП
	buffer[6] = 3 //Фаза
	buffer[7] = byte(v.tcycle & 0xff)
	s := 0
	if c == 168 {
		s = 8
	}
	for i := 0; i < 8; i++ {
		buffer[8+i] = byte(v.ton[i+s] & 0xff)
	}
	for i := 0; i < 8; i++ {
		buffer[16+i] = byte(v.toff[i+s] & 0xff)
	}
	y := 0
	for i := 0; i < 8; i++ {
		if v.yellow[i+s] {
			y += 1
		}
		y = y << 1
	}
	buffer[24] = byte(y & 0xff)
	buffer[25] = byte(time.Now().Minute())
	buffer[26] = byte(time.Now().Hour())
	buffer[27] = byte(time.Now().Weekday())
	buffer[28] = crc(buffer)
	return buffer
}
func crc(buffer [29]byte) byte {
	var cr int64 = 0
	for i := 0; i < 28; i++ {
		cr += int64(buffer[i])
	}
	return byte(cr & 0xff)
}
