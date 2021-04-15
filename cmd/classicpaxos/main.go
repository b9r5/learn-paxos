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
	"flag"
	"fmt"
	"github.com/b9r5/learn-paxos/internal/classicpaxos"
	"os"
	"time"
)

// main starts a number of proposers and acceptors, waits until all proposers
// have finished the proposer algorithm, and checks that the proposers agreed
// on the same value.
func main() {
	var nProposers = flag.Int("proposers", 10, "number of proposers")
	var nAcceptors = flag.Int("acceptors", 5, "number of acceptors")
	var proposerTimeout = flag.Duration("proposer-timeout",
		100*time.Millisecond,
		"time for proposer to wait for promise and accept messages")
	var channelTimeout = flag.Duration("channel-timeout", 10*time.Millisecond,
		"time to wait for lossy channel buffer to fill before returning a message")
	var buffer = flag.Int("buffer-size", 2,
		"number of messages to buffer before returning one selected randomly")
	var drop = flag.Float64("drop-probability", 0.1,
		"probability of lossy channel dropping a message, in range [0, 1)")

	c := classicpaxos.Config{
		NProposers:      *nProposers,
		NAcceptors:      *nAcceptors,
		ProposerTimeout: *proposerTimeout,
		ChannelTimeout:  *channelTimeout,
		Buffer:          *buffer,
		Drop:            *drop,
	}

	if err := c.Run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
