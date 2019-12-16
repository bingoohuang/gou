package pbe_test

import (
	"fmt"
	"testing"

	"github.com/bingoohuang/gou/pbe"
	"github.com/stretchr/testify/assert"
)

func TestChangePBE(t *testing.T) {
	s := `+---+---------+-----------------------------+
| # | PLAIN   | ENCRYPTED                   |
+---+---------+-----------------------------+
| 1 | 1333333 | {PBE}CBX2_bxV5SgOFhPizFrF7A |
+---+---------+-----------------------------+

+---+--------+-----------------------------+
| # | PLAIN  | ENCRYPTED                   |
+---+--------+-----------------------------+
| 1 | 444444 | {PBE}-YBVq_tS3Frr7OtjGIDjXQ |
+---+--------+-----------------------------+
`

	xx, err := pbe.Config{Passphrase: "bingoohuang"}.ChangePbe(s, "bingoohuang123")
	assert.Nil(t, err)
	fmt.Println(xx)
}

func TestFreePBE(t *testing.T) {
	s := `+---+---------+-----------------------------+
| # | PLAIN   | ENCRYPTED                   |
+---+---------+-----------------------------+
| 1 | 1333333 | {PBE}CBX2_bxV5SgOFhPizFrF7A |
+---+---------+-----------------------------+

+---+--------+-----------------------------+
| # | PLAIN  | ENCRYPTED                   |
+---+--------+-----------------------------+
| 1 | 444444 | {PBE}-YBVq_tS3Frr7OtjGIDjXQ |
+---+--------+-----------------------------+
`

	xx, err := pbe.Config{Passphrase: "bingoohuang"}.EbpText(s)
	assert.Nil(t, err)
	fmt.Println(xx)
}

func TestPbeText(t *testing.T) {
	s := `+---+---------+-----------------------------+
| # | PLAIN   | ENCRYPTED                   |
+---+---------+-----------------------------+
| 1 | {PWD:1333333} | dd |
+---+---------+-----------------------------+

+---+--------+-----------------------------+
| # | PLAIN  | ENCRYPTED                   |
+---+--------+-----------------------------+
| 1 | "PWD:444444" | x PWD:444444 x |
+---+--------+-----------------------------+
`

	c := pbe.Config{Passphrase: "bingoohuang"}
	xx, err := c.PbeText(s)
	assert.Nil(t, err)
	fmt.Println(xx)

	yy, err := c.EbpText(xx)
	assert.Nil(t, err)
	fmt.Println(yy)
}
