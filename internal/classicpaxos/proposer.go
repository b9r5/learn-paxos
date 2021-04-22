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
	"time"
)

// proposer represents the proposer role in Classic Paxos.
type proposer struct {
	input      <-chan message   // input channel
	replyTo    chan<- message   // for acceptor replies
	id         int              // proposer identifier
	nProposers int              // number of proposers
	acceptors  []chan<- message // input channels for acceptors
	timeout    time.Duration    // time to wait for promise and accept messages
	values     chan<- string    // proposer places agreed value on this channel
}

// newProposer creates a proposer with the given parameters and starts its
// goroutine.
func newProposer(id, nProposers int,
	input <-chan message,
	replyTo chan<- message,
	acceptorChannels []chan<- message,
	timeout time.Duration,
	values chan<- string) *proposer {

	p := &proposer{
		input:      input,
		replyTo:    replyTo,
		id:         id,
		nProposers: nProposers,
		acceptors:  acceptorChannels,
		timeout:    timeout,
		values:     values,
	}
	go p.run()
	return p
}

// run is a translation of the proposer algorithm for Classic Paxos. It returns
// the agreed value.
func (p *proposer) run() string {
	var epoch Epoch

	candidateValue := fmt.Sprintf("v%d", p.id)

	for {
		var value = ""                          // current proposal value
		var maxEpoch Epoch                      // maximum epoch received in phase 1
		promisedAcceptors := make(map[int]bool) // keys are acceptors that have promised
		acceptedAcceptors := make(map[int]bool) // keys are acceptors that have accepted

		// select and set the epoch
		if epoch.Nil() {
			epoch = newEpoch(p.id, p.nProposers)
		} else {
			epoch = epoch.Next()
		}

		for _, a := range p.acceptors {
			a <- prepare{epoch: epoch, replyTo: p.replyTo, proposerID: p.id}
		}

		timedOut := false

		for !timedOut && len(promisedAcceptors) < (len(p.acceptors)/2)+1 {
			select {
			case msg := <-p.input:
				fmt.Printf("proposer %d received message %s\n", p.id, msg)

				if promise, ok := msg.(promise); ok {
					promisedAcceptors[promise.acceptorID] = true
					if !promise.acceptedEpoch.Nil() &&
						(maxEpoch.Nil() || promise.acceptedEpoch.Cmp(maxEpoch) > 0) {

						// (maxEpoch, value) is the greatest proposal received
						maxEpoch = promise.acceptedEpoch
						value = promise.acceptedValue
					}
				}
			case <-time.After(p.timeout):
				timedOut = true
			}
		}

		if timedOut {
			continue
		}

		if value == "" {
			// no proposals were received thus propose candidate value
			value = candidateValue
		}

		// start phase 2 for proposal (epoch, value)
		for _, a := range p.acceptors {
			a <- propose{epoch: epoch, value: value, replyTo: p.replyTo,
				proposerID: p.id}
		}

		for !timedOut && len(acceptedAcceptors) < (len(p.acceptors)/2)+1 {
			select {
			case msg := <-p.input:
				fmt.Printf("proposer %d received message %s\n", p.id, msg)

				if accept, ok := msg.(accept); ok && epoch.Cmp(accept.epoch) == 0 {
					acceptedAcceptors[accept.acceptorID] = true
				}
			case <-time.After(p.timeout):
				timedOut = true
			}
		}

		if timedOut {
			continue
		}

		fmt.Printf("proposer %d believes value %s is decided\n", p.id, value)

		p.values <- value
	}
}
