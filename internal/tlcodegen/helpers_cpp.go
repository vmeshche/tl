// Copyright 2022 V Kontakte LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tlcodegen

const basicCPPTLCodeHeader = `%s
#pragma once

#include <stddef.h>
#include <cstring>
#include <stdexcept>
#include <array>
#include <string>
#include <vector>
#include <utility>
#include <memory>

#if __cplusplus >= 201703L
#include <optional>
#include <variant>
#endif

#define TLGEN2_EXOTIC_LAYOUT 0 // when your architecture int32_t, int64_t, float and double types are not the same as in x86
#define TLGEN2_UNLIKELY(x) (x) // __builtin_expect((x), 0) // could improve performance on your platform
#define TLGEN2_NOINLINE //__attribute__ ((noinline)) // could improve performance on your platform

namespace %s {

`

const basicCPPTLCodeFooter = `
} // namespace %s


#undef TLGEN2_NOINLINE
#undef TLGEN2_UNLIKELY
#undef TLGEN2_EXOTIC_LAYOUT
`

const basicCPPTLCodeBody = `

#if __cplusplus >= 201703L
using std::optional;
using std::variant;
using std::string_view;
#else
#error "TODO - if compiling to C++14 standard, define our toy optional and toy variant here."
#endif


enum {
	TL_MAX_TINY_STRING_LEN          = 253,
	TL_BIG_STRING_LEN               = 0xffffff,
	TL_BIG_STRING_MARKER            = 0xfe,
};

class parse_exception : public std::runtime_error {
public:
	explicit parse_exception() : std::runtime_error("TL parse exception") {}
};

class string_too_long : public std::runtime_error {
public:
	explicit string_too_long() : std::runtime_error("TL string too long") {}
};

class string_padding_wrong : public std::runtime_error {
public:
	explicit string_padding_wrong() : std::runtime_error("TL string padding non zero") {}
};

class tag_wrong : public std::runtime_error {
public:
	explicit tag_wrong() : std::runtime_error("TL tag wrong") {}
};

class sequence_length_wrong : public std::runtime_error {
public:
	explicit sequence_length_wrong() : std::runtime_error("TL sequence length wrong") {}
};

class bool_tag_wrong : public std::runtime_error {
public:
	explicit bool_tag_wrong() : std::runtime_error("TL bool tag wrong") {}
};

class union_tag_wrong : public std::runtime_error {
public:
	explicit union_tag_wrong() : std::runtime_error("TL union/enum tag wrong") {}
};

inline TLGEN2_NOINLINE void throwEOF() {
	throw parse_exception{};
}

inline TLGEN2_NOINLINE void throwStringTooLong() {
	throw string_too_long{};
}

inline TLGEN2_NOINLINE void throwStringPaddingWrong() {
	throw string_padding_wrong{};
}

inline TLGEN2_NOINLINE void throwTagWrong() {
	throw tag_wrong{};
}

inline TLGEN2_NOINLINE void throwSequenceLengthWrong() {
	throw sequence_length_wrong{};
}

inline TLGEN2_NOINLINE void throwBoolTagWrong() {
	throw bool_tag_wrong{};
}

inline TLGEN2_NOINLINE void throwUnionTagWrong() {
	throw union_tag_wrong{};
}

inline void tl_pack32(char *buf, uint32_t val) {
#if TLGEN2_EXOTIC_LAYOUT
	buf[0] = char(val);
	buf[1] = char(val >> 8);
	buf[2] = char(val >> 16);
	buf[3] = char(val >> 24);
#else
	std::memcpy(buf, &val, 4);
#endif
}

inline uint32_t tl_unpack32(const char * buf) {
#if TLGEN2_EXOTIC_LAYOUT
	return uint32_t(buf[0] | (buf[1] << 8) | (buf[2] << 16) | (buf[3] << 24));
#else
	uint32_t value = 0;
	std::memcpy(&value, buf, 4);
	return value;
#endif
}

class tl_istream { // TODO - prohibit copy/move
public:
	void nat_read(uint32_t & value) {
		value = nat_read();
	}
	uint32_t nat_read() {
		auto p = advance(4);
		return tl_unpack32(p);
	}
	void nat_read_exact_tag(uint32_t tag) {
		if (TLGEN2_UNLIKELY(tag != nat_read())) {
			throwTagWrong();
		}
	}
	void int_read(int32_t & value) {
		value = int_read();
	}
	int32_t int_read() {
		auto p = advance(4);
		return static_cast<int32_t>(tl_unpack32(p));
	}
	void long_read(int64_t & value) {
		value = long_read();
	}
	int64_t long_read() {
		auto p = advance(8);
		int64_t value = 0;
#if TLGEN2_EXOTIC_LAYOUT
		value = int64_t((long long)(p[0]) | ((long long)p[1] << 8) | ((long long)p[2] << 16) | ((long long)p[3] << 24) | ((long long)p[4] << 32) | ((long long)p[5] << 40) | ((long long)p[6] << 48) | ((long long)p[7] << 56));
#else
		std::memcpy(&value, p, 8);
#endif
		return value;
	}
	void float_read(float & value) {
		value = float_read();
	}
	float float_read() {
		auto p = advance(4);
		float value = 0;
#if TLGEN2_EXOTIC_LAYOUT
		static_assert(false, "Please define conversion from x86 float layout to your architecture");
#else
		std::memcpy(&value, p, 4);
#endif
		return value;
	}
	void double_read(double & value) {
		value = double_read();
	}
	double double_read() {
		auto p = advance(8);
		double value = 0;
#if TLGEN2_EXOTIC_LAYOUT
		static_assert(false, "Please define conversion from x86 double layout to your architecture");
#else
		std::memcpy(&value, p, 8);
#endif
		return value;
	}
	std::string string_read() {
		std::string value;
		string_read(value);
		return value;
	}
	bool bool_read(uint32_t f, uint32_t t) {
		auto tag = nat_read();
		if (tag == t) { return true; }
		if (tag != f) { throwBoolTagWrong(); }
		return false;
	}
	void string_read(std::string & value) {
		ensure(4);
		auto len = size_t(static_cast<unsigned char>(*ptr));
		if (TLGEN2_UNLIKELY(len >= TL_BIG_STRING_MARKER)) {
			if (TLGEN2_UNLIKELY(len > TL_BIG_STRING_MARKER)) {
				throwStringTooLong();
			}
			len = tl_unpack32(ptr) >> 8U;
			ptr += 4;
			value.clear();
			fetch_data_append(value, len);
			fetch_pad((-len) & 3);
			return;
		}
		auto pad = ((-(len+1)) & 3);
		auto fullLen = 1 + len + pad;
		if (TLGEN2_UNLIKELY(ptr + fullLen > end)) {
			ptr += 1;
			value.clear();
			fetch_data_append(value, len);
			fetch_pad(pad);
			return;
		}
		// fast path for short strings that fully fit in buffer
		auto x = tl_unpack32(ptr + fullLen - 4);
		if (TLGEN2_UNLIKELY((x & ~(0xFFFFFFFFU >> (8*pad))) != 0)) {
			throwStringPaddingWrong();
		}
		value.assign(ptr + 1, len);
		ptr += fullLen;
	}
protected:
	const char * ptr{};
	const char * end{};
	virtual void grow_buffer(size_t size) = 0; // after call buffer must contain at least 8 bytes
private:
	void fetch_data(char * data, size_t size) {
		for (;TLGEN2_UNLIKELY(ptr + size > end);) {
			std::memcpy(data, ptr, end - ptr);
			data += end - ptr;
			size -= end - ptr;
			ptr = end;
			grow_buffer(size);
			if (TLGEN2_UNLIKELY(ptr == end)) {
				throwEOF();
			}
		}
		std::memcpy(data, ptr, size);
		ptr += size;
	}
	void fetch_data_append(std::string & value, size_t size) {
		for (;TLGEN2_UNLIKELY(ptr + size > end);) {
			value.append(ptr, end - ptr);
			size -= end - ptr;
			ptr = end;
			grow_buffer(size);
			if (TLGEN2_UNLIKELY(ptr == end)) {
				throwEOF();
			}
		}
		value.append(ptr, size);
		ptr += size;
	}
	void fetch_pad(size_t len) {
		auto p = advance(len);
		uint32_t x = 0;
		std::memcpy(&x, p, len);
		if (TLGEN2_UNLIKELY(x != 0)) {
			throwStringPaddingWrong();
		}
	}
	const char * advance(size_t size) {
		ensure(size);
		auto p = ptr;
		ptr += size;
		return p;
	}
	void ensure(size_t size) {
		if (TLGEN2_UNLIKELY(ptr + size > end)) {
			grow_buffer(size);
			if (TLGEN2_UNLIKELY(ptr + size > end)) {
				throwEOF();
			}
		}
	}
};

class tl_ostream { // TODO - prohibit copy/move
public:
	void nat_write(uint32_t value) {
		auto p = advance(4);
		tl_pack32(p, value);
	}
	void int_write(int32_t value) {
		auto p = advance(4);
		tl_pack32(p, static_cast<uint32_t>(value));
	}
	void long_write(int64_t value) {
		auto p = advance(8);
#if TLGEN2_EXOTIC_LAYOUT
		ptr[0] = char(value);
		ptr[1] = char(value >> 8);
		ptr[2] = char(value >> 16);
		ptr[3] = char(value >> 24);
		ptr[4] = char(value >> 32);
		ptr[5] = char(value >> 40);
		ptr[6] = char(value >> 48);
		ptr[7] = char(value >> 56);
#else
		std::memcpy(p, &value, 8);
#endif
	}
	void float_write(float value) {
		auto p = advance(4);
#if TLGEN2_EXOTIC_LAYOUT
		static_assert(false, "Please define conversion to x86 float layout from your architecture");
#else
		std::memcpy(p, &value, 4);
#endif
	}
	void double_write(double value) {
		auto p = advance(8);
#if TLGEN2_EXOTIC_LAYOUT
		static_assert(false, "Please define conversion from x86 double layout to your architecture");
#else
		std::memcpy(p, &value, 8);
#endif
	}
	void string_write(const std::string & value) {
		auto len = value.size();
		if (TLGEN2_UNLIKELY(len > TL_MAX_TINY_STRING_LEN)) {
			if (TLGEN2_UNLIKELY(len > TL_BIG_STRING_LEN)) {
				throwStringTooLong();
			}
			auto p = advance(4);
			tl_pack32(p, (len << 8U) | TL_BIG_STRING_MARKER);
			store_data(value.data(), value.size());
			store_pad((-len) & 3);
			return;
		}
		auto pad = ((-(len+1)) & 3);
		auto fullLen = 1 + len + pad;
		if (TLGEN2_UNLIKELY(ptr + fullLen > end)) {
			auto p = advance(1);
			*p = static_cast<char>(len);
			store_data(value.data(), value.size());
			store_pad(pad);
			return;
		}
		tl_pack32(ptr + fullLen - 4, 0); // padding first
		tl_pack32(ptr, len);
		std::memcpy(ptr + 1, value.data(), len);
		ptr += fullLen;
	}
protected:
	char * ptr{};
	char * end{};
	virtual void grow_buffer(size_t size) = 0; // after call buffer must contain at least 8 bytes
private:
	void store_data(const char * data, size_t size) {
		for (;TLGEN2_UNLIKELY(ptr + size > end);) {
			std::memcpy(ptr, data, end - ptr);
			data += end - ptr;
			size -= end - ptr;
			ptr = end;
			grow_buffer(size);
			if (TLGEN2_UNLIKELY(ptr == end)) {
				throwEOF();
			}
		}
		std::memcpy(ptr, data, size);
		ptr += size;
	}
	void store_pad(size_t len) {
		auto p = advance(len);
		if (len != 0) {
			p[0] = 0;
			p[len-1] = 0;
			p[len/2] = 0;
		}
	}
	char * advance(size_t size) {
		ensure(size);
		auto p = ptr;
		ptr += size;
		return p;
	}
	void ensure(size_t size) {
		if (TLGEN2_UNLIKELY(ptr + size > end)) {
			grow_buffer(size);
			if (TLGEN2_UNLIKELY(ptr + size > end)) {
				throwEOF();
			}
		}
	}
};

class tl_istream_string : public tl_istream { // TODO - custom copy/move
public:
	explicit tl_istream_string(const std::string & buf) {
		ptr = buf.data();
		end = ptr + buf.size();
	}
protected:
	void grow_buffer(size_t size) override { throwEOF(); }
};

class tl_ostream_string : public tl_ostream { // TODO - custom copy/move
public:
	explicit tl_ostream_string() {
		resize(INITIAL_SIZE);
	}
	std::string & get_buffer() {
		resize(ptr - buf.data());
		return buf;
	}
protected:
	void grow_buffer(size_t size) override {
		auto pos = ptr - buf.data();
		resize(buf.size()*3/2 + INITIAL_SIZE + size); // some arbitrary strategy
		ptr += pos;
	}
private:
	enum { INITIAL_SIZE = 1024 };
	std::string buf;
	void resize(size_t size) {
		buf.resize(size);
		ptr = const_cast<char *>(buf.data()); // works for all known implementations
		end = ptr + buf.size();
	}
};

`
