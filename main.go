package main

import (
	"fmt"
	"strings"

	"github.com/dsonic0912/PolicyReporter-FSM/examples"
)

func main() {
	fmt.Println("=== Mod-Three Finite State Machine ===")
	fmt.Println("Using Generic FSM Library")
	fmt.Println()

	// Example from specification
	fmt.Println("Example 1: Input = \"110\"")
	examples.PrintModThreeTrace("110")
	fmt.Println()

	// Show the automaton configuration
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	fmt.Println("Automaton Configuration:")
	fmt.Println()
	fa := examples.NewModThreeAutomaton()
	fmt.Print(fa.String())
	fmt.Println()

	// Additional examples
	testInputs := []string{"101010", "1111", "1001", "0", "1"}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	fmt.Println("Additional Examples:")
	fmt.Println()

	for _, input := range testInputs {
		result, err := examples.ModThree(input)
		if err != nil {
			fmt.Printf("Input: %s - Error: %v\n", input, err)
			continue
		}
		fmt.Printf("Input: %s -> Output: %d (remainder when %s₂ is divided by 3)\n", input, result, input)
	}
}
