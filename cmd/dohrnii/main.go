package main

import (
	"dohrnii/internal/app/handlers"
	"path/filepath"
	"path"
	"github.com/gin-gonic/gin"
	"context"
	"bufio"
	"fmt"
	"flag"
	"log"
	"dohrnii/internal/app/block"
	"dohrnii/internal/app/node"
	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

func server() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui" + path.Join(dir, file))
		}
	})

	authorized := r.Group("/")
	authorized.GET("/blockchain", handlers.GetBlockchain)
	authorized.GET("/lastblock", handlers.GetLastBlock)
	errServer := r.Run(":3000")
	if errServer != nil {
		panic(errServer)
	}
}

func main() {
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	ha, err := node.Host(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

	go server()
	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", node.HandleStream)

		go block.Initialize()
		select {}
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", node.HandleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")

		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go node.WriteData(rw)
		go node.ReadData(rw)

		go block.Initialize()

		select {}
	}
}