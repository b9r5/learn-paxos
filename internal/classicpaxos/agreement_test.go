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

package classicpaxos

import (
	"testing"
	"time"
)

// TestAgreement tests agreement in the Classic Paxos implementation with a
// variety of numbers of proposers and acceptors; and with either no message
// loss and reordering, or a moderate amount of message loss and reordering.
func TestAgreement(t *testing.T) {
	pt := 100*time.Millisecond // proposer timeout
	ct := 10*time.Millisecond  // channel timeout

	configs := []Config{
		// no message loss or reordering:
		{NProposers: 1, NAcceptors: 1, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 1, Drop: 0.0},
		{NProposers: 2, NAcceptors: 1, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 1, Drop: 0.0},
		{NProposers: 2, NAcceptors: 3, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 1, Drop: 0.0},
		{NProposers: 2, NAcceptors: 5, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 1, Drop: 0.0},
		{NProposers: 10, NAcceptors: 5, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 1, Drop: 0.0},

		// some message loss and reordering:
		{NProposers: 1, NAcceptors: 1, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 2, Drop: 0.1},
		{NProposers: 2, NAcceptors: 5, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 2, Drop: 0.1},
		{NProposers: 10, NAcceptors: 5, ProposerTimeout: pt, ChannelTimeout: ct,
			Buffer: 2, Drop: 0.1},
	}

	for i, c := range configs {
		if err := c.Run(); err != nil {
			t.Errorf("with %d proposer/s and %d acceptor/s (test case %d), got %v",
				c.NProposers, c.NAcceptors, i, err)
		}
	}
}
