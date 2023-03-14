package main

import (
	_ "embed"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go/v6"
	"log"
	"os"
)

//go:embed haskell.wasm
var wasm []byte

func main() {
	stdoutPath := "/tmp/log.txt"

	engine := wasmtime.NewEngine()
	module, err := wasmtime.NewModule(engine, wasm)
	if err != nil {
		log.Fatal(err)
	}

	// Create a linker with WASI functions defined within it
	linker := wasmtime.NewLinker(engine)
	err = linker.DefineWasi()
	if err != nil {
		log.Fatal(err)
	}

	// Configure WASI imports to write stdout into a file, and then create
	// a `Store` using this wasi configuration.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdoutFile(stdoutPath)
	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		log.Fatal(err)
	}

	// Run the function
	nom := instance.GetFunc(store, "hs_init")
	_, err = nom.Call(store, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	nom = instance.GetFunc(store, "query")
	_, err = nom.Call(store)
	if err != nil {
		log.Fatal(err)
	}

	// Print WASM stdout
	out, err := os.ReadFile(stdoutPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(out))

}
