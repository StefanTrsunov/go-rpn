package main

import (
	"fmt"
	"strings"
)

// Documents to search through
var documents = []string{"C++ Guide", "Java guide tutorial", "Python tutorial", "C tutorial"}

// Operator precedence for boolean operations
var precedence = map[string]int{
	"NOT": 3,
	"AND": 2,
	"OR":  1,
	"(":   0,
}

// BooleanRPNProcessor represents a boolean query processor using RPN
type BooleanRPNProcessor struct {
	stack []bool
}

// NewBooleanRPNProcessor creates a new boolean RPN processor
func NewBooleanRPNProcessor() *BooleanRPNProcessor {
	return &BooleanRPNProcessor{
		stack: make([]bool, 0),
	}
}

// Push adds a boolean value to the stack
func (proc *BooleanRPNProcessor) Push(value bool) {
	proc.stack = append(proc.stack, value)
}

// Pop removes and returns the top boolean value from the stack
func (proc *BooleanRPNProcessor) Pop() (bool, error) {
	if len(proc.stack) == 0 {
		return false, fmt.Errorf("stack is empty")
	}

	index := len(proc.stack) - 1
	value := proc.stack[index]
	proc.stack = proc.stack[:index]
	return value, nil
}

// Clear empties the stack
func (proc *BooleanRPNProcessor) Clear() {
	proc.stack = proc.stack[:0]
}

// Size returns the number of elements in the stack
func (proc *BooleanRPNProcessor) Size() int {
	return len(proc.stack)
}

// ConvertOperands converts search terms in query to T/F based on document content
func convertOperands(query, document string) string {
	word := ""
	convertedQuery := query
	queryLower := strings.ToLower(query)
	documentLower := strings.ToLower(document)

	for _, char := range queryLower {
		if word == "AND" || word == "OR" || word == "NOT" {
			word = ""
			continue
		}

		if char == ' ' || char == '(' || char == ')' {
			if word != "" {
				replacement := "F"
				if strings.Contains(documentLower, word) {
					replacement = "T"
				}
				convertedQuery = strings.ReplaceAll(convertedQuery, strings.TrimSpace(word), replacement)
				word = ""
			}
			continue
		}
		word += string(char)
	}

	// Handle the last word if exists
	if word != "" {
		replacement := "F"
		if strings.Contains(documentLower, word) {
			replacement = "T"
		}
		convertedQuery = strings.ReplaceAll(convertedQuery, word, replacement)
	}

	return convertedQuery
}

// Tokenize breaks the query into tokens
func tokenize(query string) []string {
	word := ""
	tokens := []string{}

	for _, char := range query {
		if char == ' ' {
			word = ""
			continue
		}

		if word == "" && (char == '(' || char == ')' || char == 'T' || char == 'F') {
			tokens = append(tokens, string(char))
			continue
		}

		word += string(char)

		if word == "AND" || word == "OR" || word == "NOT" {
			tokens = append(tokens, word)
			word = ""
		}
	}

	return tokens
}

// BuildRPN converts infix boolean expression to RPN using Shunting Yard algorithm
func buildRPN(tokens []string) []string {
	output := []string{}
	operations := []string{}

	for _, token := range tokens {
		if token == "(" {
			operations = append(operations, token)
			continue
		}

		if token == ")" {
			for len(operations) > 0 && operations[len(operations)-1] != "(" {
				output = append(output, operations[len(operations)-1])
				operations = operations[:len(operations)-1]
			}
			// Remove the opening parenthesis
			if len(operations) > 0 {
				operations = operations[:len(operations)-1]
			}
			continue
		}

		if token == "AND" || token == "OR" || token == "NOT" {
			for len(operations) > 0 && precedence[operations[len(operations)-1]] >= precedence[token] {
				output = append(output, operations[len(operations)-1])
				operations = operations[:len(operations)-1]
			}
			operations = append(operations, token)
		} else {
			output = append(output, token)
		}
	}

	// Pop remaining operations
	for len(operations) > 0 {
		output = append(output, operations[len(operations)-1])
		operations = operations[:len(operations)-1]
	}

	return output
}

