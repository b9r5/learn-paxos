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
	"testing"
	"time"
)

// TestAgreement tests agreement in the Classic Paxos implementation with a
// variety of numbers of proposers, acceptors; and with no message loss or
// reordering, and a moderate amount of message loss and reordering.
func TestAgreement(t *testing.T) {
	type testCase struct {
		nProposers, nAcceptors, buffer int
		drop                           float64
	}

	cases := []testCase{
		// no message loss or reordering:
		{nProposers: 1, nAcceptors: 1, buffer: 1, drop: 0.0},
		{nProposers: 2, nAcceptors: 1, buffer: 1, drop: 0.0},
		{nProposers: 2, nAcceptors: 3, buffer: 1, drop: 0.0},
		{nProposers: 2, nAcceptors: 5, buffer: 1, drop: 0.0},
		{nProposers: 10, nAcceptors: 5, buffer: 1, drop: 0.0},

		// some message loss and reordering:
		{nProposers: 1, nAcceptors: 1, buffer: 2, drop: 0.1},
		{nProposers: 2, nAcceptors: 5, buffer: 2, drop: 0.1},
		{nProposers: 10, nAcceptors: 5, buffer: 2, drop: 0.1},
	}

	for _, c := range cases {
		if err := runPaxos(c.nProposers, c.nAcceptors,
			100*time.Millisecond, 10*time.Millisecond,
			c.buffer, c.drop); err != nil {

			t.Errorf("with %d proposer/s and %d acceptor/s, got %v",
				c.nProposers, c.nAcceptors, err)
		}
	}
}
