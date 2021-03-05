package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionNotFound(t *testing.T) {
	funcArgs := []string{"abc|cde", "|"}

	f := Function{
		FuncName: "notfound",
		Args:     funcArgs,
	}
	_, err := f.Run()
	assert.Equal(t, "Function doesn't exist: notfound", err.Error())
}

func TestFunctionSplit(t *testing.T) {
	funcArgs := []string{"abc|cde", "|"}

	f := Function{
		FuncName: "split",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "[\"abc\",\"cde\"]", funcRes)
}

func TestFunctionSubString(t *testing.T) {
	funcArgs := []string{"abcdefgh", "0", "6"}

	f := Function{
		FuncName: "substring",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "abcdef", funcRes)
}

func TestFunctionFormula(t *testing.T) {
	funcArgs := []string{"20", "/", "4"}

	f := Function{
		FuncName: "formula",
		Args:     funcArgs,
	}
	funcRes, err := f.Run()
	assert.Equal(t, "5", funcRes)

	funcArgs = []string{}

	f = Function{
		FuncName: "formula",
		Args:     funcArgs,
	}
	funcRes, err = f.Run()
	assert.Equal(t, "formula: empty parameters", err.Error())

	funcArgs = []string{"20", "k", "4"}
	f = Function{
		FuncName: "formula",
		Args:     funcArgs,
	}
	funcRes, err = f.Run()
	assert.Equal(t, "Invalid Formula: 20 k 4", err.Error())
}

