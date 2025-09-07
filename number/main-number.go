package main

import (
	"fmt"
	"strconv"
	"strings"
)

// RPNCalculator represents a Reverse Polish Notation calculator
type RPNCalculator struct {
	stack []float64
}

// NewRPNCalculator creates a new RPN calculator instance
func NewRPNCalculator() *RPNCalculator {
	return &RPNCalculator{
		stack: make([]float64, 0),
	}
}

// Push adds a number to the stack
func (calc *RPNCalculator) Push(value float64) {
	calc.stack = append(calc.stack, value)
}

// Pop removes and returns the top element from the stack
func (calc *RPNCalculator) Pop() (float64, error) {
	if len(calc.stack) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	index := len(calc.stack) - 1
	value := calc.stack[index]
	calc.stack = calc.stack[:index]
	return value, nil
}

// Peek returns the top element without removing it
func (calc *RPNCalculator) Peek() (float64, error) {
	if len(calc.stack) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}
	return calc.stack[len(calc.stack)-1], nil
}

// IsEmpty checks if the stack is empty
func (calc *RPNCalculator) IsEmpty() bool {
	return len(calc.stack) == 0
}

// Size returns the number of elements in the stack
func (calc *RPNCalculator) Size() int {
	return len(calc.stack)
}

// Clear empties the stack
func (calc *RPNCalculator) Clear() {
	calc.stack = calc.stack[:0]
}

// Evaluate processes a single token (number or operator)
func (calc *RPNCalculator) Evaluate(token string) error {
	switch token {
	case "+":
		return calc.performBinaryOperation(func(a, b float64) float64 { return a + b })
	case "-":
		return calc.performBinaryOperation(func(a, b float64) float64 { return a - b })
	case "*":
		return calc.performBinaryOperation(func(a, b float64) float64 { return a * b })
	case "/":
		return calc.performBinaryOperation(func(a, b float64) float64 { return a / b })
	case "^", "**":
		return calc.performBinaryOperation(func(a, b float64) float64 {
			result := 1.0
			for i := 0; i < int(b); i++ {
				result *= a
			}
			return result
		})
	default:
		if value, err := strconv.ParseFloat(token, 64); err == nil {
			calc.Push(value)
			return nil
		}
		return fmt.Errorf("unknown token: %s", token)
	}
}

// performBinaryOperation applies a binary operation to the top two stack elements
func (calc *RPNCalculator) performBinaryOperation(operation func(float64, float64) float64) error {
	if len(calc.stack) < 2 {
		return fmt.Errorf("insufficient operands for operation")
	}

	// Pop second operand first (top of stack)
	b, _ := calc.Pop()
	// Pop first operand (second from top)
	a, _ := calc.Pop()

	result := operation(a, b)
	calc.Push(result)
	return nil
}

// EvaluateExpression processes an entire RPN expression and returns the result
func (calc *RPNCalculator) EvaluateExpression(expression string) (float64, error) {
	calc.Clear()
	tokens := strings.Fields(expression)

	for _, token := range tokens {
		if err := calc.Evaluate(token); err != nil {
			return 0, err
		}
	}

	if calc.Size() != 1 {
		return 0, fmt.Errorf("invalid expression: expected 1 result, got %d", calc.Size())
	}

	return calc.Peek()
}

// PrintStack displays the current stack contents
func (calc *RPNCalculator) PrintStack() {
	fmt.Print("Stack: [")
	for i, value := range calc.stack {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%.2f", value)
	}
	fmt.Println("]")
}

func runNumbersDemo() {
	fmt.Println("=== Reverse Polish Notation Calculator Demo ===\n")

	calc := NewRPNCalculator()

	// Example 1: Simple addition - 3 + 2 + 4
	fmt.Println("Example 1:")
	fmt.Println("Input equation: 3 + 2 + 4")
	fmt.Println("Infix expression: 3 + 2 + 4")
	fmt.Println("RPN expression: 3 2 + 4 +")
	fmt.Println("Step-by-step evaluation:")

	calc.Clear()
	tokens1 := []string{"3", "2", "+", "4", "+"}

	for i, token := range tokens1 {
		fmt.Printf("  Step %d: Process '%s'", i+1, token)
		if err := calc.Evaluate(token); err != nil {
			fmt.Printf(" -> Error: %v\n", err)
			break
		}
		fmt.Print(" -> ")
		calc.PrintStack()
	}

	if !calc.IsEmpty() {
		finalResult, _ := calc.Peek()
		fmt.Printf("Final result: %.0f\n\n", finalResult)
	}

	// Example 2: (3 + 4) × (5 + 6)
	fmt.Println("Example 2:")
	fmt.Println("Input equation: (3 + 4) × (5 + 6)")
	fmt.Println("Infix expression: (3 + 4) * (5 + 6)")
	fmt.Println("RPN expression: 3 4 + 5 6 + *")
	fmt.Println("Step-by-step evaluation:")

	calc.Clear()
	tokens2 := []string{"3", "4", "+", "5", "6", "+", "*"}

	for i, token := range tokens2 {
		fmt.Printf("  Step %d: Process '%s'", i+1, token)
		if err := calc.Evaluate(token); err != nil {
			fmt.Printf(" -> Error: %v\n", err)
			break
		}
		fmt.Print(" -> ")
		calc.PrintStack()
	}

	if !calc.IsEmpty() {
		finalResult, _ := calc.Peek()
		fmt.Printf("Final result: %.0f\n\n", finalResult)
	}

	// Example 3: More complex - ((15 / 3) + 2) * (8 - 3)
	fmt.Println("Example 3:")
	fmt.Println("Input equation: ((15 / 3) + 2) * (8 - 3)")
	fmt.Println("Infix expression: ((15 / 3) + 2) * (8 - 3)")
	fmt.Println("RPN expression: 15 3 / 2 + 8 3 - *")
	fmt.Println("Step-by-step evaluation:")

	calc.Clear()
	tokens3 := []string{"15", "3", "/", "2", "+", "8", "3", "-", "*"}

	for i, token := range tokens3 {
		fmt.Printf("  Step %d: Process '%s'", i+1, token)
		if err := calc.Evaluate(token); err != nil {
			fmt.Printf(" -> Error: %v\n", err)
			break
		}
		fmt.Print(" -> ")
		calc.PrintStack()
	}

	if !calc.IsEmpty() {
		finalResult, _ := calc.Peek()
		fmt.Printf("Final result: %.0f\n", finalResult)
	}
}

func main() {
	runNumbersDemo()
}
