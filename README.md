Nexelliad
========

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](https://choosealicense.com/licenses/isc/)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/Nexell-AI-Network/nexelliad/v2/)

Nexelliad is the reference full node Nexell-AI implementation written in Go (golang).

## What is Nexell-AI

Nexell-AI is a fork of Kaspa with an concept of Distributed Artificial Intelligence Network on BlockDAG. 
Kaspa is an attempt at a proof-of-work cryptocurrency with instant confirmations and sub-second block times. 
It is based on [the PHANTOM protocol](https://eprint.iacr.org/2018/104.pdf), a generalization of Nakamoto consensus.

## Requirements

Go 1.19 or later.

## Installation

#### Build from Source

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Ensure Go was installed properly and is a supported version:

```bash
$ go version
```

- Run the following commands to obtain and install nexelliad including all dependencies:

```bash
$ git clone https://github.com/Nexell-AI-Network/nexelliad/v2/
$ cd nexelliad
$ go install . ./cmd/...
```

- Nexelliad (and utilities) should now be installed in `$(go env GOPATH)/bin`. If you did
  not already add the bin directory to your system path during Go installation,
  you are encouraged to do so now.
  
- Open your shell configuration file. For example, for Bash, you can use the following command:
  
```bash
$ nano ~/.bashrc
```
- Add the following line to the end of the file:

```bash
 export PATH=$PATH:$(go env GOPATH)/bin
```

## Getting Started

Nexelliad has several configuration options available to tweak how it runs, but all
of the basic operations work with zero configuration.

## Creating a wallet

- To create a wallet, you need to run nexelliad with utxoindex

```bash
$ nexelliad --utxoindex
```
- Open another terminal

```bash
$ nexelliawallet create
```

- You will be asked to choose a password for the wallet (a password must be at least 8 characters long, and it won't be shown on the screen you as you entering it). After that you should run this command in order to start the wallet daemon:

```bash
$ nexelliawallet start-daemon
```
- Do not close the first 2 terminals and open a new terminal and then run this in order to request an address from the wallet:

```bash
$ nexelliawallet new-address
```

- Your screen will show you something like this:

The wallet address is:
nexellia:0123456789abcdef0123456789abcdef0123456789

- To see your secret seed phrase :

```bash
$ nexelliawallet dump-unencrypted-data
```

Note: Every time you ask nexelliawallet for an address you will get a different address. This is perfectly fine. Every secret key is associated with many different public addresses and there is no reason not to use a fresh one for each transaction.

At this point your can close the wallet daemon, though you should keep it running of you want to be able to check your balance and make transactions


## Discord
Join our discord server using the following link: https://discord.gg/MHn4wStC4h


## License

Nexelliad is licensed under the copyfree [ISC License](https://choosealicense.com/licenses/isc/).
