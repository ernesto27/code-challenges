#!/usr/bin/env python3
"""
A simple calculator module that demonstrates basic arithmetic operations.
This module provides functions for addition, subtraction, multiplication, and division.
"""

import math


def add(a, b):
    """
    Add two numbers and return the result.
    
    Args:
        a (float): First number
        b (float): Second number
        
    Returns:
        float: Sum of a and b
    """
    return a + b


def subtract(a, b):
    """
    Subtract second number from first number.
    
    Args:
        a (float): First number (minuend)
        b (float): Second number (subtrahend)
        
    Returns:
        float: Difference of a and b
    """
    return a - b


def multiply(a, b):
    """
    Multiply two numbers and return the result.
    
    Args:
        a (float): First number
        b (float): Second number
        
    Returns:
        float: Product of a and b
    """
    return a * b


def divide(a, b):
    """
    Divide first number by second number.
    
    Args:
        a (float): Dividend
        b (float): Divisor
        
    Returns:
        float: Quotient of a divided by b
        
    Raises:
        ValueError: If divisor is zero
    """
    if b == 0:
        raise ValueError("Cannot divide by zero")
    return a / b


def power(base, exponent):
    """
    Calculate base raised to the power of exponent.
    
    Args:
        base (float): Base number
        exponent (float): Exponent
        
    Returns:
        float: Result of base^exponent
    """
    return math.pow(base, exponent)


def main():
    """
    Main function to demonstrate calculator operations.
    """
    # Test basic arithmetic operations
    print("Calculator Demo")
    print("=" * 20)
    
    # Addition example
    x, y = 10, 5
    result = add(x, y)
    print("{} + {} = {}".format(x, y, result))
    
    # Subtraction example
    result = subtract(x, y)
    print("{} - {} = {}".format(x, y, result))
    
    # Multiplication example
    result = multiply(x, y)
    print("{} * {} = {}".format(x, y, result))
    
    # Division example
    try:
        result = divide(x, y)
        print("{} / {} = {}".format(x, y, result))
    except ValueError as e:
        print("Error: {}".format(e))
    
    # Power example
    result = power(x, y)
    print("{} ^ {} = {}".format(x, y, result))


# Execute main function if script is run directly
if __name__ == "__main__":
    main()