package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Data structure for transaction
type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Signature []byte
}

// Structure for block view in blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}

// Structure of blockchain
type Blockchain struct {
	Blocks     []Block
	Difficulty int
}

// Generate keys for example
func generateKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey, &privateKey.PublicKey
}

// Create transaction with signature
func NewTransaction(sender, receiver string, amount float64, privateKey *rsa.PrivateKey) (*Transaction, error) {
	tx := &Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}
	signature, err := SignTransaction(tx, privateKey)
	if err != nil {
		return nil, err
	}
	tx.Signature = signature
	return tx, nil
}

// Signature  of transaction with private key
func SignTransaction(tx *Transaction, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%f", tx.Sender, tx.Receiver, tx.Amount)))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Check signature of transaction with using public key
func ValidateTransaction(tx *Transaction, publicKey *rsa.PublicKey) bool {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%f", tx.Sender, tx.Receiver, tx.Amount)))
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], tx.Signature)
	return err == nil
}

// Create new block
func (bc *Blockchain) NewBlock(transactions []Transaction, prevHash string) Block {
	block := Block{
		Index:        len(bc.Blocks),
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Nonce:        0,
	}
	block.Hash = bc.MineBlock(&block)
	return block
}

// Mining block with Proof-of-Work
func (bc *Blockchain) MineBlock(block *Block) string {
	var hash string
	target := strings.Repeat("0", bc.Difficulty)

	for {
		hash = CalculateHash(block)
		if strings.HasPrefix(hash, target) {
			break
		}
		block.Nonce++
	}
	return hash
}

// Calculate block hash
func CalculateHash(block *Block) string {
	data := fmt.Sprintf("%d%s%v%s%d", block.Index, block.Timestamp, block.Transactions, block.PrevHash, block.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Initialisation of new blockchain with genesis block
func NewBlockchain() *Blockchain {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{},
		PrevHash:     "0",
		Hash:         "",
		Nonce:        0,
	}
	genesisBlock.Hash = CalculateHash(&genesisBlock)
	return &Blockchain{
		Blocks:     []Block{genesisBlock},
		Difficulty: 4, // initial difficulty
	}
}

// Adding a block to the blockchain
func (bc *Blockchain) AddBlock(block Block) {
	if bc.ValidateBlock(block, bc.Blocks[len(bc.Blocks)-1]) {
		bc.Blocks = append(bc.Blocks, block)
	}
}

// Block validation: checking its hash and whether it matches the previous block
func (bc *Blockchain) ValidateBlock(block, prevBlock Block) bool {
	if block.PrevHash != prevBlock.Hash {
		return false
	}
	if block.Hash != CalculateHash(&block) {
		return false
	}
	return true
}

// Blockchain integrity check
func (bc *Blockchain) IsChainValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if !bc.ValidateBlock(currentBlock, prevBlock) {
			return false
		}
	}
	return true
}

// The main function of the application
func main() {
	// Initializing a new blockchain
	blockchain := NewBlockchain()

	// Generating keys for two participants
	privateKey1, publicKey1 := generateKeys()
	privateKey2, publicKey2 := generateKeys()

	// Create a transaction between two users
	tx1, err := NewTransaction("Natasha", "Ashot", 100, privateKey1)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	// Checking the validity of the transaction
	if !ValidateTransaction(tx1, publicKey1) {
		fmt.Println("Transaction is invalid!")
		return
	}

	// Create a block with transactions
	block1 := blockchain.NewBlock([]Transaction{*tx1}, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)

	// Adding a block to the blockchain
	blockchain.AddBlock(block1)

	fmt.Println("Block 1 added!")
	printBlock(block1)

	// Create another transaction between the same users
	tx2, err := NewTransaction("Ashot", "Natasha", 50, privateKey2)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	// Checking the validity of the second transaction
	if !ValidateTransaction(tx2, publicKey2) {
		fmt.Println("The second transaction is invalid!")
		return
	}

	// Create a new block with a new transaction
	block2 := blockchain.NewBlock([]Transaction{*tx2}, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)

	// Adding the second block to the blockchain
	blockchain.AddBlock(block2)

	fmt.Println("Block 2 added!")
	printBlock(block2)

	// Blockchain Integrity Check
	if blockchain.IsChainValid() {
		fmt.Println("The block chain is valid!")
	} else {
		fmt.Println("Error: Blockchain is invalid!")
	}
}

// Function for outputting information about a block
func printBlock(block Block) {
	fmt.Printf("Index: %d\n", block.Index)
	fmt.Printf("Time: %s\n", block.Timestamp)
	fmt.Printf("Transaction: %v\n", block.Transactions)
	fmt.Printf("Previous hash: %s\n", block.PrevHash)
	fmt.Printf("Hash: %s\n", block.Hash)
	fmt.Printf("Nonce: %d\n\n", block.Nonce)
}
