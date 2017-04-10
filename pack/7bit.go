package pack

func Pack7Bit(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	result := make([]byte, 0, len(data))

	i := uint(1)
	for ; i < uint(len(data)); i++ {
		bits_hi := i % 8

		// skip iteration when 7 bytes were generated
		// we already store data from 8 characters in them
		if bits_hi == 0 {
			continue
		}

		mask_hi := (byte(0xff) >> (8 - bits_hi))
		mask_lo := byte(0xff) >> bits_hi

		high := (data[i] & mask_hi) << (8 - bits_hi)

		// shift low bits to make place for high
		low := (data[i-1] >> (i%8 - 1)) & mask_lo

		result = append(result, high|low)
	}

	// generate last output byte if needed
	if i%8 != 0 {
		mask_lo := byte(0xff) >> (i % 8)
		result = append(result, (data[i-1]>>(i%8-1))&mask_lo)
	}

	return result
}

func Unpack7Bit(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	result := make([]byte, 0, 2*len(data))
	for i := uint(0); i < uint(len(data)); i++ {
		// number of bits to extract for low and high part of character
		bits_lo := i % 7
		bits_hi := 7 - bits_lo

		mask_lo := byte(0xff) >> (8 - bits_lo)
		mask_hi := byte(0xff) >> (8 - bits_hi)

		high := (data[i] & mask_hi) << bits_lo

		if bits_lo == 0 {
			result = append(result, high)
		} else {
			low := (data[i-1] >> (8 - bits_lo)) & mask_lo
			result = append(result, high|low)
		}

		if bits_lo == 6 {
			low := data[i] >> 1
			result = append(result, low)
		}
	}

	if result[len(result)-1] == 0 {
		return result[:len(result)-1]
	} else {
		return result
	}
}