// EvaluateRPN evaluates a boolean RPN expression
func (proc *BooleanRPNProcessor) EvaluateRPN(rpn []string) (bool, error) {
	proc.Clear()

	for _, token := range rpn {
		switch token {
		case "T":
			proc.Push(true)
		case "F":
			proc.Push(false)
		case "AND":
			if proc.Size() < 2 {
				return false, fmt.Errorf("insufficient operands for AND operation")
			}
			second, _ := proc.Pop()
			first, _ := proc.Pop()
			proc.Push(first && second)
		case "OR":
			if proc.Size() < 2 {
				return false, fmt.Errorf("insufficient operands for OR operation")
			}
			second, _ := proc.Pop()
			first, _ := proc.Pop()
			proc.Push(first || second)
		case "NOT":
			if proc.Size() < 1 {
				return false, fmt.Errorf("insufficient operands for NOT operation")
			}
			operand, _ := proc.Pop()
			proc.Push(!operand)
		default:
			return false, fmt.Errorf("unknown token: %s", token)
		}
	}

	if proc.Size() != 1 {
		return false, fmt.Errorf("invalid expression: expected 1 result, got %d", proc.Size())
	}

	result, _ := proc.Pop()
	return result, nil
}

// Match checks if a document matches the given boolean query
func match(query, document string) bool {
	fmt.Printf("Query: %s ===> Document: %s\n", query, document)

	// Convert search terms to T/F based on document content
	convertedQuery := convertOperands(query, document)
	fmt.Printf("Converted Query: %s\n", convertedQuery)

	// Tokenize the converted query
	tokens := tokenize(convertedQuery)
	fmt.Printf("Tokenized Query: %v\n", tokens)

	// Build RPN from tokens
	rpnQuery := buildRPN(tokens)
	fmt.Printf("RPN Query: %v\n", rpnQuery)

	// Evaluate RPN expression
	processor := NewBooleanRPNProcessor()
	result, err := processor.EvaluateRPN(rpnQuery)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}

	fmt.Printf("Result: %t\n\n", result)
	return result
}

