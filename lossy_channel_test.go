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
	"math/rand"
	"sync"
	"testing"
	"time"
)

type testMessage struct {
	number int
}

func TestThatLossyChannelDelivers1Message(t *testing.T) {
	l := lossyChannel{
		input:   make(chan message, 1),
		timeout: time.Second,
		size:    1,
		drop:    0,
	}

	msg := testMessage{number: rand.Int()}

	go func() {
		time.Sleep(100 * time.Millisecond)
		l.input <- msg
	}()

	got := (l.receive()).(testMessage)

	if got.number != msg.number {
		t.Errorf("got message %d, want %d", got.number, msg.number)
	}
}

func TestThatLossyChannelDelivers10Messages(t *testing.T) {
	l := lossyChannel{
		input:   make(chan message, 10),
		timeout: time.Millisecond,
		size:    10,
		drop:    0,
	}

	// send 10 messages
	for n := 0; n < 10; n++ {
		l.input <- testMessage{number: n}
	}

	received := make(map[int]bool) // received messages

	for n := 0; n < 10; n++ {
		msg := l.receive().(testMessage)
		received[msg.number] = true
	}

	// test all 10 messages were received
	for n := 0; n < 10; n++ {
		if _, ok := received[n]; !ok {
			t.Errorf("message %d was not received", n)
		}
	}
}

func TestThatLossyChannelDropsMessages(t *testing.T) {
	l := lossyChannel{
		input:   make(chan message, 1),
		timeout: time.Millisecond,
		size:    1,
		drop:    1, // always drop
	}

	l.input <- testMessage{}

	var mut sync.Mutex
	received := false

	go func() {
		l.receive()

		mut.Lock()
		defer mut.Unlock()

		received = true
	}()

	time.Sleep(1 * time.Second)

	mut.Lock()
	defer mut.Unlock()

	if received {
		t.Errorf("received message, wanted not to receive it")
	}
}
