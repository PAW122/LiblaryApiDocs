package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func runLuaFunctionWithContext(file string, funcCall string, req SimulateRequest, defaultDB []map[string]any, parsedBody map[string]any) (string, error) {
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

	// pree parsed body
	parsedTable := L.NewTable()
	for k, v := range parsedBody {
		switch val := v.(type) {
		case string:
			L.SetField(parsedTable, k, lua.LString(val))
		case float64:
			L.SetField(parsedTable, k, lua.LNumber(val))
		default:
			L.SetField(parsedTable, k, lua.LString(fmt.Sprintf("%v", val)))
		}
	}
	L.SetField(reqTable, "json", parsedTable)

	// Przekazanie defaultDB jako request.db
	dbTable := L.NewTable()
	for _, entry := range defaultDB {
		row := L.NewTable()
		for key, val := range entry {
			switch v := val.(type) {
			case string:
				L.SetField(row, key, lua.LString(v))
			case float64:
				L.SetField(row, key, lua.LNumber(v))
			case int:
				L.SetField(row, key, lua.LNumber(v))
			default:
				L.SetField(row, key, lua.LString(fmt.Sprintf("%v", v)))
			}
		}
		dbTable.Append(row)
	}
	L.SetField(reqTable, "db", dbTable)

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
		db := tbl.RawGetString("db") // jeśli zwrócona

		out := map[string]interface{}{
			"response": luaValueToInterface(response),
			"log":      luaValueToInterface(logs),
		}

		if db != lua.LNil {
			out["db"] = luaValueToInterface(db)
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
		// Sprawdź, czy to tablica (ciągłe indeksy liczbowe od 1)
		if isArrayLike(v) {
			// Ale uwaga: sprawdź czy to array of objects
			isArrayOfObjects := true
			v.ForEach(func(_, value lua.LValue) {
				if subtbl, ok := value.(*lua.LTable); ok {
					if !isTableObject(subtbl) {
						isArrayOfObjects = false
					}
				} else {
					isArrayOfObjects = false
				}
			})

			if isArrayOfObjects {
				var arr []interface{}
				v.ForEach(func(_, value lua.LValue) {
					arr = append(arr, luaValueToInterface(value))
				})
				return arr
			}
		}

		// W przeciwnym wypadku: obiekt (mapa)
		obj := make(map[string]interface{})
		v.ForEach(func(key, value lua.LValue) {
			obj[key.String()] = luaValueToInterface(value)
		})
		return obj

	case lua.LString:
		return v.String()
	case lua.LNumber:
		// Spróbuj zwrócić jako liczbową wartość
		if float64(int(v)) == float64(v) {
			return int(v)
		}
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LFunction, *lua.LUserData, *lua.LState:
		return nil
	default:
		return v.String()
	}
}

func isTableObject(tbl *lua.LTable) bool {
	hasKeys := false
	tbl.ForEach(func(key, _ lua.LValue) {
		if key.Type() != lua.LTNumber {
			hasKeys = true
		}
	})
	return hasKeys
}

func isArrayLike(tbl *lua.LTable) bool {
	expectedIndex := 1
	count := 0
	isObject := false

	tbl.ForEach(func(key, value lua.LValue) {
		if key.Type() != lua.LTNumber || int(key.(lua.LNumber)) != expectedIndex {
			isObject = true
		}
		expectedIndex++
		count++
	})

	return !isObject && count > 0
}
