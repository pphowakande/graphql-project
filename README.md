# The Go Programming Language

### Download and Install

#### Binary Distributions

Official binary distributions are available at https://golang.org/dl/.

After downloading a binary release, visit https://golang.org/doc/install
or load [doc/install.html](./doc/install.html) in your web browser for installation
instructions.

#### Install From Source

If a binary distribution is not available for your combination of
operating system and architecture, visit
https://golang.org/doc/install/source or load [doc/install-source.html](./doc/install-source.html)
in your web browser for source installation instructions.

### Contributing

Go is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines:
	https://golang.org/doc/contribute.html

Note that the Go project uses the issue tracker for bug reports and
proposals only. See https://golang.org/wiki/Questions for a list of
places to ask questions about the Go language.

[rf]: https://reneefrench.blogspot.com/
[cc3-by]: https://creativecommons.org/licenses/by/3.0/

# ATHLIMA

## Installation
To install the packages you need to the cmd folder and run the following command

`go get ./...`

After installing all the packages update the configuration file config.yaml with database and email configurations

## Run

To build the service execute the following command
`go build -o turf-api *.go`

To run the service execute the following command
`sudo ./turf-api`


## Run Test Cases

To run the testcases execute the following command
`go test -v ./...`