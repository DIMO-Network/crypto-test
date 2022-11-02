package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

		fmt.Printf("Hash: %s\n", h)

		cmd := exec.Command("autopi", "crypto.sign_string", strings.TrimPrefix(h.Hex(), "0x"))
		o, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		m := reg.FindSubmatch(o)

		sig := common.FromHex(string(m[1]))
		fmt.Printf("Signature: %s\n", hexutil.Encode(sig))

		func() (err error) {
			go func() {
				if err != nil {
					fmt.Printf("Error: %s\n", err)
				}
			}()

			sig := common.FromHex(string(m[1]))

			if len(sig) != 65 {
				return fmt.Errorf("signature has length %d", len(sig))
			}

			sig[64] -= 27

			uncPubKey, err := crypto.Ecrecover(h[:], sig)
			if err != nil {
				return fmt.Errorf("failed to recover: %w", err)
			}

			pubKey, err := crypto.UnmarshalPubkey(uncPubKey)
			if err != nil {
				return fmt.Errorf("failed to unmarshal public key: %w", err)
			}

			addr := crypto.PubkeyToAddress(*pubKey)

			fmt.Printf("Recovered: %s\n", addr)
			return nil
		}()
	}
}
