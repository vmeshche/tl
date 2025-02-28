// Copyright 2022 V Kontakte LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Code generated by qtc from "qt_brackets.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package tlcodegen

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

func (tuple *TypeRWBrackets) StreamGenerateCode(qw422016 *qt422016.Writer, bytesVersion bool, directImports *DirectImports) {
	goName := addBytes(tuple.goGlobalName, bytesVersion)
	natDecl := formatNatArgsDecl(tuple.wr.NatParams)
	typeString := tuple.wr.TypeString2(bytesVersion, directImports, tuple.wr.ins, false, false)
	elementTypeString := tuple.element.t.TypeString2(bytesVersion, directImports, tuple.wr.ins, false, false)

	switch {
	case tuple.dictLike:
		keyTypeString := tuple.dictKeyField.t.TypeString2(bytesVersion, directImports, tuple.wr.ins, false, false)
		valueTypeString := tuple.dictValueField.t.TypeString2(bytesVersion, directImports, tuple.wr.ins, false, false)
		valueNatArgsDecl := formatNatArgsDecl(tuple.element.t.NatParams)
		keyFieldName := tuple.dictKeyField.goName
		valueFieldName := tuple.dictValueField.goName

		if bytesVersion {
			if tuple.wr.gen.options.GenerateRandomCode {
				qw422016.N().S(`func `)
				qw422016.N().S(goName)
				qw422016.N().S(`FillRandom(rand basictl.Rand, vec *`)
				qw422016.N().S(typeString)
				qw422016.N().S(` `)
				qw422016.N().S(natDecl)
				qw422016.N().S(`) {
    l := basictl.RandomNat(rand)
    *vec = make([]`)
				qw422016.N().S(elementTypeString)
				qw422016.N().S(`, l)
    for i := range *vec {
        `)
				qw422016.N().S(tuple.element.t.TypeRandomCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", formatNatArgs(nil, tuple.element.natArgs), false))
				qw422016.N().S(`
    }
}
`)
			}
			qw422016.N().S(`
func `)
			qw422016.N().S(goName)
			qw422016.N().S(`Read(w []byte, vec *`)
			qw422016.N().S(typeString)
			qw422016.N().S(` `)
			qw422016.N().S(natDecl)
			qw422016.N().S(`) (_ []byte, err error) {
    var l uint32
    if w, err = basictl.NatRead(w, &l); err != nil {
        return w, err
    }
    if err = basictl.CheckLengthSanity(w, l, 4); err != nil {
        return w, err
    }
    if uint32(cap(*vec)) < l {
        *vec = make([]`)
			qw422016.N().S(elementTypeString)
			qw422016.N().S(`, l)
    } else {
        *vec = (*vec)[:l]
    }
    for i := range *vec {
        `)
			qw422016.N().S(tuple.element.t.TypeReadingCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
			qw422016.N().S(`
    }
    return w, nil
}

func `)
			qw422016.N().S(goName)
			qw422016.N().S(`Write(w []byte, vec `)
			qw422016.N().S(typeString)
			qw422016.N().S(` `)
			qw422016.N().S(natDecl)
			qw422016.N().S(`) (_ []byte, err error) {
    w = basictl.NatWrite(w, uint32(len(vec)))
    for _, elem := range vec {
        `)
			qw422016.N().S(tuple.element.t.TypeWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
			qw422016.N().S(`
    }
    return w, nil
}

func `)
			qw422016.N().S(goName)
			qw422016.N().S(`ReadJSON(j interface{}, vec *`)
			qw422016.N().S(typeString)
			qw422016.N().S(` `)
			qw422016.N().S(valueNatArgsDecl)
			qw422016.N().S(`) error {
    var _map map[string]interface{}
    var _mapok bool
    if j != nil {
        _map, _mapok = j.(map[string]interface{})
        if !_mapok {
            return `)
			qw422016.N().S(tuple.wr.gen.InternalPrefix())
			qw422016.N().S(`ErrorInvalidJSON(`)
			qw422016.N().Q(typeString)
			qw422016.N().S(`, "expected json object") // TODO - better name
        }
    }
      l := len(_map)
    if cap(*vec) < l {
        *vec = make([]`)
			qw422016.N().S(elementTypeString)
			qw422016.N().S(`, l)
    } else {
        *vec = (*vec)[:l]
    }
    i := 0
    arr := *vec
`)
			if tuple.dictKeyString {
				qw422016.N().S(`for key, _jvalue := range _map {
        arr[i].`)
				qw422016.N().S(tuple.dictKeyField.goName)
				qw422016.N().S(` = append(arr[i].`)
				qw422016.N().S(tuple.dictKeyField.goName)
				qw422016.N().S(`[:0], key...)
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_jvalue", "arr[i]."+valueFieldName, formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
        i++
    }
    return nil
}

func `)
				qw422016.N().S(goName)
				qw422016.N().S(`WriteJSON(w []byte, vec `)
				qw422016.N().S(typeString)
				qw422016.N().S(` `)
				qw422016.N().S(valueNatArgsDecl)
				qw422016.N().S(`) (_ []byte, err error) {
    w = append(w, '{')
    for _, elem := range vec {
        w = basictl.JSONAddCommaIfNeeded(w)
        w = basictl.JSONWriteStringBytes(w, elem.`)
				qw422016.N().S(tuple.dictKeyField.goName)
				qw422016.N().S(`)
        w = append(w, ':')
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem."+valueFieldName, formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
    }
    return append(w, '}'), nil
}
`)
			} else {
				qw422016.N().S(`for _jkey, _jvalue := range _map {
        `)
				qw422016.N().S(tuple.dictKeyField.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_jkey", "arr[i]."+keyFieldName, formatNatArgs(nil, tuple.dictKeyField.natArgs), false))
				qw422016.N().S(`
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_jvalue", "arr[i]."+valueFieldName, formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
        i++
    }
    return nil
}

func `)
				qw422016.N().S(goName)
				qw422016.N().S(`WriteJSON(w []byte, vec `)
				qw422016.N().S(typeString)
				qw422016.N().S(` `)
				qw422016.N().S(valueNatArgsDecl)
				qw422016.N().S(`) (_ []byte, err error) {
    w = append(w, '{')
    for _, elem := range vec {
        key := elem.`)
				qw422016.N().S(keyFieldName)
				qw422016.N().S(`
        w = basictl.JSONAddCommaIfNeeded(w)
        w = append(w, `)
				qw422016.N().S("`")
				qw422016.N().S(`"`)
				qw422016.N().S("`")
				qw422016.N().S(`...)
        `)
				qw422016.N().S(tuple.dictKeyField.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "key", formatNatArgs(nil, tuple.dictKeyField.natArgs), false))
				qw422016.N().S(`
        w = append(w, `)
				qw422016.N().S("`")
				qw422016.N().S(`":`)
				qw422016.N().S("`")
				qw422016.N().S(`...)
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem."+valueFieldName, formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
    }
    return append(w, '}'), nil
}
`)
			}
		} else {
			qw422016.N().S(`func `)
			qw422016.N().S(goName)
			qw422016.N().S(`Reset(m map[`)
			qw422016.N().S(keyTypeString)
			qw422016.N().S(`]`)
			qw422016.N().S(valueTypeString)
			qw422016.N().S(`) {
    for k := range m {
        delete(m, k)
    }
}

`)
			if tuple.wr.gen.options.GenerateRandomCode {
				qw422016.N().S(`func `)
				qw422016.N().S(goName)
				qw422016.N().S(`FillRandom(rand basictl.Rand, m *map[`)
				qw422016.N().S(keyTypeString)
				qw422016.N().S(`]`)
				qw422016.N().S(valueTypeString)
				qw422016.N().S(` `)
				qw422016.N().S(natDecl)
				qw422016.N().S(`) {
    l := basictl.RandomNat(rand)
    *m = make(map[`)
				qw422016.N().S(keyTypeString)
				qw422016.N().S(`]`)
				qw422016.N().S(valueTypeString)
				qw422016.N().S(`, l)
    for i := 0; i < int(l); i++ {
        var elem `)
				qw422016.N().S(elementTypeString)
				qw422016.N().S(`
        `)
				qw422016.N().S(tuple.element.t.TypeRandomCode(bytesVersion, directImports, tuple.wr.ins, "elem", formatNatArgs(nil, tuple.element.natArgs), false))
				qw422016.N().S(`
        (*m)[elem.`)
				qw422016.N().S(keyFieldName)
				qw422016.N().S(`] = elem.`)
				qw422016.N().S(valueFieldName)
				qw422016.N().S(`
    }
}
`)
			}
			qw422016.N().S(`func `)
			qw422016.N().S(goName)
			qw422016.N().S(`Read(w []byte, m *map[`)
			qw422016.N().S(keyTypeString)
			qw422016.N().S(`]`)
			qw422016.N().S(valueTypeString)
			qw422016.N().S(` `)
			qw422016.N().S(natDecl)
			qw422016.N().S(`) (_ []byte, err error) {
    var l uint32
    if w, err = basictl.NatRead(w, &l); err != nil { // No sanity check required for map
        return w, err
    }
    var data map[`)
			qw422016.N().S(keyTypeString)
			qw422016.N().S(`]`)
			qw422016.N().S(valueTypeString)
			qw422016.N().S(`
    if *m == nil {
        if l == 0 {
            return w, nil
        }
        data = make(map[`)
			qw422016.N().S(keyTypeString)
			qw422016.N().S(`]`)
			qw422016.N().S(valueTypeString)
			qw422016.N().S(`, l)
        *m = data
    } else {
        data = *m
        for k := range data {
            delete(data, k)
        }
    }
    for i := 0; i < int(l); i++ {
        var elem `)
			qw422016.N().S(elementTypeString)
			qw422016.N().S(`
        `)
			qw422016.N().S(tuple.element.t.TypeReadingCode(bytesVersion, directImports, tuple.wr.ins, "elem", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
			qw422016.N().S(`
         data[elem.`)
			qw422016.N().S(keyFieldName)
			qw422016.N().S(`] = elem.`)
			qw422016.N().S(valueFieldName)
			qw422016.N().S(`
    }
    return w, nil
}

func `)
			qw422016.N().S(goName)
			qw422016.N().S(`Write(w []byte, m map[`)
			qw422016.N().S(keyTypeString)
			qw422016.N().S(`]`)
			qw422016.N().S(valueTypeString)
			qw422016.N().S(` `)
			qw422016.N().S(natDecl)
			qw422016.N().S(`) (_ []byte, err error) {
    w = basictl.NatWrite(w, uint32(len(m)))
    for key, val := range m {
        elem := `)
			qw422016.N().S(elementTypeString)
			qw422016.N().S(`{`)
			qw422016.N().S(keyFieldName)
			qw422016.N().S(`:key, `)
			qw422016.N().S(valueFieldName)
			qw422016.N().S(`:val}
        `)
			qw422016.N().S(tuple.element.t.TypeWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
			qw422016.N().S(`
    }
    return w, nil
}

func `)
			qw422016.N().S(goName)
			qw422016.N().S(`ReadJSON(j interface{}, m *`)
			qw422016.N().S(typeString)
			qw422016.N().S(` `)
			qw422016.N().S(valueNatArgsDecl)
			qw422016.N().S(`) error {
    var _map map[string]interface{}
    var _mapok bool
    if j != nil {
        _map, _mapok = j.(map[string]interface{})
        if !_mapok {
            return `)
			qw422016.N().S(tuple.wr.gen.InternalPrefix())
			qw422016.N().S(`ErrorInvalidJSON(`)
			qw422016.N().Q(typeString)
			qw422016.N().S(`, "expected json object") // TODO - better name
        }
    }
    l := len(_map)
    var data map[`)
			qw422016.N().S(keyTypeString)
			qw422016.N().S(`]`)
			qw422016.N().S(valueTypeString)
			qw422016.N().S(`
    if *m == nil {
        if l == 0 {
            return nil
        }
        data = make(map[`)
			qw422016.N().S(keyTypeString)
			qw422016.N().S(`]`)
			qw422016.N().S(valueTypeString)
			qw422016.N().S(`, l)
        *m = data
    } else {
        data = *m
        for k := range data {
            delete(data, k)
        }
    }
`)
			if tuple.dictKeyString {
				qw422016.N().S(`for key, _jvalue := range _map {
        var value `)
				qw422016.N().S(valueTypeString)
				qw422016.N().S(`
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_jvalue", "value", formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
        data[key] = value
    }
    return nil
}

func `)
				qw422016.N().S(goName)
				qw422016.N().S(`WriteJSON(w []byte, m `)
				qw422016.N().S(typeString)
				qw422016.N().S(` `)
				qw422016.N().S(valueNatArgsDecl)
				qw422016.N().S(`) (_ []byte, err error) {
    keys := make([]`)
				qw422016.N().S(keyTypeString)
				qw422016.N().S(`, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
`)
				directImports.importSort = true

				qw422016.N().S(`    sort.Strings(keys)
    w = append(w, '{')
    for _, key := range keys {
        value := m[key]
        w = basictl.JSONAddCommaIfNeeded(w)
        w = basictl.JSONWriteString(w, key) // StringKey
        w = append(w, ':')
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "value", formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
    }
    return append(w, '}'), nil
}
`)
			} else {
				qw422016.N().S(`for _jkey, _jvalue := range _map {
        var key `)
				qw422016.N().S(keyTypeString)
				qw422016.N().S(`
        `)
				qw422016.N().S(tuple.dictKeyField.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_jkey", "key", formatNatArgs(nil, tuple.dictKeyField.natArgs), false))
				qw422016.N().S(`
        var value `)
				qw422016.N().S(valueTypeString)
				qw422016.N().S(`
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_jvalue", "value", formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
        data[key] = value
    }
    return nil
}

func `)
				qw422016.N().S(goName)
				qw422016.N().S(`WriteJSON(w []byte, m `)
				qw422016.N().S(typeString)
				qw422016.N().S(` `)
				qw422016.N().S(natDecl)
				qw422016.N().S(`) (_ []byte, err error) {
    keys := make([]`)
				qw422016.N().S(keyTypeString)
				qw422016.N().S(`, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
`)
				directImports.importSort = true

				qw422016.N().S(`    sort.Slice(keys, func(i, j int) bool {
        return keys[i] < keys[j]
    })
    w = append(w, '{')
    for _, key := range keys {
        value := m[key]
        w = basictl.JSONAddCommaIfNeeded(w)
        w = append(w, `)
				qw422016.N().S("`")
				qw422016.N().S(`"`)
				qw422016.N().S("`")
				qw422016.N().S(`...)
        `)
				qw422016.N().S(tuple.dictKeyField.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "key", formatNatArgs(nil, tuple.dictKeyField.natArgs), false))
				qw422016.N().S(`
        w = append(w, `)
				qw422016.N().S("`")
				qw422016.N().S(`":`)
				qw422016.N().S("`")
				qw422016.N().S(`...)
        `)
				qw422016.N().S(tuple.dictValueField.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "value", formatNatArgs(nil, tuple.dictValueField.natArgs), false))
				qw422016.N().S(`
    }
    return append(w, '}'), nil
}
`)
			}
		}
	case tuple.vectorLike:
		if tuple.wr.gen.options.GenerateRandomCode {
			qw422016.N().S(`func `)
			qw422016.N().S(goName)
			qw422016.N().S(`FillRandom(rand basictl.Rand, vec *`)
			qw422016.N().S(typeString)
			qw422016.N().S(` `)
			qw422016.N().S(natDecl)
			qw422016.N().S(`) {
    l := basictl.RandomNat(rand)
    *vec = make([]`)
			qw422016.N().S(elementTypeString)
			qw422016.N().S(`, l)
    for i := range *vec {
        `)
			qw422016.N().S(tuple.element.t.TypeRandomCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", formatNatArgs(nil, tuple.element.natArgs), false))
			qw422016.N().S(`
    }
}
`)
		}
		qw422016.N().S(`func `)
		qw422016.N().S(goName)
		qw422016.N().S(`Read(w []byte, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    var l uint32
    if w, err = basictl.NatRead(w, &l); err != nil {
        return w, err
    }
    if err = basictl.CheckLengthSanity(w, l, 4); err != nil {
        return w, err
    }
    if uint32(cap(*vec)) < l {
        *vec = make([]`)
		qw422016.N().S(elementTypeString)
		qw422016.N().S(`, l)
    } else {
        *vec = (*vec)[:l]
    }
    for i := range *vec {
        `)
		qw422016.N().S(tuple.element.t.TypeReadingCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
		qw422016.N().S(`
    }
    return w, nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`Write(w []byte, vec `)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    w = basictl.NatWrite(w, uint32(len(vec)))
    for _, elem := range vec {
        `)
		qw422016.N().S(tuple.element.t.TypeWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
		qw422016.N().S(`
    }
    return w, nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`ReadJSON(j interface{}, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) error {
    l, _arr, err := `)
		qw422016.N().S(tuple.wr.gen.InternalPrefix())
		qw422016.N().S(`JsonReadArray(`)
		qw422016.N().Q(typeString)
		qw422016.N().S(`, j)
    if err != nil {
        return err
    }
    if cap(*vec) < l {
        *vec = make([]`)
		qw422016.N().S(elementTypeString)
		qw422016.N().S(`, l)
    } else {
        *vec = (*vec)[:l]
    }
    for i := range *vec {
        `)
		qw422016.N().S(tuple.element.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_arr[i]", "(*vec)[i]", formatNatArgs(nil, tuple.element.natArgs), false))
		qw422016.N().S(`
    }
    return nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`WriteJSON(w []byte, vec `)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    w = append(w, '[')
    for _, elem := range vec {
        w = basictl.JSONAddCommaIfNeeded(w)
        `)
		qw422016.N().S(tuple.element.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", formatNatArgs(nil, tuple.element.natArgs), false))
		qw422016.N().S(`
    }
    return append(w, ']'), nil
}

`)
	case tuple.dynamicSize:
		if tuple.wr.gen.options.GenerateRandomCode {
			qw422016.N().S(`func `)
			qw422016.N().S(goName)
			qw422016.N().S(`FillRandom(rand basictl.Rand, vec *`)
			qw422016.N().S(typeString)
			qw422016.N().S(` `)
			qw422016.N().S(natDecl)
			qw422016.N().S(`) {
    *vec = make([]`)
			qw422016.N().S(elementTypeString)
			qw422016.N().S(`, nat_n)
    for i := range *vec {
        `)
			qw422016.N().S(tuple.element.t.TypeRandomCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", formatNatArgs(nil, tuple.element.natArgs), false))
			qw422016.N().S(`
    }
}
`)
		}
		qw422016.N().S(`
func `)
		qw422016.N().S(goName)
		qw422016.N().S(`Read(w []byte, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    if err = basictl.CheckLengthSanity(w, nat_n, 4); err != nil {
        return w, err
    }
    if uint32(cap(*vec)) < nat_n {
        *vec = make([]`)
		qw422016.N().S(elementTypeString)
		qw422016.N().S(`, nat_n)
    } else {
        *vec = (*vec)[:nat_n]
    }
    for i := range *vec {
        `)
		qw422016.N().S(tuple.element.t.TypeReadingCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
		qw422016.N().S(`
    }
    return w, nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`Write(w []byte, vec `)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    if uint32(len(vec)) != nat_n {
        return w, `)
		qw422016.N().S(tuple.wr.gen.InternalPrefix())
		qw422016.N().S(`ErrorWrongSequenceLength(`)
		qw422016.N().Q(typeString)
		qw422016.N().S(`, len(vec), nat_n)
    }
    for _, elem := range vec {
        `)
		qw422016.N().S(tuple.element.t.TypeWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
		qw422016.N().S(`
    }
    return w, nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`ReadJSON(j interface{}, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) error {
    _, _arr, err := `)
		qw422016.N().S(tuple.wr.gen.InternalPrefix())
		qw422016.N().S(`JsonReadArrayFixedSize(`)
		qw422016.N().Q(typeString)
		qw422016.N().S(`, j, nat_n)
    if err != nil {
        return err
    }
    if uint32(cap(*vec)) < nat_n {
        *vec = make([]`)
		qw422016.N().S(elementTypeString)
		qw422016.N().S(`, nat_n)
    } else {
        *vec = (*vec)[:nat_n]
    }
    for i := range *vec {
        `)
		qw422016.N().S(tuple.element.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_arr[i]", "(*vec)[i]", formatNatArgs(nil, tuple.element.natArgs), false))
		qw422016.N().S(`
    }
    return nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`WriteJSON(w []byte, vec `)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    if uint32(len(vec)) != nat_n {
        return w, `)
		qw422016.N().S(tuple.wr.gen.InternalPrefix())
		qw422016.N().S(`ErrorWrongSequenceLength(`)
		qw422016.N().Q(typeString)
		qw422016.N().S(`, len(vec), nat_n)
    }
    w = append(w, '[')
    for _, elem := range vec {
        w = basictl.JSONAddCommaIfNeeded(w)
        `)
		qw422016.N().S(tuple.element.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", formatNatArgs(nil, tuple.element.natArgs), false))
		qw422016.N().S(`
    }
    return append(w, ']'), nil
}

`)
	default:
		qw422016.N().S(`func `)
		qw422016.N().S(goName)
		qw422016.N().S(`Reset(vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(`) {
    for i := range *vec {
            `)
		qw422016.N().S(tuple.element.t.TypeResettingCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", false))
		qw422016.N().S(`
    }
}

`)
		if tuple.wr.gen.options.GenerateRandomCode {
			qw422016.N().S(`func `)
			qw422016.N().S(goName)
			qw422016.N().S(`FillRandom(rand basictl.Rand, vec *`)
			qw422016.N().S(typeString)
			qw422016.N().S(` `)
			qw422016.N().S(natDecl)
			qw422016.N().S(`) {
    for i := range *vec {
        `)
			qw422016.N().S(tuple.element.t.TypeRandomCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", formatNatArgs(nil, tuple.element.natArgs), false))
			qw422016.N().S(`
    }
}
`)
		}
		qw422016.N().S(`
func `)
		qw422016.N().S(goName)
		qw422016.N().S(`Read(w []byte, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    for i := range *vec {
        `)
		qw422016.N().S(tuple.element.t.TypeReadingCode(bytesVersion, directImports, tuple.wr.ins, "(*vec)[i]", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
		qw422016.N().S(`
    }
    return w, nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`Write(w []byte, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    for _, elem := range *vec {
        `)
		qw422016.N().S(tuple.element.t.TypeWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", tuple.element.Bare(), formatNatArgs(nil, tuple.element.natArgs), false, false))
		qw422016.N().S(`
    }
    return w, nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`ReadJSON(j interface{}, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) error {
    _, _arr, err := `)
		qw422016.N().S(tuple.wr.gen.InternalPrefix())
		qw422016.N().S(`JsonReadArrayFixedSize(`)
		qw422016.N().Q(typeString)
		qw422016.N().S(`, j, `)
		qw422016.E().V(tuple.size)
		qw422016.N().S(`)
    if err != nil {
        return err
    }
    for i := range *vec {
        `)
		qw422016.N().S(tuple.element.t.TypeJSONReadingCode(bytesVersion, directImports, tuple.wr.ins, "_arr[i]", "(*vec)[i]", formatNatArgs(nil, tuple.element.natArgs), false))
		qw422016.N().S(`
    }
    return nil
}

func `)
		qw422016.N().S(goName)
		qw422016.N().S(`WriteJSON(w []byte, vec *`)
		qw422016.N().S(typeString)
		qw422016.N().S(` `)
		qw422016.N().S(natDecl)
		qw422016.N().S(`) (_ []byte, err error) {
    w = append(w, '[')
    for _, elem := range *vec {
        w = basictl.JSONAddCommaIfNeeded(w)
        `)
		qw422016.N().S(tuple.element.t.TypeJSONWritingCode(bytesVersion, directImports, tuple.wr.ins, "elem", formatNatArgs(nil, tuple.element.natArgs), false))
		qw422016.N().S(`
    }
    return append(w, ']'), nil
}
`)
	}
}

func (tuple *TypeRWBrackets) WriteGenerateCode(qq422016 qtio422016.Writer, bytesVersion bool, directImports *DirectImports) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	tuple.StreamGenerateCode(qw422016, bytesVersion, directImports)
	qt422016.ReleaseWriter(qw422016)
}

func (tuple *TypeRWBrackets) GenerateCode(bytesVersion bool, directImports *DirectImports) string {
	qb422016 := qt422016.AcquireByteBuffer()
	tuple.WriteGenerateCode(qb422016, bytesVersion, directImports)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