func runDocumentsDemo() {
	fmt.Println("=== Boolean Query Processing with RPN ===\n")

	// Available documents for searching
	fmt.Println("Available documents:")
	for i, doc := range documents {
		fmt.Printf("  %d. %s\n", i+1, doc)
	}
	fmt.Println()

	// Example 1: Simple term search
	fmt.Println("Example 1:")
	fmt.Println("Input query: python")
	fmt.Println("Looking for documents containing 'python'")
	fmt.Println()

	query1 := "python"
	fmt.Printf("Processing query: %s\n", query1)

	matches1 := []string{}
	for _, doc := range documents {
		fmt.Printf("  Document: \"%s\"\n", doc)

		// Convert terms to T/F
		converted := convertOperands(query1, doc)
		fmt.Printf("  Converted: %s\n", converted)

		// Since it's just a single term, no RPN conversion needed
		// Just check if the term exists
		docLower := strings.ToLower(doc)
		queryLower := strings.ToLower(query1)
		result := strings.Contains(docLower, queryLower)
		fmt.Printf("  Result: %t\n", result)

		if result {
			matches1 = append(matches1, doc)
		}
		fmt.Println()
	}

	fmt.Printf("Matching documents: ")
	if len(matches1) == 0 {
		fmt.Println("None")
	} else {
		fmt.Printf("%v\n", matches1)
	}
	fmt.Println(strings.Repeat("-", 60))

	// Example 2: AND operation
	fmt.Println("Example 2:")
	fmt.Println("Input query: python AND tutorial")
	fmt.Println("Looking for documents containing both 'python' AND 'tutorial'")
	fmt.Println()

	query2 := "python AND tutorial"
	fmt.Printf("Processing query: %s\n", query2)

	// Demonstrate with one document step-by-step
	testDoc := "Python tutorial"
	fmt.Printf("Step-by-step for document: \"%s\"\n", testDoc)

	// Step 1: Convert operands
	converted2 := convertOperands(query2, testDoc)
	fmt.Printf("  Step 1 - Convert terms: %s\n", converted2)

	// Step 2: Tokenize
	tokens2 := tokenize(converted2)
	fmt.Printf("  Step 2 - Tokenize: %v\n", tokens2)

	// Step 3: Build RPN
	rpn2 := buildRPN(tokens2)
	fmt.Printf("  Step 3 - Build RPN: %v\n", rpn2)

	// Step 4: Evaluate RPN
	fmt.Println("  Step 4 - Evaluate RPN:")
	processor := NewBooleanRPNProcessor()

	for i, token := range rpn2 {
		fmt.Printf("    Step %d: Process '%s'", i+1, token)

		switch token {
		case "T":
			processor.Push(true)
			fmt.Printf(" -> Push true")
		case "F":
			processor.Push(false)
			fmt.Printf(" -> Push false")
		case "AND":
			second, _ := processor.Pop()
			first, _ := processor.Pop()
			result := first && second
			processor.Push(result)
			fmt.Printf(" -> Pop %t and %t, push %t", first, second, result)
		case "OR":
			second, _ := processor.Pop()
			first, _ := processor.Pop()
			result := first || second
			processor.Push(result)
			fmt.Printf(" -> Pop %t and %t, push %t", first, second, result)
		case "NOT":
			operand, _ := processor.Pop()
			result := !operand
			processor.Push(result)
			fmt.Printf(" -> Pop %t, push %t", operand, result)
		}

		fmt.Printf(" -> Stack: %v\n", processor.stack)
	}

	finalResult2, _ := processor.Pop()
	fmt.Printf("  Final result: %t\n\n", finalResult2)

	// Check all documents for this query
	matches2 := []string{}
	for _, doc := range documents {
		converted := convertOperands(query2, doc)
		tokens := tokenize(converted)
		rpn := buildRPN(tokens)

		proc := NewBooleanRPNProcessor()
		result, err := proc.EvaluateRPN(rpn)
		if err == nil && result {
			matches2 = append(matches2, doc)
		}
	}

	fmt.Printf("All matching documents: ")
	if len(matches2) == 0 {
		fmt.Println("None")
	} else {
		fmt.Printf("%v\n", matches2)
	}
	fmt.Println(strings.Repeat("-", 60))

	// Example 3: Complex query with OR and parentheses
	fmt.Println("Example 3:")
	fmt.Println("Input query: (python OR java) AND guide")
	fmt.Println("Looking for documents with (python OR java) AND guide")
	fmt.Println()

	query3 := "(python OR java) AND guide"
	fmt.Printf("Processing query: %s\n", query3)

	// Demonstrate with one document step-by-step
	testDoc3 := "Java guide tutorial"
	fmt.Printf("Step-by-step for document: \"%s\"\n", testDoc3)

	// Step 1: Convert operands
	converted3 := convertOperands(query3, testDoc3)
	fmt.Printf("  Step 1 - Convert terms: %s\n", converted3)

	// Step 2: Tokenize
	tokens3 := tokenize(converted3)
	fmt.Printf("  Step 2 - Tokenize: %v\n", tokens3)

	// Step 3: Build RPN
	rpn3 := buildRPN(tokens3)
	fmt.Printf("  Step 3 - Build RPN: %v\n", rpn3)

	// Step 4: Evaluate RPN
	fmt.Println("  Step 4 - Evaluate RPN:")
	processor3 := NewBooleanRPNProcessor()

	for i, token := range rpn3 {
		fmt.Printf("    Step %d: Process '%s'", i+1, token)

		switch token {
		case "T":
			processor3.Push(true)
			fmt.Printf(" -> Push true")
		case "F":
			processor3.Push(false)
			fmt.Printf(" -> Push false")
		case "AND":
			second, _ := processor3.Pop()
			first, _ := processor3.Pop()
			result := first && second
			processor3.Push(result)
			fmt.Printf(" -> Pop %t and %t, push %t", first, second, result)
		case "OR":
			second, _ := processor3.Pop()
			first, _ := processor3.Pop()
			result := first || second
			processor3.Push(result)
			fmt.Printf(" -> Pop %t and %t, push %t", first, second, result)
		case "NOT":
			operand, _ := processor3.Pop()
			result := !operand
			processor3.Push(result)
			fmt.Printf(" -> Pop %t, push %t", operand, result)
		}

		fmt.Printf(" -> Stack: %v\n", processor3.stack)
	}

	finalResult3, _ := processor3.Pop()
	fmt.Printf("  Final result: %t\n\n", finalResult3)

	// Check all documents for this query
	matches3 := []string{}
	for _, doc := range documents {
		converted := convertOperands(query3, doc)
		tokens := tokenize(converted)
		rpn := buildRPN(tokens)

		proc := NewBooleanRPNProcessor()
		result, err := proc.EvaluateRPN(rpn)
		if err == nil && result {
			matches3 = append(matches3, doc)
		}
	}

	fmt.Printf("All matching documents: ")
	if len(matches3) == 0 {
		fmt.Println("None")
	} else {
		fmt.Printf("%v\n", matches3)
	}
}

func main() {
	runDocumentsDemo()
}
