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

type message interface{}

// prepare is the message sent by the proposer in phase 1 of Classic Paxos.
type prepare struct {
	epoch      Epoch
	proposerID int
	replyTo    chan<- message
}

// String returns the string form of a prepare message.
func (p prepare) String() string {
	return fmt.Sprintf("prepare(%s) from proposer %d",
		p.epoch, p.proposerID)
}

// promise is the message sent by the acceptor in phase 1 of Classic Paxos.
type promise struct {
	epoch, acceptedEpoch Epoch
	acceptedValue        string
	acceptorID           int
}

// String returns the string form of a promise message.
func (p promise) String() string {
	var v string
	if p.acceptedValue == "" {
		v = "nil"
	} else {
		v = p.acceptedValue
	}

	return fmt.Sprintf("promise(%s, %s, %s) from acceptor %d",
		p.epoch, p.acceptedEpoch, v, p.acceptorID)
}

// propose is the message sent by the proposer in phase 2 of Classic Paxos.
type propose struct {
	epoch      Epoch
	value      string
	proposerID int
	replyTo    chan<- message
}

// String returns the string form of a propose message.
func (p propose) String() string {
	return fmt.Sprintf("propose(%s, %s) from proposer %d",
		p.epoch, p.value, p.proposerID)
}

// accept is the message sent by the acceptor in phase 2 of Classic Paxos.
type accept struct {
	epoch      Epoch
	acceptorID int
}

// String returns the string form of an accept message.
func (a accept) String() string {
	return fmt.Sprintf("accept(%s) from acceptor %d",
		a.epoch, a.acceptorID)
}
