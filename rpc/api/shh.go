// Copyright 2015 The go-pbfcoin Authors
// This file is part of the go-pbfcoin library.
//
// The go-pbfcoin library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-pbfcoin library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-pbfcoin library. If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"math/big"

	"github.com/pbfcoin/go-pbfcoin/pbf"
	"github.com/pbfcoin/go-pbfcoin/rpc/codec"
	"github.com/pbfcoin/go-pbfcoin/rpc/shared"
	"github.com/pbfcoin/go-pbfcoin/xpbf"
)

const (
	ShhApiVersion = "1.0"
)

var (
	// mapping between methods and handlers
	shhMapping = map[string]shhhandler{
		"shh_version":          (*shhApi).Version,
		"shh_post":             (*shhApi).Post,
		"shh_hasIdentity":      (*shhApi).HasIdentity,
		"shh_newIdentity":      (*shhApi).NewIdentity,
		"shh_newFilter":        (*shhApi).NewFilter,
		"shh_uninstallFilter":  (*shhApi).UninstallFilter,
		"shh_getMessages":      (*shhApi).GetMessages,
		"shh_getFilterChanges": (*shhApi).GetFilterChanges,
	}
)

func newWhisperOfflineError(method string) error {
	return shared.NewNotAvailableError(method, "whisper offline")
}

// net callback handler
type shhhandler func(*shhApi, *shared.Request) (interface{}, error)

// shh api provider
type shhApi struct {
	xpbf    *xpbf.Xpbf
	pbfcoin *pbf.pbfcoin
	methods map[string]shhhandler
	codec   codec.ApiCoder
}

// create a new whisper api instance
func NewShhApi(xpbf *xpbf.Xpbf, pbf *pbf.pbfcoin, coder codec.Codec) *shhApi {
	return &shhApi{
		xpbf:    xpbf,
		pbfcoin: pbf,
		methods: shhMapping,
		codec:   coder.New(nil),
	}
}

// collection with supported methods
func (self *shhApi) methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *shhApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.method]; ok {
		return callback(self, req)
	}

	return nil, shared.NewNotImplementedError(req.method)
}

func (self *shhApi) Name() string {
	return shared.ShhApiName
}

func (self *shhApi) ApiVersion() string {
	return ShhApiVersion
}

func (self *shhApi) Version(req *shared.Request) (interface{}, error) {
	w := self.xpbf.Whisper()
	if w == nil {
		return nil, newWhisperOfflineError(req.method)
	}

	return w.Version(), nil
}

func (self *shhApi) Post(req *shared.Request) (interface{}, error) {
	w := self.xpbf.Whisper()
	if w == nil {
		return nil, newWhisperOfflineError(req.method)
	}

	args := new(WhisperMessageArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, err
	}

	err := w.Post(args.Payload, args.To, args.From, args.Topics, args.Priority, args.Ttl)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (self *shhApi) HasIdentity(req *shared.Request) (interface{}, error) {
	w := self.xpbf.Whisper()
	if w == nil {
		return nil, newWhisperOfflineError(req.method)
	}

	args := new(WhisperIdentityArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, err
	}

	return w.HasIdentity(args.Identity), nil
}

func (self *shhApi) NewIdentity(req *shared.Request) (interface{}, error) {
	w := self.xpbf.Whisper()
	if w == nil {
		return nil, newWhisperOfflineError(req.method)
	}

	return w.NewIdentity(), nil
}

func (self *shhApi) NewFilter(req *shared.Request) (interface{}, error) {
	args := new(WhisperFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, err
	}

	id := self.xpbf.NewWhisperFilter(args.To, args.From, args.Topics)
	return newHexNum(big.NewInt(int64(id)).Bytes()), nil
}

func (self *shhApi) UninstallFilter(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, err
	}
	return self.xpbf.UninstallWhisperFilter(args.Id), nil
}

func (self *shhApi) GetFilterChanges(req *shared.Request) (interface{}, error) {
	w := self.xpbf.Whisper()
	if w == nil {
		return nil, newWhisperOfflineError(req.method)
	}

	// Retrieve all the new messages arrived since the last request
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, err
	}

	return self.xpbf.WhisperMessagesChanged(args.Id), nil
}

func (self *shhApi) GetMessages(req *shared.Request) (interface{}, error) {
	w := self.xpbf.Whisper()
	if w == nil {
		return nil, newWhisperOfflineError(req.method)
	}

	// Retrieve all the cached messages matching a specific, existing filter
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, err
	}

	return self.xpbf.WhisperMessages(args.Id), nil
}
