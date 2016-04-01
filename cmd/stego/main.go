package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/toashd/stego"
)

var (
	encode   = flag.Bool("e", false, "Encode a message into the file")
	decode   = flag.Bool("d", false, "Decode a message from the file")
	message  = flag.String("m", "", "Message to be encoded")
	image    = flag.String("p", "", "Image file name to encode or decode")
	output   = flag.String("o", "", "Output file name")
	password = flag.String("pwd", "", "Password to encrypt the message")
)

var usage = `Usage: stego [options...] <values>
Options:
  -d    Decode a message from the file
  -e    Encode a message into the file
  -m    Message to be encoded
  -p    Image file name to encode or decode
  -o    Output file name
  -pwd  Password to encrypt the message

Examples:
  $ stego -e -p lena.png -m "Lena is beautiful." -o out
  $ stego -d -p out.png
    Lena is beautiful
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()

	if *encode {
		if *message == "" || *image == "" {
			exitWithUsage("Error: need to specify -p and -m")
		}

		f, err := os.Open(*image)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}

		if len(*message)*8 > int(fi.Size())-54 {
			err := errors.New("source image not large enough to hold the message")
			fmt.Fprintf(os.Stderr, err.Error(), "\n")
			os.Exit(1)
		}

		if *output == "" {
			base := filepath.Base(f.Name())
			*output = strings.TrimSuffix(base, filepath.Ext(base)) + "-enc"
		}

		o, _ := os.Create(*output + ".png")
		defer o.Close()

		w := bufio.NewWriter(o)
		err = stego.Encode(w, f, &stego.Payload{Data: []byte(*message), Secret: *password}, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}
		w.Flush()

		os.Exit(0)
	}
	if *decode {
		if *image == "" {
			exitWithUsage("Error: need to specidy \"-p\"")
		}

		f, err := os.Open(*image)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}
		defer f.Close()

		w := bufio.NewWriter(os.Stdout)
		_, err = stego.Decode(w, f, *password)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}
		w.Flush()
		os.Exit(0)
	}
	if *image == "" {
		exitWithUsage("Error: please specify -p")
	}
	printCap(*image)
}

// exitWithUsage prints usage information and terminates.
func exitWithUsage(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message+"\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

// printCap prints the number of characters the specifed file can hold (capacity).
func printCap(filename string) {
	image, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening " + filename)
	}
	fmt.Printf("%s can hold %v characters\n", filename, (len(image)/8)-54)
}
