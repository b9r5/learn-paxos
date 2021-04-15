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
	"fmt"
	"sync"
	"time"
)

// Config represents configuration for Classic Paxos, including number of
// proposers, number of acceptors, proposer timeout, and lossyChannel
// parameters.
type Config struct {
	// number of proposers
	NProposers int

	// number of acceptors
	NAcceptors int

	// how long proposer waits for acceptor responses before re-proposing
	ProposerTimeout time.Duration

	// how long lossyChannel waits for buffer to fill before returning message
	ChannelTimeout time.Duration

	// how many messages lossyChannel buffers
	Buffer int

	// probability that lossyChannel drops a message
	Drop float64
}

// Run runs Classic Paxos for the scenario given by the configuration c.
func (c *Config) Run() error {
	// 1. create acceptors
	acceptorChannels := c.newAcceptors()

	// 2. create proposers
	valueChannel := c.newProposers(acceptorChannels)

	// 3. check whether proposers agreed on same value
	return c.checkValues(valueChannel)
}

// newAcceptors creates c.NAcceptors acceptors. It returns the channels that are
// inputs to the acceptors' lossy channels for use by the proposers.
func (c *Config) newAcceptors() []chan<- message {
	channels := make([]*lossyChannel, c.NAcceptors)
	for i := 0; i < c.NAcceptors; i++ {
		channels[i] = newLossyChannel(c.Buffer, c.ChannelTimeout, c.Drop)
	}

	acceptors := make([]*acceptor, c.NAcceptors)
	for i := 0; i < c.NAcceptors; i++ {
		acceptors[i] = newAcceptor(i, channels[i].output)
	}

	result := make([]chan<- message, c.NAcceptors)
	for i := 0; i < c.NAcceptors; i++ {
		result[i] = channels[i].input
	}

	return result
}

// newProposers creates c.NProposers proposers. acceptorChannels are the
// acceptors' input channels. It returns a channel of the values that each
// proposers believes was agreed.
func (c *Config) newProposers(acceptorChannels []chan<- message) <-chan string {
	proposers := make([]*proposer, c.NProposers)

	var group sync.WaitGroup
	group.Add(c.NProposers)

	channels := make([]*lossyChannel, c.NProposers)
	for i := 0; i < c.NProposers; i++ {
		channels[i] = newLossyChannel(c.Buffer, c.ChannelTimeout, c.Drop)
	}

	valueChannel := make(chan string, c.NProposers)

	for i := 0; i < c.NProposers; i++ {
		proposers[i] = newProposer(i, c.NProposers, channels[i].output,
			channels[i].input, acceptorChannels, c.ProposerTimeout,
			valueChannel)
	}

	return valueChannel
}

// checkValues waits until n values appear on the values channel, checks whether
// the values are identical, and returns a non-nil error if they differ.
func (c *Config) checkValues(values <-chan string) error {
	vals := make([]string, 0, c.NAcceptors)

	problem := false

	for i := 0; i < c.NProposers; i++ {
		vals = append(vals, <-values)
		if i > 0 && vals[i] != vals[0] {
			fmt.Printf(
				"uh oh! 2 proposers believe different values were agreed (%s versus %s)\n",
				vals[0], vals[i])
			problem = true
		}
	}

	if !problem {
		fmt.Printf(
			"yay! %d values were agreed, and they were all the same (%s)\n",
			c.NProposers, vals[0])
		return nil
	} else {
		return fmt.Errorf(
			"the following values were agreed, and some of them differ: %v\n", vals)
	}
}
