package types

func Shield(bytes []byte, err error) []byte {
	if err != nil {
		return []byte{}
	}
	return bytes
}
