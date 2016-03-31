# Stego
Stego is a very straightforward implementation of lsb image steganography. It encodes and decodes data into an image. The data or message is hidden in the LSBs of the Red, Green and Blue image components. The package provides both, encoder and decoder, as well as an additional command line tool to easily hide data within and extract hidden data from image files.

## Installation

```bash
$ go get github.com/toashd/stego
```

## Usage

The command line tool source provides the basic usage pattern of stego.

```bash
  $ stego -e -p lena.png -m "Lena is beautiful" -o out
  $ stego -d -p out.png
    Lena is beautiful
`
```

API overview:

```go
func Encode(w io.Writer, r io.Reader, p *Payload, o *Options) error
func Decode(w io.Writer, r io.Reader, pwd string) (int64, error)
```

## Todo

* Support encoding/decoding of audio and video files
* Fix issues with encoding/decoding go jpeg and gif files
* Allow encoding of remote images (e.g. url flag)
* Add web handler to easily plug into web applications

## Related Work

Tools
* [Steghide](http://steghide.sourceforge.net/) is a steganography program that is able to hide data in various kinds of image- and audio-files. The color- respectivly sample-frequencies are not changed thus making the embedding resistant against first-order statistical tests.
* [Stegano.js](https://github.com/tuseroni/stegano.js) - steganographic encoder and decoder for javascript

Papers
* [Defeating Statistical Steganalysis](http://www.citi.umich.edu/u/provos/stego/)
* [Detecting Steganographic Content on the Internet](https://www.citi.umich.edu/techreports/reports/citi-tr-01-11.pdf)


## Contribution

Please feel free to suggest any kind of improvements and/or refactorings - just file an
issue or fork and submit a pull requests.

## License

Stego is available under the MIT license. See the LICENSE file for more info.

