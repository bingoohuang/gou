package gou

// Checksum 返回data的校验和
func Checksum(data []byte) string {
	return fmt.Sprint(adler32.Checksum(data))
}