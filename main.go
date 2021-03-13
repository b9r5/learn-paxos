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
	"os"
	"sync"
	"time"
)

// main starts a number of proposers and acceptors, stops a minority of
// acceptors, and waits until all proposers have finished the proposer algorithm
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

	err := runPaxos(*nProposers, *nAcceptors, *proposerTimeout, *channelTimeout,
		*buffer, *drop)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func runPaxos(nProposers, nAcceptors int,
	proposerTimeout, channelTimeout time.Duration,
	buffer int, drop float64) error {

	// 1. create acceptors
	acceptorChannels := newAcceptors(nAcceptors, buffer, channelTimeout, drop)

	// 2. create proposers
	valueChannel := newProposers(nProposers, buffer, channelTimeout,
		proposerTimeout, drop, acceptorChannels)

	// 3. check whether proposers agreed on same value
	return checkValues(nProposers, valueChannel)
}

// newAcceptors creates n acceptors. The acceptors' lossy channels have the
// supplied parameters buffer, timeout, and drop. It returns a slice of input
// channels for use by the proposers.
func newAcceptors(n, buffer int, timeout time.Duration, drop float64) []chan<- message {
	channels := make([]*lossyChannel, n)
	for i := 0; i < n; i++ {
		channels[i] = newLossyChannel(buffer, timeout, drop)
	}

	acceptors := make([]*acceptor, n)
	for i := 0; i < n; i++ {
		acceptors[i] = newAcceptor(i, channels[i].output)
	}

	result := make([]chan<- message, n)
	for i := 0; i < n; i++ {
		result[i] = channels[i].input
	}

	return result
}

// newProposers creates n proposers. The proposers' lossy channels have the
// supplied parameters buffer, channelTimeout, and drop. proposerTimeout is used
// for the proposers' timeout values. acceptorChannels are the acceptors' input
// channels. It returns a *sync.WaitGroup to be used for waiting for all
// proposers to finish the proposer algorithm.
func newProposers(n, buffer int,
	channelTimeout, proposerTimeout time.Duration,
	drop float64,
	acceptorChannels []chan<- message) <-chan string {

	proposers := make([]*proposer, n)

	var group sync.WaitGroup
	group.Add(n)

	channels := make([]*lossyChannel, n)
	for i := 0; i < n; i++ {
		channels[i] = newLossyChannel(buffer, channelTimeout, drop)
	}

	valueChannel := make(chan string, n)

	for i := 0; i < n; i++ {
		proposers[i] = newProposer(i, n, channels[i].output, channels[i].input,
			acceptorChannels, proposerTimeout, valueChannel)
	}

	return valueChannel
}

// checkValues checks whether n identical values appear on the values channel,
// and returns a non-nil error if not.
func checkValues(n int, values <-chan string) error {
	vals := make([]string, 0, n)

	problem := false

	for i := 0; i < n; i++ {
		vals = append(vals, <-values)
		if i > 0 && vals[i] != vals[0] {
			fmt.Printf("uh oh! 2 proposers believe different values were agreed (%s versus %s)\n",
				vals[0], vals[i])
			problem = true
		}
	}

	if !problem {
		fmt.Printf("yay! %d values were agreed, and they were all the same (%s)\n",
			n, vals[0])
		return nil
	} else {
		return fmt.Errorf("the following values were agreed, and some of them differ: %v\n",
			vals)
	}
}
