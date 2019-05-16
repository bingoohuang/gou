package gou_test

func TestChecksum(t *testing.T) {
	crc := gou.Checksum([]byte("bigoohuang"))
	if crc != "380372004" {
		t.Error("Checksum failed")
	}
}
