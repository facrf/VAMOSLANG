//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"vamos/transpiler"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("vamosWasmLoaded", js.ValueOf(true))

	// 1. Transpilar VAMOS -> Go
	js.Global().Set("transpileVamosWasm", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return js.ValueOf(map[string]interface{}{"success": false, "error": "código não informado"})
		}
		source := args[0].String()
		res, err := transpiler.TranspileSource(source)
		if err != nil {
			return js.ValueOf(map[string]interface{}{"success": false, "error": err.Error()})
		}
		return js.ValueOf(map[string]interface{}{"success": true, "output": res})
	}))

	// 2. Formatar código VAMOS
	js.Global().Set("formatVamosWasm", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return js.ValueOf(map[string]interface{}{"success": false, "error": "código não informado"})
		}
		source := args[0].String()
		res, err := transpiler.FormatSource(source)
		if err != nil {
			return js.ValueOf(map[string]interface{}{"success": false, "error": err.Error()})
		}
		return js.ValueOf(map[string]interface{}{"success": true, "output": res})
	}))

	// 3. Conversão Reversa Go -> VAMOS
	js.Global().Set("go2vamosWasm", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return js.ValueOf(map[string]interface{}{"success": false, "error": "código não informado"})
		}
		source := args[0].String()
		res, err := transpiler.TranspileGoToVamos(source)
		if err != nil {
			return js.ValueOf(map[string]interface{}{"success": false, "error": err.Error()})
		}
		return js.ValueOf(map[string]interface{}{"success": true, "output": res})
	}))

	// 4. Validação e Diagnósticos Sintáticos
	js.Global().Set("validateVamosWasm", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return js.ValueOf("[]")
		}
		source := args[0].String()
		diags, _ := transpiler.ValidateSource(source)
		data, _ := json.Marshal(diags)
		return js.ValueOf(string(data))
	}))

	<-c
}
