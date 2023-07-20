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

type RpcDestActorFlags struct {
	ActorId int64
	Extra   RpcInvokeReqExtra
}

func (RpcDestActorFlags) TLName() string { return "rpcDestActorFlags" }
func (RpcDestActorFlags) TLTag() uint32  { return 0xf0a5acf7 }

func (item *RpcDestActorFlags) Reset() {
	item.ActorId = 0
	item.Extra.Reset()
}

func (item *RpcDestActorFlags) Read(w []byte) (_ []byte, err error) {
	if w, err = basictl.LongRead(w, &item.ActorId); err != nil {
		return w, err
	}
	return item.Extra.Read(w)
}

func (item *RpcDestActorFlags) Write(w []byte) (_ []byte, err error) {
	w = basictl.LongWrite(w, item.ActorId)
	return item.Extra.Write(w)
}

func (item *RpcDestActorFlags) ReadBoxed(w []byte) (_ []byte, err error) {
	if w, err = basictl.NatReadExactTag(w, 0xf0a5acf7); err != nil {
		return w, err
	}
	return item.Read(w)
}

func (item *RpcDestActorFlags) WriteBoxed(w []byte) ([]byte, error) {
	w = basictl.NatWrite(w, 0xf0a5acf7)
	return item.Write(w)
}

func (item RpcDestActorFlags) String() string {
	w, err := item.WriteJSON(nil)
	if err != nil {
		return err.Error()
	}
	return string(w)
}

func RpcDestActorFlags__ReadJSON(item *RpcDestActorFlags, j interface{}) error {
	return item.readJSON(j)
}
func (item *RpcDestActorFlags) readJSON(j interface{}) error {
	_jm, _ok := j.(map[string]interface{})
	if j != nil && !_ok {
		return ErrorInvalidJSON("rpcDestActorFlags", "expected json object")
	}
	_jActorId := _jm["actor_id"]
	delete(_jm, "actor_id")
	if err := JsonReadInt64(_jActorId, &item.ActorId); err != nil {
		return err
	}
	_jExtra := _jm["extra"]
	delete(_jm, "extra")
	for k := range _jm {
		return ErrorInvalidJSONExcessElement("rpcDestActorFlags", k)
	}
	if err := RpcInvokeReqExtra__ReadJSON(&item.Extra, _jExtra); err != nil {
		return err
	}
	return nil
}

func (item *RpcDestActorFlags) WriteJSON(w []byte) (_ []byte, err error) {
	w = append(w, '{')
	if item.ActorId != 0 {
		w = basictl.JSONAddCommaIfNeeded(w)
		w = append(w, `"actor_id":`...)
		w = basictl.JSONWriteInt64(w, item.ActorId)
	}
	w = basictl.JSONAddCommaIfNeeded(w)
	w = append(w, `"extra":`...)
	if w, err = item.Extra.WriteJSON(w); err != nil {
		return w, err
	}
	return append(w, '}'), nil
}

func (item *RpcDestActorFlags) MarshalJSON() ([]byte, error) {
	return item.WriteJSON(nil)
}

func (item *RpcDestActorFlags) UnmarshalJSON(b []byte) error {
	j, err := JsonBytesToInterface(b)
	if err != nil {
		return ErrorInvalidJSON("rpcDestActorFlags", err.Error())
	}
	if err = item.readJSON(j); err != nil {
		return ErrorInvalidJSON("rpcDestActorFlags", err.Error())
	}
	return nil
}
