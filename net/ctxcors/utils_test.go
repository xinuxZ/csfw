// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ctxcors

import (
	"strings"
	"testing"
)

func TestWildcard(t *testing.T) {
	w := wildcard{"foo", "bar"}
	if !w.match("foobar") {
		t.Error("foo*bar should match foobar")
	}
	if !w.match("foobazbar") {
		t.Error("foo*bar should match foobazbar")
	}
	if w.match("foobaz") {
		t.Error("foo*bar should not match foobaz")
	}

	w = wildcard{"foo", "oof"}
	if w.match("foof") {
		t.Error("foo*oof should not match foof")
	}
}

func TestConvert(t *testing.T) {
	s := convert([]string{"A", "b", "C"}, strings.ToLower)
	e := []string{"a", "b", "c"}
	if s[0] != e[0] || s[1] != e[1] || s[2] != e[2] {
		t.Errorf("%v != %v", s, e)
	}
}

func TestParseHeaderList(t *testing.T) {
	h := parseHeaderList("header, second-header, THIRD-HEADER, Numb3r3d-H34d3r")
	e := []string{"Header", "Second-Header", "Third-Header", "Numb3r3d-H34d3r"}
	if h[0] != e[0] || h[1] != e[1] || h[2] != e[2] {
		t.Errorf("%v != %v", h, e)
	}
}

func TestParseHeaderListEmpty(t *testing.T) {
	if len(parseHeaderList("")) != 0 {
		t.Error("should be empty sclice")
	}
	if len(parseHeaderList(" , ")) != 0 {
		t.Error("should be empty sclice")
	}
}

var parseHeaderListResult []string

func BenchmarkParseHeaderListConvert(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseHeaderListResult = parseHeaderList("header, second-header, THIRD-HEADER")
	}
}

func BenchmarkParseHeaderListSingle(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseHeaderListResult = parseHeaderList("header")
	}
}

func BenchmarkParseHeaderListNormalized(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseHeaderListResult = parseHeaderList("Header1, Header2, Third-Header")
	}
}
