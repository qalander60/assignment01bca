package assignment01bca

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

// WriteLastHash writes the hash string to a file named "last_hash.txt"
func WriteLastHash(hash string) error {
	return os.WriteFile("last_hash.txt", []byte(hash), 0644)
}

// ReadLastHash reads the hash string from the file named "last_hash.txt"
func ReadLastHash() (string, error) {
	data, err := os.ReadFile("last_hash.txt")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// new block structure
type Block struct {
	transaction  string
	nonce        int
	previousHash string
	currentHash  string
	id           int
	next         *Block
}

// Getter methods
func (b *Block) GetTransaction() string {
	return b.transaction
}

func (b *Block) GetNonce() int {
	return b.nonce
}

func (b *Block) GetPreviousHash() string {
	return b.previousHash
}

func (b *Block) GetCurrentHash() string {
	return b.currentHash
}

func (b *Block) GetID() int {
	return b.id
}

// Hash function that serializes the block, calculates the hash, and returns it as a string
func CalculateHash(stringToHash string) string {
	hash := sha256.Sum256([]byte(stringToHash))
	return fmt.Sprintf("%x", hash[:])
}

type Chain struct {
	head   *Block
	length int
}

// ProofOfWork adjusts the nonce until the hash meets the difficulty criteria
func (b *Block) ProofOfWork(difficulty int) {
	prefix := strings.Repeat("a", difficulty)
	for {
		blockData := fmt.Sprintf("%s%d%s%d", b.transaction, b.nonce, b.previousHash, b.id)
		b.currentHash = CalculateHash(blockData)
		if strings.HasPrefix(b.currentHash, prefix) {
			break
		}
		b.nonce++
	}
}

// function to add new block to the blockchain
func (bc *Chain) NewBlock(transaction string) *Block {
	//define the new block
	id := bc.length

	newBlock := &Block{transaction, 0, "", "", id, nil}
	//checking if teh head is empty then pasted the genesis block there
	if bc.head == nil {
		bc.head = newBlock
		// Serialize block data for hashing
		newBlock.ProofOfWork(2)
		fmt.Println("Head block")
		fmt.Println()
		bc.length++
		// Write the current hash of the new block to the file
		err := WriteLastHash(newBlock.currentHash)
		if err != nil {
			fmt.Println("Error writing last hash to file:", err)
		}
		return newBlock
	}
	currBlock := bc.head
	//iterate the loop until the blockchain ends
	for currBlock.next != nil {
		currBlock = currBlock.next
	}
	//add the new block to the end of the blockchain
	newBlock.previousHash = currBlock.currentHash
	//proof of work for the new block addition
	newBlock.ProofOfWork(2)
	currBlock.next = newBlock
	bc.length++
	// Write the current hash of the new block to the file
	err := WriteLastHash(newBlock.currentHash)
	if err != nil {
		fmt.Println("Error writing last hash to file:", err)
	}
	return newBlock
}

// function to list all blocks
func (bc *Chain) ListBlocks() {
	if bc.head == nil {
		fmt.Println("Blockchain is empty.")
		return
	}
	currBlock := bc.head
	for currBlock != nil {
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Block ID       : %d\n", currBlock.id)
		fmt.Printf("Transaction    : %s\n", currBlock.transaction)
		fmt.Printf("Nonce          : %d\n", currBlock.nonce)
		fmt.Printf("Previous Hash  : %s\n", currBlock.previousHash)
		fmt.Printf("Current Hash   : %s\n", currBlock.currentHash)
		fmt.Println("--------------------------------------------------")
		fmt.Println()
		currBlock = currBlock.next
	}
}

// function to change the block transaction of given block reference without varying other blocks
func (bc *Chain) ChangeBlockUn(id int, transaction string) {
	if bc.head == nil {
		return
	}
	currBlock := bc.head
	for currBlock != nil {
		if currBlock.id == id {
			currBlock.transaction = transaction
			currBlock.ProofOfWork(2)
			return
		}
		currBlock = currBlock.next
	}
	fmt.Printf("Block with the ID : %d not found", id)
}

// function to change the block transaction of given block reference
func (bc *Chain) ChangeBlock(id int, transaction string) {
	if bc.head == nil {
		return
	}
	currBlock := bc.head
	for currBlock != nil {
		if currBlock.id == id {
			currBlock.transaction = transaction
			currBlack.ProofOfWork(2)
			nextBlock := currBlock.next
			for nextBlock != nil {
				nextBlock.previousHash = currBlock.currentHash
				nextBlock.ProofOfWork(2)
				currBlock = nextBlock
				nextBlock = nextBlock.next
			}
			// Write the current hash of the new block to the file
			err := WriteLastHash(currBlock.currentHash)
			if err != nil {
				fmt.Println("Error writing last hash to file:", err)
			}
			return
		}
		currBlock = currBlock.next
	}
	fmt.Printf("Block with the ID : %d not found", id)
}

// function to verify blockchain in case any changes are made
func (bc *Chain) VerifyChain() {
	fmt.Println("***Intitiating blockchian verification***")
	if bc.head == nil {
		return
	}
	currBlock := bc.head
	nextBlock := currBlock.next
	for nextBlock != nil {
		currBlock.ProofOfWork(2)

		if nextBlock.previousHash != currBlock.currentHash {
			fmt.Printf("previous block hash value stored in %d doesn't match the actual hash value of %d\n", nextBlock.id, currBlock.id)
		}
		nextBlock.previousHash = currBlock.currentHash
		currBlock = nextBlock
		nextBlock = nextBlock.next
	}
	currBlock.ProofOfWork(2)
	fmt.Println("Comparing the last node hash with the stored one in teh file")
	lastHash, err := ReadLastHash()
	if err == nil {
		if lastHash == currBlock.currentHash {
			fmt.Println("[+] The last hash value matched with the value stored in the file.")
		} else {
			fmt.Println("[-] The last hash value doesn't matched with the value stored in the file.")
		}
	}
	fmt.Println("Blockchain verification complete.")
}
