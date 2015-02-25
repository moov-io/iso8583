package iso8583

func lbcd(data []byte) []byte {
	if len(data)%2 != 0 {
		return bcd(append(data, "\x00"...))
	}
	return bcd(data)
}

func rbcd(data []byte) []byte {
	if len(data)%2 != 0 {
		return bcd(append([]byte("\x00"), data...))
	}
	return bcd(data)
}

// Encode numeric in ascii into bsd (be sure len(data) % 2 == 0)
func bcd(data []byte) []byte {
	if len(data)%2 != 0 {
		panic("length of raw data must be even")
	}
	out := make([]byte, len(data)/2)
	for i, j := 0, 0; i < len(out); i++ {
		out[i] = ((data[j] & 0x0f) << 4) | (data[j+1] & 0x0f)
		j += 2
	}
	return out
}

func bcdl2Ascii(data []byte, length int) []byte {
	return bcd2Ascii(data)[:length]
}

func bcdr2Ascii(data []byte, length int) []byte {
	out := bcd2Ascii(data)
	return out[len(out)-length:]
}

func bcd2Ascii(data []byte) []byte {
	outLen := len(data) * 2
	out := make([]byte, outLen)
	for i := 0; i < outLen; i++ {
		bcdIndex := i / 2
		if i%2 == 0 {
			// higher order bits to ascii:
			out[i] = (data[bcdIndex] >> 4) + '0'
		} else {
			// lower order bits to ascii:
			out[i] = (data[bcdIndex] & 0x0f) + '0'
		}
	}
	return out
}
