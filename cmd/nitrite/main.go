package main

import (
	"encoding/base64"
	"encoding/json"
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
	debug := false

	if len(args) > 1 && args[1].Type() == js.TypeBoolean {
		debug = args[1].Bool()
	}

	currentTime := time.Now()
	if debug {
		currentTime = time.Date(2024, time.March, 5, 16, 0, 0, 0, time.UTC)
	}

	type response struct {
		Result string `json:"result"`
		Err    string `json:"error"`
	}

	document, err := base64.StdEncoding.DecodeString(fDocument)
	if err != nil {
		errorResponse := response{
			Result: "",
			Err:    "Provided attestation document is not encoded as a valid standard Base64 string",
		}
		errorMessage, _ := json.Marshal(errorResponse)
		return js.ValueOf(string(errorMessage))
	}

	res, verifyErr := nitrite.Verify(
		document,
		nitrite.VerifyOptions{
			CurrentTime: currentTime,
		},
	)

	resJSON := ""
	if res != nil {
		enc, marshalErr := json.Marshal(res.Document)
		if marshalErr != nil {
			errorResponse := response{
				Result: "",
				Err:    "Error marshalling the verification result: " + marshalErr.Error(),
			}
			errorMessage, _ := json.Marshal(errorResponse)
			return js.ValueOf(string(errorMessage))
		}
		resJSON = string(enc)
	}

	if verifyErr != nil {
		errorResponse := response{
			Result: "",
			Err:    "Attestation verification failed with error: " + verifyErr.Error(),
		}
		errorMessage, _ := json.Marshal(errorResponse)
		return js.ValueOf(string(errorMessage))
	}

	successResponse := response{
		Result: resJSON,
		Err:    "",
	}
	finalMessage, _ := json.Marshal(successResponse)
	return js.ValueOf(string(finalMessage))
}
