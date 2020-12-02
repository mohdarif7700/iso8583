// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package fuzzreader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCorpusSymlinks(t *testing.T) {
	// avoid symbolic link error on windows
	if runtime.GOOS == "windows" {
		t.Skip()
	}
	fds, err := ioutil.ReadDir("corpus")
	if err != nil {
		t.Fatal(err)
	}
	if len(fds) == 0 {
		t.Fatal("no file descriptors found in corpus/")
	}

	for i := range fds {
		if fds[i].Mode()&os.ModeSymlink != 0 {
			if path, err := os.Readlink(filepath.Join("corpus", fds[i].Name())); err != nil {
				t.Errorf("broken symlink: %v", err)
			} else {
				if _, err := os.Stat(filepath.Join("corpus", path)); err != nil {
					t.Errorf("broken symlink: %v", err)
				}
			}
		} else {
			t.Errorf("%s isn't a symlink, move outside corpus/ and symlink into directory", fds[i].Name())
		}
	}
}

func TestFuzzWithValidData(t *testing.T) {
	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "financial_transaction_message.dat"))
	if err != nil {
		t.Fatal(err)
	}

	if ret := Fuzz(byteData); ret != 1 {
		t.Errorf("Expected value is 1 (got %v)", ret)
	}

	byteData, err = ioutil.ReadFile(filepath.Join(basePath, "..", "testdata", "iso_reversal_message_error_date.dat"))
	if err != nil {
		t.Fatal(err)
	}

	if ret := Fuzz(byteData); ret != 0 {
		t.Errorf("Expected value is 0 (got %v)", ret)
	}

	if ret := Fuzz(byteData); ret != 0 {
		t.Errorf("Expected value is 0 (got %v)", ret)
	}
}
