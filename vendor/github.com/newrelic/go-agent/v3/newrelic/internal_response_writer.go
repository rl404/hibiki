// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package newrelic

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type replacementResponseWriter struct {
	thd      *thread
	original http.ResponseWriter
}

func (rw *replacementResponseWriter) Header() http.Header {
	return rw.original.Header()
}

func (rw *replacementResponseWriter) Write(b []byte) (n int, err error) {
	hdr := rw.original.Header()

	// This is safe to call unconditionally, even if Write() is called multiple
	// times; see also the commentary in addCrossProcessHeaders().
	addCrossProcessHeaders(rw.thd.txn, hdr)

	n, err = rw.original.Write(b)

	headersJustWritten(rw.thd, http.StatusOK, hdr)

	secureAgent.SendEvent("INBOUND_WRITE", string(b), hdr)
	return
}

func (rw *replacementResponseWriter) WriteHeader(code int) {
	hdr := rw.original.Header()

	addCrossProcessHeaders(rw.thd.txn, hdr)

	rw.original.WriteHeader(code)

	headersJustWritten(rw.thd, code, hdr)
}

func (rw *replacementResponseWriter) CloseNotify() <-chan bool {
	return rw.original.(http.CloseNotifier).CloseNotify()
}
func (rw *replacementResponseWriter) Flush() {
	rw.original.(http.Flusher).Flush()
}
func (rw *replacementResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.original.(http.Hijacker).Hijack()
}
func (rw *replacementResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	return rw.original.(io.ReaderFrom).ReadFrom(r)
}

func upgradeResponseWriter(rw *replacementResponseWriter) http.ResponseWriter {
	// GENERATED CODE DO NOT MODIFY
	// This code generated by internal/tools/interface-wrapping
	var (
		i0 int32 = 1 << 0
		i1 int32 = 1 << 1
		i2 int32 = 1 << 2
		i3 int32 = 1 << 3
	)
	var interfaceSet int32
	if _, ok := rw.original.(http.CloseNotifier); ok {
		interfaceSet |= i0
	}
	if _, ok := rw.original.(http.Flusher); ok {
		interfaceSet |= i1
	}
	if _, ok := rw.original.(http.Hijacker); ok {
		interfaceSet |= i2
	}
	if _, ok := rw.original.(io.ReaderFrom); ok {
		interfaceSet |= i3
	}
	switch interfaceSet {
	default: // No optional interfaces implemented
		return struct {
			http.ResponseWriter
		}{rw}
	case i0:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
		}{rw, rw}
	case i1:
		return struct {
			http.ResponseWriter
			http.Flusher
		}{rw, rw}
	case i0 | i1:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
		}{rw, rw, rw}
	case i2:
		return struct {
			http.ResponseWriter
			http.Hijacker
		}{rw, rw}
	case i0 | i2:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
		}{rw, rw, rw}
	case i1 | i2:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
		}{rw, rw, rw}
	case i0 | i1 | i2:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
			http.Hijacker
		}{rw, rw, rw, rw}
	case i3:
		return struct {
			http.ResponseWriter
			io.ReaderFrom
		}{rw, rw}
	case i0 | i3:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			io.ReaderFrom
		}{rw, rw, rw}
	case i1 | i3:
		return struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
		}{rw, rw, rw}
	case i0 | i1 | i3:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
			io.ReaderFrom
		}{rw, rw, rw, rw}
	case i2 | i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw}
	case i0 | i2 | i3:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw, rw}
	case i1 | i2 | i3:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw, rw}
	case i0 | i1 | i2 | i3:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw, rw, rw}
	}
}
