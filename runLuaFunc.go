package main

import (
	"bytes"
	"encoding/json"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func runLuaFunctionWithContext(file string, funcCall string, req SimulateRequest) (string, error) {
	L := lua.NewState()
	defer L.Close()

	// Ustaw globalną zmienną "request"
	reqTable := L.NewTable()
	L.SetField(reqTable, "method", lua.LString(req.Method))
	L.SetField(reqTable, "url", lua.LString(req.Endpoint))
	L.SetField(reqTable, "body", lua.LString(req.Body))

	headersTable := L.NewTable()
	for k, v := range req.Headers {
		L.SetField(headersTable, k, lua.LString(v))
	}
	L.SetField(reqTable, "headers", headersTable)
	L.SetGlobal("request", reqTable)

	if err := L.DoFile(file); err != nil {
		return "", err
	}

	funcName := strings.TrimSuffix(strings.TrimSuffix(funcCall, "()"), ";")
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal(funcName),
		NRet:    1,
		Protect: true,
	}); err != nil {
		return "", err
	}

	ret := L.Get(-1)
	L.Pop(1)

	if tbl, ok := ret.(*lua.LTable); ok {
		response := tbl.RawGetString("response")
		logs := tbl.RawGetString("log")

		// konwersja do JSON
		out := map[string]interface{}{
			"response": luaValueToInterface(response),
			"log":      luaValueToInterface(logs),
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(out)
		return buf.String(), nil
	}

	return ret.String(), nil
}

func luaValueToInterface(val lua.LValue) interface{} {
	switch v := val.(type) {
	case *lua.LTable:
		result := map[string]interface{}{}
		v.ForEach(func(key, value lua.LValue) {
			result[key.String()] = luaValueToInterface(value)
		})
		// sprawdź, czy to może być tablica
		if isArrayLike(v) {
			arr := []interface{}{}
			v.ForEach(func(_, value lua.LValue) {
				arr = append(arr, luaValueToInterface(value))
			})
			return arr
		}
		return result
	default:
		return v.String()
	}
}

func isArrayLike(tbl *lua.LTable) bool {
	i := 1
	tbl.ForEach(func(_, value lua.LValue) {
		if tbl.RawGetInt(i) == lua.LNil {
			i = -1 // Mark as invalid
		}
		i++
	})
	if i == -1 {
		return false
	}
	return true
}
