package sm2

import (
	"fmt"
	"io/ioutil"
	"os"
	"rogchap.com/v8go"
	"strings"
)

var privateKey, publicKey string

// DefaultSm2Path sm.js路径
var defaultSm2Path string

// GetPublicKey 获取公钥
func GetPublicKey() string {
	if publicKey != "" {
		return publicKey
	}
	privateKey, publicKey = generateKeyPairHexFunc()
	return publicKey
}

// DoDecrypt 解密
func DoDecrypt(encryptData string, cipherMode int) string {
	var iso = v8go.NewIsolate()
	var ctx = v8go.NewContext(iso)
	//预编译 vm2
	script, _ := iso.CompileUnboundScript(scriptStr(defaultSm2Path), "math.js", v8go.CompileOptions{}) // compile script to get cached data
	_, err := script.Run(ctx)
	if err != nil {
		fmt.Errorf("script exec error,err:%v", err)
		return ""
	}
	code := fmt.Sprintf("const res = doDecrypt(\"%s\",\"%s\",%d)", encryptData, privateKey, cipherMode)
	ctx.RunScript(code, "sm2.js")
	result, _ := ctx.RunScript("res", "doDecrypt.js")
	return result.String()
}

func generateKeyPairHexFunc() (privateKey, publicKey string) {
	// 创建一个新的 V8 虚拟机实例和上下文对象
	var iso = v8go.NewIsolate()
	var ctx = v8go.NewContext(iso)
	// 预编译 vm1 JavaScript 代码
	script, _ := iso.CompileUnboundScript(scriptStr(defaultSm2Path), "math.js", v8go.CompileOptions{}) // compile script to get cached data
	_, err := script.Run(ctx)
	if err != nil {
		fmt.Errorf("script exec error,err:%v", err)
		return
	}
	ctx.RunScript("const result = generateKeyPairHex()", "sm2.js")
	keyPair, _ := ctx.RunScript("result", "value.js")
	fmt.Printf("result:%s", keyPair.String())
	split := strings.Split(keyPair.String(), " ")
	privateKey = split[0]
	publicKey = split[1]
	return
}

func scriptStr(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// 定义一个 JavaScript 函数，计算两个数的和
	jsCode := fmt.Sprintf(`%s`, file)
	return jsCode
}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defaultSm2Path = fmt.Sprintf("%s/sm2.js", dir)
}

// SetSm2FilePath 设置sm路径
func SetSm2FilePath(path string) {
	defaultSm2Path = path
}
