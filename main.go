package main

import (
	"fmt"

	"github.com/qalander60/bcaassignment01"
)

func main() {
	bc := &bcaassignment01.Chain{}

	// Define the difficulty for Proof-of-Work
	// Add blocks to the blockchain
	bc.NewBlock("Genesis block transaction")
	bc.NewBlock("Alice to Bob")
	bc.NewBlock("Bob to Charlie")

	// List all blocks
	fmt.Println("\nInitial Blockchain:")
	bc.ListBlocks()

	// Change a block's transaction
	bc.ChangeBlock(1, "Alice to Alice")
	// Verify the blockchain
}
