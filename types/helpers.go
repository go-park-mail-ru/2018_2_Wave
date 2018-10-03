package types

func Must(bytes []byte, err error) []byte {
	if err != nil {
		return []byte{}
	}
	return bytes
}
