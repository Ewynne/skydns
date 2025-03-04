// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import (
	"github.com/miekg/dns"
	"github.com/tyler-cromwell/skydns/msg"
)

type Backend interface {
	HasSynced() bool
	Records(name string, requesterIP string, exact bool) ([]msg.Service, error)
	ReverseRecord(name string) (*msg.Service, error)
	TranslateForwardedRequest(name string, requesterIP string, req *dns.Msg)(*dns.Msg, error)
	TranslateForwardedResponse(name string, responderIP string, req, resp *dns.Msg)(*dns.Msg, error)
}

// FirstBackend exposes the Backend interface over multiple Backends, returning
// the first Backend that answers the provided record request. If no Backend answers
// a record request, the last error seen will be returned.
type FirstBackend []Backend

// FirstBackend implements Backend
var _ Backend = FirstBackend{}

func (g FirstBackend) Records(name string, requesterIP string, exact bool) (records []msg.Service, err error) {
	var lastError error
	for _, backend := range g {
		if records, err = backend.Records(name, requesterIP, exact); err == nil && len(records) > 0 {
			return records, nil
		}
		if err != nil {
			lastError = err
		}
	}
	return nil, lastError
}

func (g FirstBackend) ReverseRecord(name string) (record *msg.Service, err error) {
	var lastError error
	for _, backend := range g {
		if record, err = backend.ReverseRecord(name); err == nil && record != nil {
			return record, nil
		}
		if err != nil {
			lastError = err
		}
	}
	return nil, lastError
}

func (g FirstBackend) HasSynced() bool {
	// Stub implementation only to satisfy interface.
	return true
}

func (g FirstBackend) TranslateForwardedRequest(name string, requesterIP string, req *dns.Msg)(*dns.Msg, error) {
	return req, nil
}

func (g FirstBackend) TranslateForwardedResponse(name string, responderIP string, req, resp *dns.Msg)(*dns.Msg, error) {
	return resp, nil
}
