//go:build ignore

package main

import (
	"fmt"
)

//go:generate

// Calculator represents a simple calculator
type Calculator struct {
	// result holds the current calculation result
	result float64
}

// NewCalculator creates a new calculator instance
func NewCalculator() *Calculator {
	return &Calculator{
		result: 0, // Initialize result to zero
	}
}

// Add performs addition operation
func (c *Calculator) Add(a, b float64) float64 {
	// Add two numbers and return the result
	c.result = a + b
	return c.result
}

// Subtract performs subtraction operation
func (c *Calculator) Subtract(a, b float64) float64 {
	c.result = a - b // Simple subtraction
	return c.result
}

/*
Multiply function performs multiplication
It takes two float64 parameters and returns their product
*/
func (c *Calculator) Multiply(a, b float64) float64 {
	c.result = a * b
	return c.result
}

/* Divide function performs division */
func (c *Calculator) Divide(a, b float64) float64 {
	if b == 0 {
		// Cannot divide by zero
		fmt.Println("Error: Division by zero")
		return 0
	}
	c.result = a / b /* Calculate division */
	return c.result
}

// GetResult returns the current result
func (c *Calculator) GetResult() float64 {
	return c.result // Return stored result
}

/*
Main function demonstrates calculator usage
This is a multi-line comment block
*/
func main() {
	calc := NewCalculator() // Create calculator instance

	// Perform some calculations
	fmt.Println("Addition:", calc.Add(10, 5))           // Should print 15
	fmt.Println("Subtraction:", calc.Subtract(10, 3))   // Should print 7
	fmt.Println("Multiplication:", calc.Multiply(4, 6)) // Should print 24
	fmt.Println("Division:", calc.Divide(20, 4))        // Should print 5

	/* Display final result */
	fmt.Printf("Final result: %.2f\n", calc.GetResult())
}
