package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"syscall/js"
	"time"

	"github.com/hf/nitrite"
)

func main() {
	js.Global().Set("validateAttestation", js.FuncOf(validateAttestation))
	<-make(chan bool)
}

func validateAttestation(this js.Value, args []js.Value) interface{} {
	fDocument := args[0].String()

	document, err := base64.StdEncoding.DecodeString(fDocument)
	if nil != err {
		fmt.Printf("Provided attestation document is not encoded as a valid standard Base64 string\n")
		os.Exit(2)
	}

	res, err := nitrite.Verify(
		document,
		nitrite.VerifyOptions{
			CurrentTime: time.Now(),
		},
	)

	resJSON := ""

	if nil != res {
		enc, err := json.Marshal(res.Document)
		if nil != err {
			panic(err)
		}

		resJSON = string(enc)
	}

	if nil != err {
		fmt.Printf("Attestation verification failed with error %v\n", err)
		os.Exit(2)
	}

	return js.ValueOf(resJSON)
}
