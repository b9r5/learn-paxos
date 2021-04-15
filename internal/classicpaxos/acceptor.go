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

import "fmt"

// acceptor represents the acceptor role in Classic Paxos.
type acceptor struct {
	input <-chan message // input channel
	id    int            // acceptor identifier
}

// newAcceptor creates an acceptor with the given id and input channel, and
// starts its goroutine.
func newAcceptor(id int, input <-chan message) *acceptor {
	a := &acceptor{input: input, id: id}
	go a.run()
	return a
}

// run is a translation of the acceptor algorithm for Classic Paxos.
func (a *acceptor) run() {
	promisedEpoch := nilEpoch
	acceptedEpoch := nilEpoch
	acceptedValue := ""

	for {
		m := <-a.input

		fmt.Printf("acceptor %d received message %s\n", a.id, m)

		switch msg := m.(type) {
		case prepare:
			epoch := msg.epoch
			if promisedEpoch.nil() || epoch.cmp(promisedEpoch) >= 0 {
				promisedEpoch = epoch
				msg.replyTo <- promise{acceptorID: a.id, epoch: epoch,
					acceptedEpoch: acceptedEpoch, acceptedValue: acceptedValue}
			}
		case propose:
			epoch := msg.epoch
			value := msg.value
			if promisedEpoch.nil() || epoch.cmp(promisedEpoch) >= 0 {
				promisedEpoch = epoch
				acceptedValue, acceptedEpoch = value, epoch
				msg.replyTo <- accept{acceptorID: a.id, epoch: epoch}
			}
		}
	}
}
