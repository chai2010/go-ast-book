# Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# MnBook: Mini Markdown Book
# https://github.com/wa-lang/mnbook
#

default:
	mnbook serve

build:
	-rm book
	mnbook build
	-rm book/.gitignore
	-rm -rf book/.git

deploy:
	-@make clean
	mnbook build
	-rm book/.gitignore
	-rm -rf book/.git
	-rm -rf book/examples

	cd book && git init
	cd book && git add .
	cd book && git commit -m "first commit"
	cd book && git branch -M gh-pages
	cd book && git remote add origin git@github.com:chai2010/go-ast-book.git
	cd book && git push -f origin gh-pages

clean:
	-rm -rf book
