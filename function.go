package rules

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/rwtodd/Go.Sed/sed"
	"github.com/tidwall/gjson"
)

type (
	// Function ..
	Function struct {
		FuncName string
		Args     []string
	}
)

// Run ..
func (f *Function) Run() (string, error) {
	switch f.FuncName {
	case "formula":
		if len(f.Args) < 1 {
			return "", fmt.Errorf("formula: empty parameters")
		}
		stringExp := strings.Join(f.Args, " ")
		exp, errExp := govaluate.NewEvaluableExpression(stringExp)
		if errExp != nil {
			return "", fmt.Errorf("Invalid Formula: %s", stringExp)
		}
		result, err := exp.Evaluate(nil)
		if err != nil {
			return "", fmt.Errorf("Formula can't be processed: %s", stringExp)
		}

		return fmt.Sprint(result), nil
	case "md5":
		md5 := md5.Sum([]byte(f.Args[0]))
		return fmt.Sprintf("%x", md5), nil
	case "sha1":
		enc := sha1.New()
		enc.Write([]byte(f.Args[0]))
		encResult := enc.Sum(nil)
		return fmt.Sprintf("%x", encResult), nil
	case "sha256":
		sha256 := sha256.Sum256([]byte(f.Args[0]))
		return fmt.Sprintf("%x", sha256), nil
	case "hs256":
		if len(f.Args) < 1 {
			return "", fmt.Errorf("hs256: first argument (data) can't be empty")
		}
		lenArgs := len(f.Args)
		if lenArgs < 2 {
			return "", fmt.Errorf("hs256: second argument (secret key) can't be empty")
		}

		src := f.Args[0]
		secret := f.Args[1]
		key := []byte(secret)
		h := hmac.New(sha256.New, key)
		h.Write([]byte(src))
		hs256 := base64.StdEncoding.EncodeToString(h.Sum(nil))
		return fmt.Sprintf("%v", hs256), nil
	case "base64enc":
		base64 := base64.StdEncoding.EncodeToString([]byte(f.Args[0]))
		return fmt.Sprintf("%v", base64), nil
	case "base64dec":
		base64, err := base64.StdEncoding.DecodeString(f.Args[0])
		return fmt.Sprintf("%s", base64), err
	case "rsaGenerateSign":
		if f.Args[0] == "" {
			return "", fmt.Errorf("rsaGenerateSign: first argument (data) can't be empty")
		}
		lenArgs := len(f.Args)
		if lenArgs < 2 {
			return "", fmt.Errorf("rsaGenerateSign: second argument (private key) can't be empty")
		}

		// skip when parsing for log
		if f.Args[1] == "*****" {
			return "*****", nil
		}

		// Use bufer to trim space /t /n
		buffer := new(bytes.Buffer)
		json.Compact(buffer, []byte(f.Args[0]))
		payload := fmt.Sprintf("%v", buffer)
		sEnc := ""
		// Load key file.

		key := []byte(f.Args[1])
		block, _ := pem.Decode(key)
		if block == nil {
			return "", fmt.Errorf("rsaGenerateSign: Failed to decode key file")
		}
		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes) //x509.ParseCertificate(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("rsaGenerateSign: Failed to parse private key")
		}
		hashed := sha256.Sum256([]byte(payload))
		s, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
		if err != nil {
			return "", fmt.Errorf("rsaGenerateSign: Failed to generate signature")
		}
		sEnc = base64.StdEncoding.EncodeToString(s)
		return sEnc, nil
	case "rsaVerifySign":
		if f.Args[0] == "" {
			return "", fmt.Errorf("rsaVerifySign: first argument (data) can't be empty")
		}
		lenArgs := len(f.Args)
		if lenArgs < 3 {
			return "", fmt.Errorf("rsaVerifySign: second & third argument (public key & signature) can't be empty")
		}

		// skip when parsing for log
		if f.Args[1] == "*****" {
			return "*****", nil
		}

		// Use bufer to trim space /t /n
		buffer := new(bytes.Buffer)
		json.Compact(buffer, []byte(f.Args[0]))
		response := fmt.Sprintf("%v", buffer)

		block, _ := pem.Decode([]byte(f.Args[1]))
		if block == nil {
			return "", fmt.Errorf("rsaVerifySign: Failed to decode pem file key")
		}
		pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("rsaVerifySign: Failed to parse public key")
		}
		sDec, err := base64.StdEncoding.DecodeString(f.Args[2])
		if err != nil {
			return "", fmt.Errorf("rsaVerifySign: Failed to decode signature")
		}
		hashed := sha256.Sum256([]byte(response))
		err = rsa.VerifyPKCS1v15(pubkey.(*rsa.PublicKey), crypto.SHA256, hashed[:], sDec)
		return "", err
	case "length":
		length := len(f.Args[0])
		return fmt.Sprint(length), nil
	case "lowercase":
		return strings.ToLower(f.Args[0]), nil
	case "uppercase":
		return strings.ToUpper(f.Args[0]), nil
	case "replace":
		if len(f.Args) < 1 {
			return "", fmt.Errorf("replace: first argument can't be empty")
		}
		if len(f.Args) < 2 {
			return "", fmt.Errorf("replace: second argument can't be empty")
		}
		if len(f.Args) < 3 {
			return "", fmt.Errorf("replace: third argument can't be empty")
		}

		result := strings.ReplaceAll(f.Args[0], f.Args[1], f.Args[2])
		return result, nil
	case "sed":
		if len(f.Args) < 1 {
			return "", fmt.Errorf("sed: first argument (data) can't be empty")
		}
		if len(f.Args) < 2 {
			return "", fmt.Errorf("sed: second argument (expressions) can't be empty")
		}

		exps := strings.Split(f.Args[1], "  ")

		result := strings.TrimSpace(f.Args[0])
		engine, err := sed.New(strings.NewReader(exps[0]))
		exp := ""

		for i := range exps {
			exp = strings.TrimSpace(exps[i])
			engine, err = sed.New(strings.NewReader(exp))
			if err != nil {
				return "", fmt.Errorf("sed: init expression failed, " + err.Error())
			}
			result, err = engine.RunString(result)
			if err != nil {
				return "", fmt.Errorf("sed: run expression failed, " + err.Error())
			}
			result = strings.TrimRight(result, "\r\n")
		}

		return result, nil
	case "date_format":
		lenArgs := len(f.Args)
		if len(f.Args) < 1 {
			return "", fmt.Errorf("date_format: argument can't be empty")
		}
		if lenArgs < 1 {
			return "", fmt.Errorf("date_format: argument can't be empty")
		}
		if lenArgs > 1 && lenArgs != 3 {
			return "", fmt.Errorf("date_format: maximum argument 3")
		}
		if f.Args[0] == "" {
			return "", fmt.Errorf("date_format: first argument can't be empty")
		}
		if lenArgs == 3 {
			if f.Args[1] == "" || f.Args[2] == "" {
				return "", fmt.Errorf("date_format: second and third argument can't be empty")
			}
		}
		_time := time.Now()
		if lenArgs == 3 {
			var err error
			_time, err = DateTimeParse(f.Args[1], f.Args[2])
			if err != nil {
				return "", fmt.Errorf("date_format: invalid parse format time %s -> %s", f.Args[1], f.Args[2])
			}
		}
		return DateTimeFormat(f.Args[0], _time), nil
	case "array_map":
		var firstParam []interface{}
		err := json.Unmarshal([]byte(f.Args[0]), &firstParam)
		if err != nil {
			return "", fmt.Errorf("array_map: invalid first argument %s", fmt.Sprint(f.Args[0]))
		}
		arrCount := len(firstParam)
		if arrCount == 0 {
			return "", fmt.Errorf("array_map: array first argument can't be empty")
		}
		secondParam := f.Args[1]
		if secondParam == "" {
			return "", fmt.Errorf("array_map: second param can't be empty")
		}

		var errMsg string
		var result string

		dataToMap := getAllPath(secondParam)
		var body = map[string]interface{}{
			"body": firstParam,
		}
		_body, _ := json.Marshal(body)
		_firstParam := string(_body)

	firstLoop:
		for idx, row := range firstParam {
			copyData := secondParam
			_idx := strconv.Itoa(idx)
			if len(dataToMap) > 0 {
				for _, key := range dataToMap {
					var data, _keyTemp string

					_key := strings.TrimSpace(key)
					_hasFunc := strings.Index(_key, "|")
					if _hasFunc > -1 {
						_keyTemp = _key
						_key = _key[0:_hasFunc]
					}

					var value gjson.Result
					if _key == "_row_" {
						value = gjson.Get(_firstParam, "body."+_idx)

					} else if strings.HasPrefix(_key, "_row_.") {
						_key = _key[6:len(_key)]
						value = gjson.Get(_firstParam, "body."+_idx+"."+_key)
					}
					if !value.Exists() {
						errMsg = fmt.Sprintf("array_map: key [%s] on index-%s doesn't exist -> %s", _key, _idx, fmt.Sprint(row))
						break firstLoop
					}
					data = value.String()

					if _hasFunc > -1 {
						_key = _keyTemp
						_func := _key[_hasFunc+1:]
						if _func != "" {
							var _funcName string
							var _funcArgs []string
							_funcIdx := strings.Index(_func, ",")
							_funcName = strings.TrimSpace(_func)
							if _funcIdx > -1 {
								_funcName = strings.TrimSpace(_func[0:_funcIdx])
								_funcArgs = append(_funcArgs, strings.TrimSpace(_func[_funcIdx+1:]))
								_funcArgs = append(_funcArgs, strings.Replace(_funcArgs[0], "_VALUE_", value.String(), -1))
							}
							f := Function{
								FuncName: _funcName,
								Args:     _funcArgs,
							}
							funcRes, err := f.Run()
							if err == nil {
								data = funcRes
							}
						}
					}
					copyData = strings.Replace(copyData, "["+key+"]", data, 1)
				}
			}
			if idx != (arrCount - 1) {
				result += copyData + ","
			} else {
				result += copyData
			}
		}
		if errMsg != "" {
			return "", fmt.Errorf(errMsg)
		}

		result = "[" + result + "]"
		var resultCheck []interface{}
		err = json.Unmarshal([]byte(result), &resultCheck)
		if err != nil {
			return "", fmt.Errorf("array_map: invalid JSON -> %s", result)
		}
		return result, nil
	case "array_join":
		if f.Args[0] == "" {
			return "", fmt.Errorf("array_join: invalid first argument %s", fmt.Sprint(f.Args[0]))
		}

		var arrData []string
		err := json.Unmarshal([]byte(f.Args[0]), &arrData)
		if err != nil {
			return "", fmt.Errorf("array_join: invalid array format -> %s", f.Args[0])
		}

		result := strings.Join(arrData, f.Args[1])
		return result, nil
	case "split":
		if f.Args[0] == "" {
			return "", fmt.Errorf("split: invalid first argument %s", fmt.Sprint(f.Args[0]))
		}
		if f.Args[1] == "" {
			f.Args[1] = ","
		}

		res := strings.Split(f.Args[0], f.Args[1])
		result, _ := json.Marshal(res)

		return string(result), nil
	case "substring":
		if f.Args[0] == "" {
			return "", fmt.Errorf("substring: first argument can't be empty")
		}
		lenArgs := len(f.Args)
		if lenArgs < 2 {
			return "", fmt.Errorf("substring: second argument can't be empty")
		}

		first := []rune(f.Args[0]) //support for special char
		lenFirstArgs := len(first)
		second, _ := strconv.Atoi(f.Args[1])
		third := 0

		if second < 0 || second >= lenFirstArgs {
			return "", fmt.Errorf("substring: second argument index out of range")
		}
		if lenArgs > 2 {
			third, _ = strconv.Atoi(f.Args[2])
			if third < 1 || third > lenFirstArgs {
				return "", fmt.Errorf("substring: third argument index out of range")
			}
			if second == third {
				return "", fmt.Errorf("substring: second and third argument can't be the same")
			}
			if second >= third {
				return "", fmt.Errorf("substring: second argument index out of range")
			}
		}

		var result string
		if third < 1 {
			result = string(first[second:])
		} else {
			result = string(first[second:third])
		}
		return result, nil

	}
	return "", fmt.Errorf("Function doesn't exist: %s", f.FuncName)
}
