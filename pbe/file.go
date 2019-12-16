package pbe

import (
	"fmt"
	"regexp"
	"strings"
)

// Config configs the passphrase.
type Config struct {
	Passphrase string
}

// Pbe encrypts p by PBEWithMD5AndDES with 19 iterations.
// it will prompt password if viper get none.
func (c Config) Pbe(p string) (string, error) {
	pwd := c.Passphrase
	if pwd == "" {
		pwd = GetPbePwd()
	}

	if pwd == "" {
		return "", fmt.Errorf("pbepwd is requird")
	}

	encrypt, err := Encrypt(p, pwd, iterations)
	if err != nil {
		return "", err
	}

	return pbePrefix + encrypt, nil
}

// Ebp decrypts p by PBEWithMD5AndDES with 19 iterations.
func (c Config) Ebp(p string) (string, error) {
	if !strings.HasPrefix(p, pbePrefix) {
		return p, nil
	}

	pwd := c.Passphrase
	if pwd == "" {
		pwd = GetPbePwd()
	}

	if pwd == "" {
		return "", fmt.Errorf("pbepwd is requird")
	}

	return Decrypt(p[len(pbePrefix):], pwd, iterations)
}

var pbeRe = regexp.MustCompile(`\{PBE\}[\w_-]+`) // nolint

// ChangePbe changes the {PBE}xxx to {PBE} yyy with a new passphase
func (c Config) ChangePbe(s, newPassphrase string) (string, error) {
	var err error

	m := make(map[string]string)
	f := func(old string) string {
		if v, ok := m[old]; ok {
			return v
		}

		raw := ""
		raw, err = c.Ebp(old)

		if err != nil {
			return ""
		}

		newPwd := ""
		newPwd, err = Config{Passphrase: newPassphrase}.Pbe(raw)

		if err != nil {
			return ""
		}

		m[old] = newPwd

		return newPwd
	}

	return pbeRe.ReplaceAllStringFunc(s, f), err
}

// EbpText free the {PBE}xxx to yyy with a  passphrase
func (c Config) EbpText(s string) (string, error) {
	var err error

	m := make(map[string]string)

	f := func(old string) string {
		if v, ok := m[old]; ok {
			return v
		}

		raw := ""
		raw, err = c.Ebp(old)

		if err != nil {
			return ""
		}

		m[old] = raw

		return raw
	}

	return pbeRe.ReplaceAllStringFunc(s, f), err
}

// PbeText will PBE encrypt the passwords in the text
// passwords should be as any of following format and its converted pattern
// 1. {PWD:clear} -> {PBE:cyphered}
// 2. [PWD:clear] -> {PBE:cyphered}
// 3. (PWD:clear) -> {PBE:cyphered}
// 4. "PWD:clear" -> "{PBE:cyphered}"
// 5.  PWD:clear  ->  {PBE:cyphered}
func (c Config) PbeText(s string) (string, error) {
	m := make(map[string]string)

	pbed := ""
	src := s

	var err error

	for {
		pos := strings.Index(src, "PWD:")
		if pos <= 0 {
			pbed += src
			break
		}

		left := src[pos-1]
		if Alphanumeric(left) {
			pbed += "PWD:"
			src = src[4:]

			continue
		}

		expectRight := string(getRight(left))
		rpos := strings.Index(src[pos+4:], expectRight)

		if rpos < 0 {
			pbed += src[0:4]
			src = src[4:]

			continue
		}

		raw := src[pos+4 : pos+rpos+4]

		pwd := ""
		if v, ok := m[raw]; ok {
			pwd = v
		} else {
			if pwd, err = c.Pbe(raw); err != nil {
				return "", err
			}

			m[raw] = pwd
		}

		switch left {
		case '(', '{', '[':
		default:
			pwd = string(left) + pwd + expectRight
		}

		pbed += src[0:pos-1] + pwd
		src = src[pos+rpos+5:]
	}

	return pbed, nil
}

func getRight(left uint8) uint8 {
	switch left {
	case '(':
		return ')'
	case '{':
		return '}'
	case '[':
		return ']'
	default:
		return left
	}
}

// Alphanumeric tells u is letter or digit char.
func Alphanumeric(u uint8) bool {
	return u >= '0' && u <= '9' || u >= 'a' && u <= 'z' || u >= 'A' && u <= 'Z'
}
