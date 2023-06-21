// Copyright 2022 V Kontakte LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tlcodegen

const binaryStringObjectKey = "base64"

const basicTLCodeHeader = `%s
package %s

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
	"unicode/utf8"
)

`

const basicTLCodeBody = `

const (
	tinyStringLen    = 253
	bigStringMarker  = 0xfe
	hugeStringMarker = 0xff
	bigStringLen     = (1 << 24) - 1
	hugeStringLen    = (1 << 56) - 1
)

var errBadPadding = fmt.Errorf("non-canonical non-zero string padding")

func CheckLengthSanity(r []byte, natParam uint32, minObjectSize uint32) error {
	if uint64(len(r)) < uint64(natParam)*uint64(minObjectSize) { // Must wrap io.ErrUnexpectedEOF
		return fmt.Errorf("invalid length: %d for remaining reader length: %d and min object size %d: %w", natParam, len(r), minObjectSize, io.ErrUnexpectedEOF)
	}
	return nil
}

func ReadBool(r []byte, v *bool, falseTag uint32, trueTag uint32) ([]byte, error) {
	tag, r, err := NatReadTag(r)
	if err != nil {
		return r, err
	}
	switch tag {
	case falseTag:
		*v = false
	case trueTag:
		*v = true
	default:
		return r, fmt.Errorf("invalid bool tag: 0x%x", tag)
	}
	return r, nil
}

func NatRead(r []byte, dst *uint32) ([]byte, error) {
	if len(r) < 4 {
		return r, io.ErrUnexpectedEOF
	}
	*dst = binary.LittleEndian.Uint32(r)
	return r[4:], nil
}

func NatWrite(w []byte, v uint32) []byte {
	return append(w, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func nat64Write(w []byte, v uint64) []byte {
	return append(w, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
}

func IntRead(r []byte, dst *int32) ([]byte, error) {
	if len(r) < 4 {
		return r, io.ErrUnexpectedEOF
	}
	*dst = int32(binary.LittleEndian.Uint32(r))
	return r[4:], nil
}

func IntWrite(w []byte, v int32) []byte {
	return append(w, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func LongRead(r []byte, dst *int64) ([]byte, error) {
	if len(r) < 8 {
		return r, io.ErrUnexpectedEOF
	}
	*dst = int64(binary.LittleEndian.Uint64(r))
	return r[8:], nil
}

func LongWrite(w []byte, v int64) []byte {
	return nat64Write(w, uint64(v))
}

func FloatRead(r []byte, dst *float32) ([]byte, error) {
	if len(r) < 4 {
		return r, io.ErrUnexpectedEOF
	}
	*dst = math.Float32frombits(binary.LittleEndian.Uint32(r))
	return r[4:], nil
}

func FloatWrite(w []byte, v float32) []byte {
	return NatWrite(w, math.Float32bits(v))
}

func DoubleRead(r []byte, dst *float64) ([]byte, error) {
	if len(r) < 8 {
		return r, io.ErrUnexpectedEOF
	}
	*dst = math.Float64frombits(binary.LittleEndian.Uint64(r))
	return r[8:], nil
}

func DoubleWrite(w []byte, v float64) []byte {
	return nat64Write(w, math.Float64bits(v))
}

func StringRead(r []byte, dst *string) ([]byte, error) {
	var b []byte
	r, err := StringReadBytes(r, &b)
	if err != nil {
		return r, err
	}
	*dst = string(b)
	return r, nil
}

func StringWrite(w []byte, v string) ([]byte, error) {
	return writeString(w, v), nil
}

func StringWriteTruncated(w []byte, v string) []byte {
	return writeString(w, v)
}

func StringReadBytes(r []byte, dst *[]byte) ([]byte, error) {
	if len(r) == 0 {
		return r, io.ErrUnexpectedEOF
	}
	b0 := r[0]

	var l int
	var p int
	switch {
	case b0 <= tinyStringLen:
		l = int(b0)
		r = r[1:]
		p = l + 1
	case b0 == bigStringMarker:
		if len(r) < 4 {
			return r, io.ErrUnexpectedEOF
		}
		l = (int(r[3]) << 16) + (int(r[2]) << 8) + (int(r[1]) << 0)
		r = r[4:]
		p = l // +4
		if l <= tinyStringLen {
			return r, fmt.Errorf("non-canonical (big) string format for length: %d", l)
		}
	default: // hugeStringMarker
		if len(r) < 8 {
			return r, io.ErrUnexpectedEOF
		}
		l64 := (int64(r[7]) << 48) + (int64(r[6]) << 40) + (int64(r[5]) << 32) + (int64(r[4]) << 24) + (int64(r[3]) << 16) + (int64(r[2]) << 8) + (int64(r[1]) << 0)
		if l64 > math.MaxInt {
			return r, fmt.Errorf("string length cannot be represented on 32-bit platform: %d", l64)
		}
		l = int(l64)
		r = r[8:]
		p = l // +8
		if l <= bigStringLen {
			return r, fmt.Errorf("non-canonical (huge) string format for length: %d", l)
		}
	}

	if l > 0 {
		if len(r) < l {
			return r, io.ErrUnexpectedEOF
		}
		// Allocate only after we know there is enough bytes in reader
		if cap(*dst) < l {
			*dst = make([]byte, l)
		} else {
			*dst = (*dst)[:l]
		}
		copy(*dst, r)
	} else {
		*dst = (*dst)[:0]
	}
	padding := paddingLen(p)
	if len(r) < l+padding {
		return r, io.ErrUnexpectedEOF
	}
	for i := 0; i < padding; i++ {
		if r[l+i] != 0 {
			return r, errBadPadding
		}
	}
	return r[l+padding:], nil
}

func StringWriteBytes(w []byte, v []byte) ([]byte, error) {
	return writeStringBytes(w, v), nil
}

func StringWriteBytesTruncated(w []byte, v []byte) []byte {
	return writeStringBytes(w, v)
}

func NatPeekTag(r []byte) (uint32, error) {
	if len(r) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.LittleEndian.Uint32(r), nil
}

func NatReadTag(r []byte) (uint32, []byte, error) {
	if len(r) < 4 {
		return 0, r, io.ErrUnexpectedEOF
	}
	return binary.LittleEndian.Uint32(r), r[4:], nil
}

func NatReadExactTag(r []byte, tag uint32) ([]byte, error) {
	if len(r) < 4 {
		return r, io.ErrUnexpectedEOF
	}
	if t := binary.LittleEndian.Uint32(r); t != tag {
		return r, fmt.Errorf("read tag #%08x instead of #%08x", t, tag)
	}
	return r[4:], nil
}

func paddingLen(l int) int {
	return int(-uint(l) % 4)
}

func writeString(w []byte, v string) []byte {
	l := len(v)
	var p int
	switch {
	case l <= tinyStringLen:
		w = append(w, byte(l))
		p = l + 1
	case l <= bigStringLen:
		w = append(w, bigStringMarker, byte(l), byte(l>>8), byte(l>>16))
		p = l // +4
	default:
		if l > hugeStringLen { // for correctness only, we do not expect strings so huge
			l = hugeStringLen
			v = v[:l]
		}
		w = append(w, hugeStringMarker, byte(l), byte(l>>8), byte(l>>16), byte(l>>24), byte(l>>32), byte(l>>40), byte(l>>48))
		p = l // +8
	}
	w = append(w, v...)

	switch uint(p) % 4 {
	case 1:
		w = append(w, 0, 0, 0)
	case 2:
		w = append(w, 0, 0)
	case 3:
		w = append(w, 0)
	}
	return w
}

func writeStringBytes(w []byte, v []byte) []byte {
	l := len(v)
	var p int
	switch {
	case l <= tinyStringLen:
		w = append(w, byte(l))
		p = l + 1
	case l <= bigStringLen:
		w = append(w, bigStringMarker, byte(l), byte(l>>8), byte(l>>16))
		p = l // +4
	default:
		if l > hugeStringLen { // for correctness only, we do not expect strings so huge
			l = hugeStringLen
			v = v[:l]
		}
		w = append(w, hugeStringMarker, byte(l), byte(l>>8), byte(l>>16), byte(l>>24), byte(l>>32), byte(l>>40), byte(l>>48))
		p = l // +8
	}
	w = append(w, v...)

	switch uint(p) % 4 {
	case 1:
		w = append(w, 0, 0, 0)
	case 2:
		w = append(w, 0, 0)
	case 3:
		w = append(w, 0)
	}
	return w
}

func JSONWriteBool(w []byte, v bool) []byte {
	return append(w, strconv.FormatBool(v)...)
}

func JSONWriteUint32(w []byte, v uint32) []byte {
	return strconv.AppendUint(w, uint64(v), 10)
}

func JSONWriteInt32(w []byte, v int32) []byte {
	return strconv.AppendInt(w, int64(v), 10)
}

func JSONWriteUint64(w []byte, v uint64) []byte {
	return strconv.AppendUint(w, v, 10)
}

func JSONWriteInt64(w []byte, v int64) []byte {
	return strconv.AppendInt(w, v, 10)
}

func JSONWriteFloat32(w []byte, v float32) []byte {
	return strconv.AppendFloat(w, float64(v), 'f', -1, 32)
}

func JSONWriteFloat64(w []byte, v float64) []byte {
	return strconv.AppendFloat(w, v, 'f', -1, 64)
}

func JSONAddCommaIfNeeded(w []byte) []byte { // Never called on empty buffer, so panic is ok
	lastChar := w[len(w)-1]
	if lastChar != '{' && lastChar != '[' {
		return append(w, ',')
	}
	return w
}

// JSON string escaping is below, keep in sync with go/src/encoding/json/encode.go

const hex = "0123456789abcdef"

var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'\x60':   true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

const (
	binaryJSONStringStart = "{\"` + binaryStringObjectKey + `\":\""
	binaryJSONStringEnd   = "\"}"
)

func alloc(buf []byte, size int) []byte {
	if cap(buf) >= len(buf) + size {
		return buf[:len(buf)+size]
	}
	return append(buf, make([]byte, size)...)
}

// NOTE: keep in sync with stringBytes below.

func JSONWriteString(w []byte, s string) []byte {
	if !utf8.ValidString(s) {
		w = append(w, binaryJSONStringStart...)
		beforeAllocation := len(w)
		w = alloc(w, base64.StdEncoding.EncodedLen(len(s)))
		base64.StdEncoding.Encode(w[beforeAllocation:], []byte(s))
		return append(w, binaryJSONStringEnd...)
	}
	w = append(w, '"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				w = append(w, s[start:i]...)
			}
			w = append(w, '\\')
			switch b {
			case '\\', '"':
				w = append(w, b)
			case '\n':
				w = append(w, 'n')
			case '\r':
				w = append(w, 'r')
			case '\t':
				w = append(w, 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w = append(w, "u00"...)
				w = append(w, hex[b>>4])
				w = append(w, hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w = append(w, s[start:i]...)
			}
			w = append(w, "\\ufffd"...)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w = append(w, s[start:i]...)
			}
			w = append(w, "\\u202"...)
			w = append(w, hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		w = append(w, s[start:]...)
	}
	return append(w, '"')
}

// NOTE: keep in sync with string above.
func JSONWriteStringBytes(w []byte, s []byte) []byte {
	if !utf8.Valid(s) {
		w = append(w, binaryJSONStringStart...)
		beforeAllocation := len(w)
		w = alloc(w, base64.StdEncoding.EncodedLen(len(s)))
		base64.StdEncoding.Encode(w[beforeAllocation:], s)
		return append(w, binaryJSONStringEnd...)
	}
	w = append(w, '"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				w = append(w, s[start:i]...)
			}
			w = append(w, '\\')
			switch b {
			case '\\', '"':
				w = append(w, b)
			case '\n':
				w = append(w, 'n')
			case '\r':
				w = append(w, 'r')
			case '\t':
				w = append(w, 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w = append(w, "u00"...)
				w = append(w, hex[b>>4])
				w = append(w, hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w = append(w, s[start:i]...)
			}
			w = append(w, "\\ufffd"...)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w = append(w, s[start:i]...)
			}
			w = append(w, "\\u202"...)
			w = append(w, hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		w = append(w, s[start:]...)
	}
	return append(w, '"')
}

type Rand interface {
	Uint32() uint32
	Int31() int32
	Int63() int64
	NormFloat64() float64
}

const RandomNatConstraint = 10

func RandomNat(rand Rand) uint32 {
	return rand.Uint32() % RandomNatConstraint
}

func RandomInt(rand Rand) int32 {
	return rand.Int31()
}

func RandomLong(rand Rand) int64 {
	return rand.Int63()
}

func RandomFloat(rand Rand) float32 {
	return float32(rand.NormFloat64())
}

func RandomDouble(rand Rand) float64 {
    return rand.NormFloat64()
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const lenLetters uint32 = uint32(len(letters))

func RandomString(rand Rand) string {
	res := make([]byte, rand.Uint32()%RandomNatConstraint)
	for i := range res {
		res[i] = letters[rand.Uint32()%lenLetters]
	}
	return string(res)
}

func RandomStringBytes(rand Rand) []byte {
	res := make([]byte, rand.Uint32()%RandomNatConstraint)
	for i := range res {
		res[i] = letters[rand.Uint32()%lenLetters]
	}
	return res
}
`

