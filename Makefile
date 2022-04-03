default:
	mdbook serve

build:
	-rm book
	mdbook build
	-rm book/.gitignore
	-rm -rf book/.git

clean: