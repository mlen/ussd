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

	// length of the output sequence
	// 7 bit encoding packs 8 characters in 7 bytes
	// we need to calculate the output size based on that
	length := (uint(len(data)) * 8) / 7

	result := make([]byte, 0, length)

	// first iteration is straightforward, hardcode it
	result = append(result, data[0]&0x7f)
	for i := uint(1); i < length; i++ {
		// during each cycle 8 characters are generated
		// after each cycle, we need to go back 1 character in the input sequence
		// because it still contains data for the next cycle
		cycle := i / 8

		// number of bits to extract for low and high part of character
		bits_lo := i % 8
		bits_hi := 7 - bits_lo

		mask_lo := byte(0xff) >> (8 - bits_lo)
		mask_hi := byte(0xff) >> (8 - bits_hi)

		low := (data[i-cycle-1] >> (8 - bits_lo)) & mask_lo
		high := ((data[i-cycle] & mask_hi) << bits_lo) // FIXME there is a bug

		result = append(result, high|low)
	}

	return result
}
