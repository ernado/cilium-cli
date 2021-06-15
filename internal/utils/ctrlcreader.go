// Copyright 2020-2021 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"io"
	"sync"
)

// CtrlCReader implements a simple Reader/Closer that returns Ctrl-C and EOF
// on Read() after it has been closed, and nothing before it.
type CtrlCReader struct {
	once sync.Once
	done chan struct{}
}

// NewCtrlCReader returns a new CtrlCReader instance
func NewCtrlCReader() *CtrlCReader {
	return &CtrlCReader{
		done: make(chan struct{}),
	}
}

// Read implements io.Reader.
// Blocks until we are done.
func (cc *CtrlCReader) Read(p []byte) (n int, err error) {
	select {
	case <-cc.done:
		if len(p) > 0 {
			p[0] = byte(3) // Ctrl-C
			if len(p) == 1 {
				return 1, io.EOF
			}
			// Add Ctrl-D for the case Ctrl-C alone is ineffective.
			// We skip this in the odd case where the buffer is too small
			p[1] = byte(4) // Ctrl-D
			return 2, io.EOF
		}
	}
	return 0, nil
}

// Close implements io.Closer. Note that we do not return an error on
// second close, not do we wait for the close to have any effect.
func (cc *CtrlCReader) Close() error {
	cc.once.Do(func() {
		close(cc.done)
	})
	return nil
}