const internalTLCodeHeader = `%s
package %s

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
)
`

const internalTLCodeBody = `

type UnionElement struct {
	TLTag    uint32
	TLName   string
	TLString string
}

func ErrorClientWrite(typeName string, err error) error {
	return fmt.Errorf("failed to serialize %s request: %w", typeName, err)
}

func ErrorClientDo(typeName string, network string, actorID uint64, address string, err error) error {
	return fmt.Errorf("%s request to %s://%d@%s failed: %w", typeName, network, actorID, address, err)
}

func ErrorClientReadResult(typeName string, network string, actorID uint64, address string, err error) error {
	return fmt.Errorf("failed to deserialize %s response from %s://%d@%s: %w", typeName, network, actorID, address, err)
}

func ErrorServerHandle(typeName string, err error) error {
	return fmt.Errorf("failed to handle %s: %w", typeName, err)
}

func ErrorServerRead(typeName string, err error) error {
	return fmt.Errorf("failed to deserialize %s request: %w", typeName, err)
}

func ErrorServerWriteResult(typeName string, err error) error {
	return fmt.Errorf("failed to serialize %s response: %w", typeName, err)
}

func ErrorInvalidEnumTag(typeName string, tag uint32) error {
	return fmt.Errorf("invalid enum %q tag: 0x%x", typeName, tag)
}

func ErrorInvalidUnionTag(typeName string, tag uint32) error {
	return fmt.Errorf("invalid union %q tag: 0x%x", typeName, tag)
}

func ErrorWrongSequenceLength(typeName string, actual int, expected uint32) error {
	return fmt.Errorf("wrong sequence %q length: %d expected: %d", typeName, actual, expected)
}

func ErrorInvalidEnumTagJSON(typeName string, tag string) error {
	return fmt.Errorf("invalid enum %q tag: %q", typeName, tag)
}

func ErrorInvalidUnionTagJSON(typeName string, tag string) error {
	return fmt.Errorf("invalid union %q tag: %q", typeName, tag)
}

func ErrorInvalidJSON(typeName string, msg string) error {
	return fmt.Errorf("invalid json for type %q - %s", typeName, msg)
}

func ErrorInvalidJSONExcessElement(typeName string, key string) error {
	return fmt.Errorf("invalid json object key %q", key)
}

func JsonReadUnionType(typeName string, j interface{}) (map[string]interface{}, string, error) {
	if j == nil {
		return nil, "", ErrorInvalidJSON(typeName, "expected json object")
	}
	jm, ok := j.(map[string]interface{})
	if !ok {
		return nil, "", ErrorInvalidJSON(typeName, "expected json object")
	}
	jtype, ok := jm["type"]
	if !ok {
		return nil, "", ErrorInvalidJSON(typeName, "expected 'type' key")
	}
	var ret string
	if err := JsonReadString(jtype, &ret); err != nil {
		return nil, "", err
	}
	delete(jm, "type")
	return jm, ret, nil
}

func JsonReadMaybe(typeName string, j interface{}) (bool, interface{}, error) {
	if j == nil {
		return false, nil, nil
	}
	jm, ok := j.(map[string]interface{})
	if !ok {
		return false, nil, ErrorInvalidJSON(typeName, "expected json object")
	}
	jvalue := jm["value"]
	delete(jm, "value")
	jok, ok := jm["ok"]
	delete(jm, "ok")
	var dst bool
	if !ok {
		if jvalue != nil {
			dst = true
		}
	} else {
		if err := JsonReadBool(jok, &dst); err != nil {
			return false, nil, err
		}
		if !dst && jvalue != nil {
			return false, nil, ErrorInvalidJSON(typeName, "if 'ok' is set to false, 'value' should be omitted")
		}
	}
	for k := range jm {
		return false, nil, ErrorInvalidJSONExcessElement(typeName, k)
	}
	return dst, jvalue, nil
}

func JsonReadArray(typeName string, j interface{}) (int, []interface{}, error) {
	var arr []interface{}
	var arrok bool
	if j != nil {
		arr, arrok = j.([]interface{})
		if !arrok {
			return 0, nil, ErrorInvalidJSON(typeName, "expected json array")
		}
	}
	return len(arr), arr, nil
}

func JsonReadArrayFixedSize(typeName string, j interface{}, expectLength uint32) (int, []interface{}, error) {
	l, arr, err := JsonReadArray(typeName, j)
	if err == nil && l != int(expectLength) {
		return 0, nil, ErrorWrongSequenceLength(typeName, l, expectLength)
	}
	return l, arr, err
}

func JsonReadBool(j interface{}, dst *bool) error {
	if j == nil {
		*dst = false
		return nil
	}
	jj, ok := j.(bool)
	if !ok {
		return fmt.Errorf("invalid json for bool")
	}
	*dst = jj
	return nil
}

func JsonReadString(j interface{}, dst *string) error {
	if j == nil {
		*dst = ""
		return nil
	}
	switch jj := j.(type) {
	case string:
		*dst = jj
		return nil
	case map[string]interface{}:
		iface, ok := jj["` + binaryStringObjectKey + `"]
		if !ok {
			return fmt.Errorf("invalid json for string: base64 encoded didn't match as string")
		}
		str, ok := iface.(string)
		if !ok {
			return fmt.Errorf("invalid json for string: unexpected binary string's object")
		}
		buf, err := base64.StdEncoding.DecodeString(str)
		*dst = string(buf)
		return err
	default:
		return fmt.Errorf("invalid json for string")
	}
}

func JsonReadStringBytes(j interface{}, dst *[]byte) error {
	if j == nil {
		*dst = nil
		return nil
	}
	switch jj := j.(type) {
	case string:
		*dst = append((*dst)[:0], jj...)
		return nil
	case map[string]interface{}:
		iface, ok := jj["` + binaryStringObjectKey + `"]
		if !ok {
			return fmt.Errorf("invalid json for string: base64 encoded didn't match as string")
		}
		str, ok := iface.(string)
		if !ok {
			return fmt.Errorf("invalid json for string: unexpected binary string's object")
		}
		buf, err := base64.StdEncoding.DecodeString(str)
		*dst = buf
		return err
	default:
		return fmt.Errorf("invalid json for string")
	}
}

// We allow to specify numbers as "123", so that JS can pass through int64 and bigger numbers
func jsonNumberOrString(j interface{}) (string, bool) {
	jn, ok := j.(json.Number)
	if ok {
		return string(jn), ok
	}
	js, ok := j.(string)
	return js, ok
}

func JsonReadUint32(j interface{}, dst *uint32) error {
	if j == nil {
		*dst = 0
		return nil
	}
	jj, ok := jsonNumberOrString(j)
	if !ok {
		return fmt.Errorf("invalid json for uint32")
	}
	val, err := strconv.ParseUint(jj, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid number format for uint32 %w", err)
	}
	*dst = uint32(val)
	return nil
}

func JsonReadInt32(j interface{}, dst *int32) error {
	if j == nil {
		*dst = 0
		return nil
	}
	jj, ok := jsonNumberOrString(j)
	if !ok {
		return fmt.Errorf("invalid json for int32")
	}
	val, err := strconv.ParseInt(jj, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid number format for int32 %w", err)
	}
	*dst = int32(val)
	return nil
}

func JsonReadInt64(j interface{}, dst *int64) error {
	if j == nil {
		*dst = 0
		return nil
	}
	jj, ok := jsonNumberOrString(j)
	if !ok {
		return fmt.Errorf("invalid json for int64")
	}
	val, err := strconv.ParseInt(jj, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number format for int64 %w", err)
	}
	*dst = val
	return nil
}

func JsonReadFloat32(j interface{}, dst *float32) error {
	if j == nil {
		*dst = 0
		return nil
	}
	jj, ok := jsonNumberOrString(j)
	if !ok {
		return fmt.Errorf("invalid json for float32")
	}
	val, err := strconv.ParseFloat(jj, 32)
	if err != nil {
		return fmt.Errorf("invalid number format for float32 %w", err)
	}
	*dst = float32(val)
	return nil
}

func JsonReadFloat64(j interface{}, dst *float64) error {
	if j == nil {
		*dst = 0
		return nil
	}
	jj, ok := jsonNumberOrString(j)
	if !ok {
		return fmt.Errorf("invalid json for float64")
	}
	val, err := strconv.ParseFloat(jj, 64)
	if err != nil {
		return fmt.Errorf("invalid number format for float64 %w", err)
	}
	*dst = val
	return nil
}

func JsonBytesToInterface(b []byte) (interface{}, error) {
	var j interface{}
	d := json.NewDecoder(bytes.NewBuffer(b))
	d.UseNumber()
	if err := d.Decode(&j); err != nil {
		return j, err
	}
	return j, nil
}
`
