// Copyright 2022 V Kontakte LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Code generated by vktl/cmd/tlgen2; DO NOT EDIT.
package gen_tlo

type TlsCombinatorRight struct {
	Value TlsTypeExprUnion
}

func (TlsCombinatorRight) TLName() string { return "tls.combinatorRight" }
func (TlsCombinatorRight) TLTag() uint32  { return 0x2c064372 }

func (item *TlsCombinatorRight) Reset() {
	item.Value.Reset()
}

func (item *TlsCombinatorRight) Read(w []byte) (_ []byte, err error) {
	return item.Value.ReadBoxed(w)
}

func (item *TlsCombinatorRight) Write(w []byte) (_ []byte, err error) {
	return item.Value.WriteBoxed(w)
}

func (item *TlsCombinatorRight) ReadBoxed(w []byte) (_ []byte, err error) {
	if w, err = NatReadExactTag(w, 0x2c064372); err != nil {
		return w, err
	}
	return item.Read(w)
}

func (item *TlsCombinatorRight) WriteBoxed(w []byte) ([]byte, error) {
	w = NatWrite(w, 0x2c064372)
	return item.Write(w)
}

func (item TlsCombinatorRight) String() string {
	w, err := item.WriteJSON(nil)
	if err != nil {
		return err.Error()
	}
	return string(w)
}

func TlsCombinatorRight__ReadJSON(item *TlsCombinatorRight, j interface{}) error {
	return item.readJSON(j)
}
func (item *TlsCombinatorRight) readJSON(j interface{}) error {
	_jm, _ok := j.(map[string]interface{})
	if j != nil && !_ok {
		return ErrorInvalidJSON("tls.combinatorRight", "expected json object")
	}
	_jValue := _jm["value"]
	delete(_jm, "value")
	for k := range _jm {
		return ErrorInvalidJSONExcessElement("tls.combinatorRight", k)
	}
	if err := TlsTypeExprUnion__ReadJSON(&item.Value, _jValue); err != nil {
		return err
	}
	return nil
}

func (item *TlsCombinatorRight) WriteJSON(w []byte) (_ []byte, err error) {
	w = append(w, '{')
	w = JSONAddCommaIfNeeded(w)
	w = append(w, `"value":`...)
	if w, err = item.Value.WriteJSON(w); err != nil {
		return w, err
	}
	return append(w, '}'), nil
}

func (item *TlsCombinatorRight) MarshalJSON() ([]byte, error) {
	return item.WriteJSON(nil)
}

func (item *TlsCombinatorRight) UnmarshalJSON(b []byte) error {
	j, err := JsonBytesToInterface(b)
	if err != nil {
		return ErrorInvalidJSON("tls.combinatorRight", err.Error())
	}
	if err = item.readJSON(j); err != nil {
		return ErrorInvalidJSON("tls.combinatorRight", err.Error())
	}
	return nil
}
