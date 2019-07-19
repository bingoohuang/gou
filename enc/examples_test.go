package enc_test

import (
	"fmt"
	"github.com/bingoohuang/gou/enc"
	"testing"
)

func TestEnc(t *testing.T) {
	fmt.Println(enc.Base64("不忘初心牢记使命!")) // 5LiN5b-Y5Yid5b-D54mi6K6w5L2_5ZG9IQ
	fmt.Println(enc.Base64Decode("5LiN5b-Y5Yid5b-D54mi6K6w5L2_5ZG9IQ")) // 不忘初心牢记使命!

	fmt.Println(enc.CBCEncrypt("16/24/32bytesxxx", "新时代中国特色社会主义!"))
	fmt.Println(enc.CBCDecrypt("16/24/32bytesxxx", "HK5Ptmtt3V16mIBhJqNeQS_SbTn5kNmE4FSKoxx5t_I9fbIkf2GnjTF6T9KtuWuA8WZYWLMYZeAGsuHyycz9UA=="))
}
