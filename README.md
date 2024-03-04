Nitrite WASM
=======
Build:
```
GOOS=js GOARCH=wasm go build -o nitrite.wasm main.go
```

To run example inside `/dist` you can use
`python3 -m http.server`. 

Then `validateAttestation(base64document: string)` should be available globally in the JS console.

In Fedora, `wasm_exec.js` is found in `/usr/lib/golang/misc/wasm` after installing the `golang-misc` package. 