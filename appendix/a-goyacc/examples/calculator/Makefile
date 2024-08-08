# Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a Apache
# license that can be found in the LICENSE file.

run:
	make flex
	make goyacc

	@go fmt
	go run .

flex:
	flex --prefix=yy --header-file=calc.lex.h -o calc.lex.c calc.l

goyacc:
	goyacc -o calc.y.go -p "calc" calc.y

clean:
