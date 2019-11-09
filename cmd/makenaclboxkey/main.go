package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/nacl/box"
)

// KeyPair is a public and private key usable with NACL BOX
type KeyPair struct {
	PublicKey, PrivateKey string
}

func main() {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	keypair := KeyPair{
		PublicKey:  base64.StdEncoding.EncodeToString(pub[:]),
		PrivateKey: base64.StdEncoding.EncodeToString(priv[:]),
	}
	out, err := json.MarshalIndent(keypair, "", "    ")
	if err != nil {
		panic(err)
	}
	os.MkdirAll(os.Getenv("HOME")+"/.sops/naclbox", 0750)
	h := sha256.Sum256(pub[:])
	path := fmt.Sprintf("%s/.sops/naclbox/%x.key", os.Getenv("HOME"), h[:])
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		panic("file " + path + " already exists")
	}
	fmt.Printf("%s\n", out)
	err = ioutil.WriteFile(path, out, 0400)
	if err != nil {
		panic(err)
	}
	fmt.Println("keypair written to", path)
}
