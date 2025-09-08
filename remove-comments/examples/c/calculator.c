/**
 * Simple Calculator in C
 * Author: Developer
 * Version: 1.0
 * 
 * This program demonstrates basic arithmetic operations
 */

#include <stdio.h>
#include <stdlib.h>

/* Define constants for operations */
#define MAX_INPUT 100
#define PI 3.14159

// Function prototypes
int add(int a, int b);       // Addition function
int subtract(int a, int b);  // Subtraction function
int multiply(int a, int b);  /* Multiplication function */
float divide(int a, int b);  /* Division function with float return */

/*
 * Main function - entry point of the program
 * Demonstrates basic calculator operations
 */
int main() {
    int num1, num2;
    char operation;
    
    // Print welcome message
    printf("Simple Calculator\n");
    printf("Enter first number: ");
    scanf("%d", &num1); // Read first number
    
    printf("Enter operation (+, -, *, /): ");
    scanf(" %c", &operation); // Note the space before %c
    
    printf("Enter second number: ");
    scanf("%d", &num2); // Read second number
    
    /* Process the operation using switch statement */
    switch(operation) {
        case '+':
            printf("Result: %d\n", add(num1, num2)); // Call add function
            break;
        case '-':
            printf("Result: %d\n", subtract(num1, num2)); // Call subtract
            break;
        case '*': /* Multiplication case */
            printf("Result: %d\n", multiply(num1, num2));
            break;
        case '/': /* Division case */
            if (num2 != 0) {
                printf("Result: %.2f\n", divide(num1, num2));
            } else {
                printf("Error: Division by zero!\n"); // Handle division by zero
            }
            break;
        default:
            printf("Error: Invalid operation!\n"); // Handle invalid input
            break;
    }
    
    return 0; // Success
}

/**
 * Add two integers
 * @param a First integer
 * @param b Second integer  
 * @return Sum of a and b
 */
int add(int a, int b) {
    return a + b; // Simple addition
}

/*
 * Subtract two integers
 */
int subtract(int a, int b) {
    int result = a - b; // Calculate difference
    return result; // Return the result
}

// Multiply two integers
int multiply(int a, int b) {
    /* Use a loop for demonstration */
    int result = 0;
    int i;
    
    // Add 'a' to result 'b' times
    for(i = 0; i < b; i++) {
        result += a; // Accumulate the sum
    }
    
    return result; // Return product
}

/**
 * Divide two integers and return float result
 */
float divide(int a, int b) {
    // Cast to float for accurate division
    return (float)a / (float)b; /* Return floating point result */
}