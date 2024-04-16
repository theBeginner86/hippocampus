//  Copyright 2024 Pranav Singh

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package aof

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/thebeginner86/hippocampus/resp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	// opens file if exists in the path or creats a new one with specified permissions
	fileH, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: fileH,
		rd:   bufio.NewReader(fileH),
	}

	// uses a go routine and an infinite loop
	// that syncs the changes to disk every second.
	// reason being that the executed commands to be in-sync
	// with max reliability onto the disk. If syncing every sec
	// is skiiped perform then it would be upto the OS to commit
	// the changes and could lead to data miss.
	//
	// use of go routine and syncing every second could be replaced
	// with syncing only when a command is executed. This has pros to
	// the extent that system won't be syncing every second but would
	// lead to poor performance and I/O operations are expensive and
	// hence would reduce scalability of DB.
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}

	}()
	return aof, nil
}

// Close() ensures that file is properly closed when system is shutdown
// here use of mutex locks ensure that file is prevented from
// concurrent access
func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

// Write() ensures that the command is persisted into file in the exact
// same RESP format that we recieve. As later on system reboot these would
// be re-read and bring back system to its previous state, that is, the one
// before failover
func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

// Read() ensures that the commands persisted in file are read
// and then sent to func sent as args. This func should have the logic
// to handle the specific RESP commands. It should make use of the handler map
// that is tracks mapping between command and handling logic
func (aof *Aof) Read(fn func(value resp.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	// brings the file ptr to the starting of the file
	aof.file.Seek(0, io.SeekStart)
	reader := resp.NewResp(aof.file)

	for {
		// reads each RESP command until EOF is found
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// after reading sends it to the handler func passed as arg
		fn(value)
	}

	return nil
}
