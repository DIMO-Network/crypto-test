package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var reg = regexp.MustCompile(`^.+(0x[0-9a-z]+)`)

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 20; i++ {
		var h common.Hash
		for i := range h {
			h[i] = byte(rand.Uint32())
		}

		fmt.Printf("Hash: %s\n", h)

		cmd := exec.Command("sudo", "/home/pi/edge-identity/edge-identity", "--lib", "/usr/lib/softhsm/libsofthsm2.so", "sign", "--token", "dimo", "--label", "clitest", "--hash", h.Hex(), "--pin", "1234")
		o, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Command: %s\n", cmd.String())
		fmt.Printf("Output: %s\n", o)
		m := reg.FindSubmatch(o)
		sig := common.FromHex(string(m[1]))
		fmt.Printf("Signature: %s\n", hexutil.Encode(sig))

		func() (err error) {
			defer func() {
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
