package main

import (
	"bufio"
	"github.com/virtyx-technologies/sago/globals"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ########## Cipher suite information
//
const TLS13_OSSL_CIPHERS = "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:TLS_AES_128_CCM_SHA256:TLS_AES_128_CCM_8_SHA256"

type Cipher struct {
	hexc                 string
	TLS_CIPHER_OSSL_NAME string
	TLS_CIPHER_RFC_NAME  string
	TLS_CIPHER_SSLVERS   string
	TLS_CIPHER_KX        string
	TLS_CIPHER_AUTH      string
	TLS_CIPHER_ENC       string
	mac                  string
	TLS_CIPHER_EXPORT    string
}

var Ciphers []Cipher

func LoadCiphers() {
	filename := filepath.Join(globals.DataDir, "Cipher-mapping.txt")

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		words = append(words[:1], words[2:]...)        // delete words[1]
		c := Cipher{}
		c.hexc = words[0]
		c.TLS_CIPHER_OSSL_NAME = words[1]
		c.TLS_CIPHER_RFC_NAME  = words[2]
		c.TLS_CIPHER_SSLVERS   = words[3]
		c.TLS_CIPHER_KX        = words[4]
		c.TLS_CIPHER_AUTH      = words[5]
		c.TLS_CIPHER_ENC       = words[6]
		c.mac                  = words[7]
		if len(words) > 8 {
			c.TLS_CIPHER_EXPORT    = words[8]
		}

		Ciphers = append(Ciphers, c)
	}
}