func TestFunctionOneArgs(t *testing.T) {

	funcArgs := []string{"abc"}
	f := Function{
		FuncName: "md5",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "900150983cd24fb0d6963f7d28e17f72", funcRes)

	f = Function{
		FuncName: "sha1",
		Args:     funcArgs,
	}
	funcRes, _ = f.Run()
	assert.Equal(t, "a9993e364706816aba3e25717850c26c9cd0d89d", funcRes)

	f = Function{
		FuncName: "sha256",
		Args:     funcArgs,
	}
	funcRes, _ = f.Run()
	assert.Equal(t, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", funcRes)

	f = Function{
		FuncName: "base64enc",
		Args:     funcArgs,
	}
	funcRes, _ = f.Run()
	assert.Equal(t, "YWJj", funcRes)

	f = Function{
		FuncName: "length",
		Args:     funcArgs,
	}
	funcRes, _ = f.Run()
	assert.Equal(t, "3", funcRes)

	f = Function{
		FuncName: "lowercase",
		Args:     funcArgs,
	}
	funcRes, _ = f.Run()
	assert.Equal(t, "abc", funcRes)

	f = Function{
		FuncName: "uppercase",
		Args:     funcArgs,
	}
	funcRes, _ = f.Run()
	assert.Equal(t, "ABC", funcRes)

	funcArgs = []string{"YWJj"}
	f = Function{
		FuncName: "base64dec",
		Args:     funcArgs,
	}
	funcRes, _ = f.Run()
	assert.Equal(t, "abc", funcRes)
}

func TestFunctionHS256(t *testing.T) {

	funcArgs := []string{"abc", "secret"}
	f := Function{
		FuncName: "hs256",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "mUba1OAOkT/Ivo5dP34RCkqegy+D+wnDRShdeGONig4=", funcRes)

	funcArgs = []string{}
	f = Function{
		FuncName: "hs256",
		Args:     funcArgs,
	}
	funcRes, err := f.Run()
	assert.Equal(t, "hs256: first argument (data) can't be empty", err.Error())

	funcArgs = []string{"abc"}
	f = Function{
		FuncName: "hs256",
		Args:     funcArgs,
	}
	funcRes, err = f.Run()
	assert.Equal(t, "hs256: second argument (secret key) can't be empty", err.Error())

}

func TestFunctionThreeArgs(t *testing.T) {

	funcArgs := []string{"abcdef", "ab", "11"}
	f := Function{
		FuncName: "replace",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "11cdef", funcRes)

}

func TestFunctionDateFormat(t *testing.T) {

	funcArgs := []string{"2014-04-04", "YYYY-MM-DD", "DD-YYYY-MM"}
	f := Function{
		FuncName: "date_format",
		Args:     funcArgs,
	}
	funcRes, err := f.Run()
	assert.Nil(t, err)
	assert.Equal(t, "04-2014-04", funcRes)

}

func TestFunctionSed(t *testing.T) {

	funcArgs := []string{"abcdef", "s/abc/def/  s/ef/jj/"}
	f := Function{
		FuncName: "sed",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "djjdef", funcRes)

}

func TestFunctionRSAGenerate(t *testing.T) {
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIICWQIBAAKBgFLyE8GuDMhNVZhZBzWKLb2+i3wqUj2eMKGbRSKe4OLXibNkE5Mz
HbgLzZD7TdOdKKwDiBTqtgC2+tH5Wf3g89qG0sL9JXF9OcxK28B4FHsOlhXKhYpR
93Za+p9KYGaJ4XF5ZMIDfDxscGwASeiMg3pH2vvOP5ejWuyVGEHaHkULAgMBAAEC
gYASHHWf5sc3vVshRt9CG4fdVIvUctE+TxpDT0oLQzHLllCk8QctLw4gL8OVEqpt
uHU3ChZeqtlO0qV1z8KMot/b6W0y2At4UH2KadpjtxS+7SEnxtQuTQBfn59KpVQx
6g2+L633fkHLez5tQLwsH8kh4OI7oyWU1m703dpuWc3n0QJBAKJz+Q9MBfj8Gi3c
rNsVqK9JCv8KnWlRd57QJgyhUWysZlzaUOJMPuyxFpSSU/d2O3DUC5+r1JNhycil
pP2ZxaUCQQCCtYM+NCZI2aQZQtB3GNVFyfeCmUgf4SB5kILRpHnvc6cVT4dF37TZ
RL7OMRI0I6Y0qMIS3MBzQ5rb6rhDl8DvAkBxgRI9i+KIaqxn6s2jbWikwCY8uE/v
bApmHgzXukbH5VTH/4mP87HrcnfSasLcHfG+DYnpkAdAyoxP8txqjGw5Aj8HJeYH
gNKXKU/QEddUrAb9yg2/FqLbG3SrMTv2OwhwD+MTR0YejB1XxGqq3AQi1dBBEPmM
DoZ3xzqwzCVHjQUCQHDnpB9tuIk7cPE0tjg9v0GrPj9RiZWnbm3K5/LdNwe/tqlQ
5+IKwVzzMgeBehNa+UvwCx2QoeFlPgvTWZmyqvc=
-----END RSA PRIVATE KEY-----`
	funcArgs := []string{"abcdef", privateKey}
	f := Function{
		FuncName: "rsaGenerateSign",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "EdDhV1IFyip34AKb6KyjQtQRbkEx1nqm6K+kt/ybVKwfclBcOoGhJb13sA30u6QUAIejtqMLYhmNHEQJxTgUTm/hVcPQ94B3lR0K9/PqaGFiFqoOaE5JW0PQItdDFeWYOyXNtpWZN/1sY9MuxBtoN583dFgWJgp5gNO3Z1YfbZU=", funcRes)
}

func TestFunctionRSAVerify(t *testing.T) {
	publicKey := `-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgFLyE8GuDMhNVZhZBzWKLb2+i3wq
Uj2eMKGbRSKe4OLXibNkE5MzHbgLzZD7TdOdKKwDiBTqtgC2+tH5Wf3g89qG0sL9
JXF9OcxK28B4FHsOlhXKhYpR93Za+p9KYGaJ4XF5ZMIDfDxscGwASeiMg3pH2vvO
P5ejWuyVGEHaHkULAgMBAAE=
-----END PUBLIC KEY-----`
	signature := "EdDhV1IFyip34AKb6KyjQtQRbkEx1nqm6K+kt/ybVKwfclBcOoGhJb13sA30u6QUAIejtqMLYhmNHEQJxTgUTm/hVcPQ94B3lR0K9/PqaGFiFqoOaE5JW0PQItdDFeWYOyXNtpWZN/1sY9MuxBtoN583dFgWJgp5gNO3Z1YfbZU="
	funcArgs := []string{"abcdef", publicKey, signature}
	f := Function{
		FuncName: "rsaVerifySign",
		Args:     funcArgs,
	}
	_, err := f.Run()
	assert.Nil(t, err)
}
