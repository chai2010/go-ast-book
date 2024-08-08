// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	br := bufio.NewReader(os.Stdin)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			// ^D: Input end of file on Unix/Linux
			// ^Z: Input end of file on Windows
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}

		// quit
		if s := strings.TrimSpace(line); s == "q" || s == "quit" || s == "exit" {
			return
		}

		// flex + goyacc
		calcParse(newCalcLexer([]byte(line)))
	}
}
