// Copyright 2021 Benjamin Horowitz
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//               http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"math/big"
)

// epoch is used to order proposals in Classic Paxos. The algorithm requires the
// set of epochs be infinite and totally ordered. The epochs that may be used
// by different proposers must be distinct. Finally, no proposer may re-use in
// phase 1 any epoch that it has already used.
//
// To implement these requirements, we use a big.Int to represent an epoch.
// Given n proposers, proposer p (0 <= p < n) uses the integers I equivalent to
// p modulo n, i.e., the integers p, n + p, 2n + p, 3n + p, ...
//
// A typical Paxos implementation would use a 32- or 64-bit integer instead of a
// big.Int. To sidestep the small complexities of overflow, we use a big.Int
// instead.
type epoch struct {
	i          *big.Int
	nProposers int // total number of proposers
}

// newEpoch returns the initial epoch for the proposer numbered proposerID.
// nProposers is the number of proposers.
func newEpoch(proposerID, nProposers int) epoch {
	return epoch{
		i:          big.NewInt(int64(proposerID)),
		nProposers: nProposers,
	}
}

// next returns the next epoch, equal to e + e.proposers.
func (e epoch) next() epoch {
	return epoch{
		i:          e.i.Add(e.i, big.NewInt(int64(e.nProposers))),
		nProposers: e.nProposers,
	}
}

// nil returns true if and only if e is a nil epoch.
func (e epoch) nil() bool {
	return e.i == nil
}

// cmp compares e and f and returns:
//
//   -1 if e <  f
//    0 if e == f
//   +1 if e >  f
func (e epoch) cmp(f epoch) int {
	return e.i.Cmp(f.i)
}

// String returns the string version of an epoch.
func (e epoch) String() string {
	if e.nil() {
		return "nil"
	}
	return fmt.Sprintf("%s", e.i)
}

// nilEpoch is a nil epoch.
var nilEpoch = epoch{
	i:          nil,
	nProposers: 0,
}
