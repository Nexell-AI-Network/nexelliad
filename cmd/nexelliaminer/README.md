# nexelliaminer

nexelliaminer is a CPU-based miner for nexelliad

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
$ cd nexelliad/cmd/nexelliaminer
$ go install .
```

- Kapaminer should now be installed in `$(go env GOPATH)/bin`. If you did
  not already add the bin directory to your system path during Go installation,
  you are encouraged to do so now.
  
## Usage

The full nexelliaminer configuration options can be seen with:

```bash
$ nexelliaminer --help
```

But the minimum configuration needed to run it is:
```bash
$ nexelliaminer --miningaddr=<YOUR_MINING_ADDRESS>
```