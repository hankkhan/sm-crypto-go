package main

import (
	_ "crypto/rand"
	"net/http"
	"sm-crypto-go/sm2"
	"strconv"
)

func main() {
	http.HandleFunc("/getPublicKey", getPublicKey)
	http.HandleFunc("/decryptData", decryptData)
	http.ListenAndServe(":8080", nil)
}

func decryptData(writer http.ResponseWriter, request *http.Request) {
	encryptText := request.FormValue("pwd")
	cipherMode := request.FormValue("cipherMode")
	cipherModeInt, err := strconv.Atoi(cipherMode)
	if err != nil {
		return
	}
	decrypt := sm2.DoDecrypt(encryptText, cipherModeInt)
	writer.Write([]byte(decrypt))
}

func getPublicKey(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(sm2.GetPublicKey()))
}
