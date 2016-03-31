package stego

import (
	"bytes"
	"os"
	"testing"
)

// TestEncode tests encoding data into an image.
func TestEncode(t *testing.T) {
	f, _ := os.Open("testdata/lena.png")
	defer f.Close()

	p := &Payload{Data: []byte("Hide me!"), Secret: ""}

	enc := new(bytes.Buffer)
	err := Encode(enc, f, p, nil)
	if err != nil {
		t.Error("TestEncode: expected no encoding error, got ", err)
	}
}

// TestDecode teset the decoding of data from an image.
func TestDecode(t *testing.T) {
	f, _ := os.Open("testdata/lena.png")
	defer f.Close()

	p := &Payload{Data: []byte("Hide me!"), Secret: ""}

	enc := new(bytes.Buffer)
	Encode(enc, f, p, nil)

	dec := new(bytes.Buffer)
	n, err := Decode(dec, enc, "")
	if err != nil {
		t.Error("TestDecode: expected no decoding error, got ", err)
	}
	if int(n) != len(p.Data) {
		t.Errorf("TestEncode: lenght of hidden data, want %v, got %v", len(p.Data), n)
	}
	if dec.String() != string(p.Data) {
		t.Errorf("TestDecode: hidden data, want %v, got %v", string(p.Data), dec.String())
		t.Errorf("TestDecode: hidden data, want %v, got %v", p.Data, dec.Bytes())
	}
}
