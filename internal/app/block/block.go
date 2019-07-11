package block

import (
	"encoding/json"
	"sync"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"dohrnii/internal/app/twitter"
	"time"
	"fmt"
	"log"
)

// Block represents what is inside a block
type Block struct {
	Height 		int				`json:"height"`
	PrevHash 	string			`json:"prevhash"`
	Hash 		string			`json:"hash"`
	Tweets		[]twitter.Tweet	`json:"tweets"`
	Nonce		int				`json:"nonce"`
	Timestamp 	string			`json:"timestamp"`
	Difficulty 	float64			`json:"difficulty"`
}

// Blockchain represents the blockchain
type Blockchain struct {
	Chain 		[]Block		`json:"chain"`
	Nodes		[]string	`json:"nodes"`
}

var targetMax = "0x0000FFFFFFFFFFFFFFFFFFFFFFF000000000000000000000000000000000000"
var tweetPool []twitter.Tweet

// Mutex securely prevent race conditions
var Mutex = &sync.Mutex{}

// Bc is the blockchain of the current instance of the program
var Bc Blockchain

// Initialize the values to get the blockchain starting
func Initialize() {
	fmt.Println("Wallou")
	go fetchTweets()
	currentTime, _ := time.Parse("2006-01-02T15:04:05", time.Now().Format("2006-01-02T15:04:05"))
	genesis := Block{0, "", "", tweetPool, 0, currentTime.String(), 1}
	Bc.Chain = append(Bc.Chain, genesis)
	blockchainToJSON()
	newBlock := CreateBlock(genesis)
	Bc.Chain = append(Bc.Chain, newBlock)
	for {
		newBlock := CreateBlock(GetLastBlock(Bc.Chain))
		Bc.Chain = append(Bc.Chain, newBlock)
		blockchainToJSON()
	}
}

func fetchTweets() {
	for {
		tweets := twitter.GetTweets()
		for i := 0; i < len(tweets); i++ {
			tweetPool = append(tweetPool, tweets[i])
		}
		time.Sleep(60 * time.Second)
	}
}

// CalculateDifficulty calculates the difficulty
func (b Block) calculateDifficulty(nDifficulty *float64) float64 {
	return *nDifficulty
}

func proofOfWork(block Block) string {
	t := time.Now()
	block.Timestamp = t.String()
	data := string(block.Height) + block.PrevHash + block.Timestamp
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// CalculateBlockHash calculate the block hash
func (b Block) calculateBlockHash(block Block) string {
	var hash string
	var hashInt *big.Int

	hash = proofOfWork(block)
	target, err := new(big.Int).SetString(targetMax, 0)
	hashInt, err = new(big.Int).SetString(hash, 16)

	if !err {
		log.Println("Error on big Int")
	}

	for hashInt.Cmp(target) != -1 {
		hash = proofOfWork(block)
		hashInt, _ = new(big.Int).SetString(hash, 16)
		//fmt.Printf("\r%s", hash)
		if hashInt.Cmp(target) == -1 {
			break
		}
	}
	fmt.Println()
	return hash
}

// CreateBlock intialize a block
func CreateBlock(lastBlock Block) Block {
	var newBlock Block

	newBlock.Tweets = tweetPool
	newBlock.Height = lastBlock.Height + 1
	newBlock.PrevHash = lastBlock.Hash
	newBlock.Hash = newBlock.calculateBlockHash(newBlock)
	currentTime, _ := time.Parse("2006-01-02T15:04:05", time.Now().Format("2006-01-02T15:04:05"))
	newBlock.Timestamp = currentTime.String()
	return newBlock
}

// GetLastBlock fetch the latest block of the current blockchain
func GetLastBlock(b []Block) Block {
	index := len(b) - 1
	return b[index]
}

func blockchainToJSON() {
	bytes, err := json.MarshalIndent(GetLastBlock(Bc.Chain), "", "  ")
	if err != nil {

		log.Fatal(err)
	}
	fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
}