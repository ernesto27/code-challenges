/**
 * Calculator module for basic arithmetic operations
 * @author Developer
 * @version 1.0
 */

// Import required modules
const readline = require('readline');

/* 
   Multi-line comment explaining the purpose
   This calculator supports basic operations
*/
class Calculator {
    /**
     * Constructor for Calculator class
     */
    constructor() {
        this.history = []; // Store calculation history
    }

    // Single line comment before method
    add(a, b) {
        const result = a + b; // Inline comment
        /* Block comment in method */
        this.history.push(`${a} + ${b} = ${result}`);
        return result;
    }

    /*
     * Subtract two numbers
     * @param {number} a - First number
     * @param {number} b - Second number  
     * @returns {number} The difference
     */
    subtract(a, b) {
        const result = a - b;
        this.history.push(`${a} - ${b} = ${result}`);
        return result;
    }

    multiply(a, b) {
        // Calculate multiplication
        const result = a * b;
        return result; // Return the product
    }

    divide(a, b) {
        if (b === 0) {
            throw new Error("Division by zero"); // Error handling
        }
        /* Calculate division result */
        const result = a / b;
        return result;
    }

    // Method to display calculation history
    showHistory() {
        console.log("// Calculation History //");
        this.history.forEach(entry => {
            console.log(entry); // Print each entry
        });
    }
}

/* 
   Main function to demonstrate calculator usage
*/
function main() {
    const calc = new Calculator(); // Create calculator instance
    
    // Perform some calculations
    console.log(calc.add(5, 3)); // Should output 8
    console.log(calc.subtract(10, 4)); // Should output 6
    console.log(calc.multiply(7, 2)); // Should output 14
    
    /* Display results */
    calc.showHistory();
}

// Export the Calculator class
module.exports = Calculator;

// Run main function if this file is executed directly
if (require.main === module) {
    main(); // Execute main function
}