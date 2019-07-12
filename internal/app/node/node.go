package node

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strings"
	"time"


	"dohrnii/internal/app/block"
	"github.com/davecgh/go-spew/spew"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"

)

var currentBc block.Blockchain

// Host creates a new basic host
func Host(listenPort int, secio bool, randseed int64) (host.Host, error) {
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Creates a new RSA key pair for this host.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil
}

// HandleStream handles incoming stream from new peers
func HandleStream(s net.Stream) {
	log.Println("Got a new stream!")

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)
}

// ReadData handles the buffer stream to read peers data
func ReadData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			bcReceived := []block.Block{}
			s := string(str)
			json.Unmarshal([]byte(s), &bcReceived)

			if err := json.Unmarshal([]byte(str), &bcReceived); err != nil {
				log.Fatal(err)
			}

			block.Mutex.Lock()
			if bcReceived[len(bcReceived) - 1].Height > len(block.Bc.Chain) {
				block.Bc.Chain = bcReceived
				bytes, err := json.MarshalIndent(bcReceived, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			block.Mutex.Unlock()
		}
	}
}

// WriteData write into others buffer stream 
func WriteData(rw *bufio.ReadWriter) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			block.Mutex.Lock()
			bytes, err := json.Marshal(block.Bc.Chain)
			if err != nil {
				log.Println(err)
			}
			block.Mutex.Unlock()

			block.Mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			block.Mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		if err != nil {
			log.Fatal(err)
		}
		newBlock := block.CreateBlock(currentBc.GetLastBlock())

		block.Mutex.Lock()
		block.Bc.Chain = append(block.Bc.Chain, newBlock)
		block.Mutex.Unlock()

		bytes, err := json.Marshal(block.Bc.Chain)
		if err != nil {
			log.Println(err)
		}

		spew.Dump(block.Bc.Chain)

		block.Mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		block.Mutex.Unlock()
	}

}