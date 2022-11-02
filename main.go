package main

import (
	"log"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var reg = regexp.MustCompile(`value: '0x([0-9a-z]+)'`)

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 20; i++ {
		var h common.Hash
		for i := range h {
			h[i] = byte(rand.Uint32())
		}

		cmd := exec.Command("autopi", "crypto.sign_string", strings.TrimPrefix(h.Hex(), "0x"))
		o, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		m := reg.FindSubmatch(o)
		sig := common.FromHex(string(m[1]))

		sig[64] -= 27

		uncPubKey, err := crypto.Ecrecover(h[:], sig)
		if err != nil {
			log.Printf("Failed to recover: %s", err)
		}

		pubKey, err := crypto.UnmarshalPubkey(uncPubKey)
		if err != nil {
			log.Printf("Couldn't load public key: %s", err)
		}

		addr := crypto.PubkeyToAddress(*pubKey)

		log.Printf("Got: %s", addr)
	}
}
