// Copyright 2022 V Kontakte LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Code generated by vktl/cmd/tlgen2; DO NOT EDIT.
package internal

import (
	"github.com/vkcom/tl/internal/basictl"
)

var _ = basictl.NatWrite

type RpcReqResultHeader struct {
	QueryId int64
}

func (RpcReqResultHeader) TLName() string { return "rpcReqResultHeader" }
func (RpcReqResultHeader) TLTag() uint32  { return 0x63aeda4e }

func (item *RpcReqResultHeader) Reset() {
	item.QueryId = 0
}

func (item *RpcReqResultHeader) Read(w []byte) (_ []byte, err error) {
	return basictl.LongRead(w, &item.QueryId)
}

func (item *RpcReqResultHeader) Write(w []byte) (_ []byte, err error) {
	return basictl.LongWrite(w, item.QueryId), nil
}

func (item *RpcReqResultHeader) ReadBoxed(w []byte) (_ []byte, err error) {
	if w, err = basictl.NatReadExactTag(w, 0x63aeda4e); err != nil {
		return w, err
	}
	return item.Read(w)
}

func (item *RpcReqResultHeader) WriteBoxed(w []byte) ([]byte, error) {
	w = basictl.NatWrite(w, 0x63aeda4e)
	return item.Write(w)
}

func (item RpcReqResultHeader) String() string {
	w, err := item.WriteJSON(nil)
	if err != nil {
		return err.Error()
	}
	return string(w)
}

func RpcReqResultHeader__ReadJSON(item *RpcReqResultHeader, j interface{}) error {
	return item.readJSON(j)
}
func (item *RpcReqResultHeader) readJSON(j interface{}) error {
	_jm, _ok := j.(map[string]interface{})
	if j != nil && !_ok {
		return ErrorInvalidJSON("rpcReqResultHeader", "expected json object")
	}
	_jQueryId := _jm["query_id"]
	delete(_jm, "query_id")
	if err := JsonReadInt64(_jQueryId, &item.QueryId); err != nil {
		return err
	}
	for k := range _jm {
		return ErrorInvalidJSONExcessElement("rpcReqResultHeader", k)
	}
	return nil
}

func (item *RpcReqResultHeader) WriteJSON(w []byte) (_ []byte, err error) {
	w = append(w, '{')
	if item.QueryId != 0 {
		w = basictl.JSONAddCommaIfNeeded(w)
		w = append(w, `"query_id":`...)
		w = basictl.JSONWriteInt64(w, item.QueryId)
	}
	return append(w, '}'), nil
}

func (item *RpcReqResultHeader) MarshalJSON() ([]byte, error) {
	return item.WriteJSON(nil)
}

func (item *RpcReqResultHeader) UnmarshalJSON(b []byte) error {
	j, err := JsonBytesToInterface(b)
	if err != nil {
		return ErrorInvalidJSON("rpcReqResultHeader", err.Error())
	}
	if err = item.readJSON(j); err != nil {
		return ErrorInvalidJSON("rpcReqResultHeader", err.Error())
	}
	return nil
}
