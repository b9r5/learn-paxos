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
	"time"
)

// lossyChannel is a chan-like struct that is allowed to rearrange and drop
// messages.
type lossyChannel struct {
	// input channel for receiving messages that get buffered
	input chan message

	// buffer of messages
	buf []message

	// amount of time to wait for buffer to fill before returning a message from
	// receive
	timeout time.Duration

	// target number of messages to buffer before returning one from receive
	size int

	// probability of dropping a message in the range [0, 1)
	drop float64

	// output channel onto which non-dropped, possibly reordered messages get
	// placed
	output chan message
}

// newLossyChannel returns a new lossyChannel with the given parameters.
// size is the number of messages to buffer before returning one from receive.
// timeout is the amount of time to wait for the lossy channel to contain
// size messages before returning a message. drop is the probability in the
// range [0, 1) of dropping a message.
func newLossyChannel(size int, timeout time.Duration, drop float64) *lossyChannel {
	l := &lossyChannel{
		input:   make(chan message, size),
		buf:     make([]message, 0, size),
		timeout: timeout,
		size:    size,
		drop:    drop,
		output:  make(chan message, size),
	}

	go l.run()

	return l
}

// run repeatedly invokes receive and places the result on l.output.
// if the result is nil (as would happen if close is invoked), then run returns.
func (l *lossyChannel) run() {
	for {
		msg := l.receive()
		if msg == nil {
			return
		}
		l.output <- msg
	}
}

// receive returns a message from l.buf. It waits up to l.timeout for the
// channel to contain l.size messages. Once either the channel contains l.size
// messages or the timeout expires: then if the channel contains one or more
// messages, it returns one of them selected pseudo-randomly; else (if the
// channel contains zero messages), it waits to receive a message, and returns
// that message. It drops incoming messages with probability l.drop.
func (l *lossyChannel) receive() message {
	start := time.Now()

	for {
		remaining := l.timeout - time.Since(start) // time left

		// if there are at least l.size messages in the buffer, or there is no
		// time left and we have at least 1 message in the buffer, then permute
		// l.buf and return one of its messages
		if len(l.buf) >= l.size || remaining <= 0 && len(l.buf) > 0 {
			rand.Shuffle(len(l.buf), func(i, j int) {
				l.buf[i], l.buf[j] = l.buf[j], l.buf[i]
			})

			msg := l.buf[0]

			buf := make([]message, len(l.buf)-1)
			copy(buf, l.buf[1:])
			l.buf = buf

			return msg
		}

		select {
		case msg := <-l.input:
			if rand.Float64() >= l.drop {
				l.buf = append(l.buf, msg) // yay! buffer the message
			}
		case <-time.After(remaining):
			if len(l.buf) <= 0 {
				// time's up! return the next message that's not dropped
				for {
					msg := <-l.input
					if rand.Float64() >= l.drop {
						return msg
					}
				}
			} // else resume at beginning of for loop
		}
	}
}

// close closes l's input and output channels.
func (l *lossyChannel) close() {
	close(l.input)
	close(l.output)
}
