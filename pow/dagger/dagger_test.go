// Copyright 2014 The go-pbfcoin Authors
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

package dagger

import (
	"math/big"
	"testing"

	"github.com/pbfcoin/go-pbfcoin/common"
)

func BenchmarkDaggerSearch(b *testing.B) {
	hash := big.NewInt(0)
	diff := common.BigPow(2, 36)
	o := big.NewInt(0) // nonce doesn't matter. We're only testing against speed, not validity

	// Reset timer so the big generation isn't included in the benchmark
	b.ResetTimer()
	// Validate
	DaggerVerify(hash, diff, o)
}
