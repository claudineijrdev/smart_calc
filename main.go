package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var operationSymbols = "+-/*%^"

type Stack []string

func (st *Stack) IsEmpty() bool {
	return len(*st) == 0
}

func (st *Stack) Push(str string) {
	*st = append(*st, str)
}

func (st *Stack) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(*st) - 1
		*st = (*st)[:index]
		return true
	}
}

func (st *Stack) Top() string {
	if st.IsEmpty() {
		return ""
	} else {
		index := len(*st) - 1
		element := (*st)[index]
		return element
	}
}

func prec(s string) int {
	if s == "^" {
		return 3
	} else if (s == "/") || (s == "*") {
		return 2
	} else if (s == "+") || (s == "-") {
		return 1
	} else {
		return -1
	}
}

func infixToPostfix(infix []string) []string {
	var sta Stack
	var postfix []string
	for _, char := range infix {
		opchar := string(char)
		_, err := strconv.Atoi(opchar)
		if err == nil {
			postfix = append(postfix, opchar)
			continue
		}
		if (char >= "a" && char <= "z") || (char >= "A" && char <= "Z") || (char >= "0" && char <= "9") {
			postfix = append(postfix, opchar)
		} else if char == "(" {
			sta.Push(opchar)
		} else if char == ")" {
			for sta.Top() != "(" {
				postfix = append(postfix, sta.Top())
				sta.Pop()
			}
			sta.Pop()
		} else {
			for !sta.IsEmpty() && prec(opchar) <= prec(sta.Top()) {
				postfix = append(postfix, sta.Top())
				sta.Pop()
			}
			sta.Push(opchar)
		}
	}
	for !sta.IsEmpty() {
		postfix = append(postfix, sta.Top())
		sta.Pop()
	}
	return postfix
}

func normalize(inputs []string) ([]string, error) {
	var result []string

	for _, char := range inputs {
		if len(char) == 1 {
			result = append(result, char)
			continue
		}
		op, err := normalizeOperations(char)
		if err != nil {
			result = append(result, char)
			return result, err
		}
		result = append(result, op)
	}
	return result, nil
}

func normalizeOperations(operation string) (string, error) {
	minus := -1
	plus := 1
	resp := 1

	if !strings.ContainsAny(operation, "+-/*%") && operation != "" {
		return operation, nil
	}

	if len(operation) == 1 {
		return operation, nil
	}

	for _, char := range operation {
		switch char {
		case '+':
			resp *= plus
		case '-':
			resp *= minus
		default:
			return "", fmt.Errorf("Invalid operation")
		}
	}

	if resp == plus {
		return "+", nil
	}
	return "-", nil
}

var operations = map[string]func(int, int) int{
	"+": func(a, b int) int {
		return a + b
	},
	"-": func(a, b int) int {
		return a - b
	},
	"*": func(a, b int) int {
		return a * b
	},
	"/": func(a, b int) int {
		return int(a / b)
	},
	"%": func(a, b int) int {
		return a % b
	},
	"^": func(a, b int) int {
		return int(math.Pow(float64(a), float64(b)))
	},
}

func calc(inputs []string) (int, error) {
	sum := 0

	if strings.Count(strings.Join(inputs, ""), "(") != strings.Count(strings.Join(inputs, ""), ")") {
		return 0, fmt.Errorf("Invalid expression")
	}

	expression, err := normalize(inputs)
	if err != nil {
		return 0, err
	}

	postfix := infixToPostfix(expression)

	var stack Stack

	for _, char := range postfix {
		if strings.ContainsAny(char, operationSymbols) {
			if stack.IsEmpty() {
				return 0, fmt.Errorf("Invalid expression")
			}
			operation := operations[char]
			b, err := strconv.Atoi(stack.Top())
			if err != nil {
				return 0, err
			}
			stack.Pop()
			a, err := strconv.Atoi(stack.Top())
			if err != nil {
				return 0, err
			}
			stack.Pop()
			sum = operation(a, b)
			stack.Push(strconv.Itoa(sum))
			continue
		}
		stack.Push(char)
	}

	if stack.IsEmpty() {
		return 0, fmt.Errorf("invalid expression")
	}

	sum, err = strconv.Atoi(stack.Top())
	if err != nil {
		return 0, err
	}

	return sum, nil
}

type Mode int

const (
	Normal Mode = iota
	Assignment
	Result
	Calculation
	CalculationAssignment
	Command
	Error
)

func getMode(expression []string) Mode {
	if strings.HasPrefix(expression[0], "/") {
		return Command
	}
	strExpression := strings.Join(expression, "")

	if len(expression) == 1 {
		if IsVariable(strExpression) || unicode.IsDigit(rune(strExpression[0])) {
			return Result
		}
	}

	if len(expression) == 2 {
		return Error
	}

	if expression[1] == "=" {
		if IsVariable(string(expression[0][0])) {
			if strings.ContainsAny(strExpression, operationSymbols) {
				return CalculationAssignment
			}
			if len(expression) == 3 {
				return Assignment
			}
		}
		return Error
	}

	return Calculation
}

func Assign(input []string, variables *map[string]int) error {
	if len(input) != 3 {
		return fmt.Errorf("Invalid assignment")
	}

	if !IsValidVariableName(input[0]) {
		return fmt.Errorf("Invalid identifier")
	}

	_, err := strconv.Atoi(input[2])
	if err != nil {
		if !IsValidVariableName(input[2]) {
			return fmt.Errorf("Invalid assignment")
		}
		if _, ok := (*variables)[input[2]]; ok {
			(*variables)[input[0]] = (*variables)[input[2]]
			return nil
		}
		return fmt.Errorf("Invalid assignment")
	}
	(*variables)[input[0]], _ = strconv.Atoi(input[2])
	return nil
}

func PreCalc(input []string, variables *map[string]int) ([]string, error) {
	var result []string
	for _, char := range input {
		if IsVariable(char) {
			v, ok := (*variables)[char]
			if ok {
				result = append(result, strconv.Itoa(v))
				continue
			}
			return nil, fmt.Errorf("unknown variable")
		}
		result = append(result, char)
	}
	return result, nil
}

func IsVariable(input string) bool {
	return unicode.IsLetter(rune(input[0]))
}

func IsValidVariableName(name string) bool {
	for _, char := range name {
		if !unicode.In(char, unicode.Letter, unicode.Latin) {
			return false
		}
	}
	return true
}

func CreateExpression(input string) []string {
	if strings.HasPrefix(input, "/") {
		return []string{input}
	}

	if strings.HasPrefix(input, "-") {
		input = "0" + input
	}

	input = strings.ReplaceAll(input, " ", "")
	var result []string
	var number string
	var operations string
	var varName string
	for _, char := range input {
		if IsVariable(string(char)) {
			if operations != "" {
				result = append(result, operations)
				operations = ""
			}
			if number != "" {
				result = append(result, number)
				number = ""
			}

			varName += string(char)
			continue
		}

		if unicode.IsDigit(char) {
			if operations != "" {
				result = append(result, operations)
				operations = ""
			}
			if varName != "" {
				result = append(result, varName)
				varName = ""
			}

			number += string(char)
			continue
		}
		if strings.ContainsAny(string(char), operationSymbols) {
			if number != "" {
				result = append(result, number)
				number = ""
			}
			if varName != "" {
				result = append(result, varName)
				varName = ""
			}
			operations += string(char)
			continue
		}

		if operations != "" {
			result = append(result, operations)
			operations = ""
		}
		if number != "" {
			result = append(result, number)
			number = ""
		}
		if varName != "" {
			result = append(result, varName)
			varName = ""
		}

		result = append(result, string(char))
	}
	if operations != "" {
		result = append(result, operations)
	}
	if number != "" {
		result = append(result, number)
	}
	if varName != "" {
		result = append(result, varName)
	}
	return result
}

func run(text string, variables *map[string]int) (int, bool, Mode, error) {
	var exitCommand = "/exit"
	var helpCommand = "/help"
	var helpMessage = "The program calculates expressions"

	expression := CreateExpression(text)
	mode := getMode(expression)

	switch mode {
	case Error:
		{
			if len(expression) > 0 {
				return 0, true, mode, fmt.Errorf("Invalid expression")
			}
			return 0, true, mode, nil
		}
	case Command:
		{
			if expression[0] == exitCommand {
				return 0, false, mode, nil
			}
			if expression[0] == helpCommand {
				return 0, true, mode, errors.New(helpMessage)
			}
			return 0, true, mode, fmt.Errorf("Unknown command")
		}
	case Result:
		{
			if IsVariable(expression[0]) {
				if _, ok := (*variables)[expression[0]]; ok {
					return (*variables)[expression[0]], true, mode, nil
				}
				return 0, true, mode, fmt.Errorf("Unknown variable")
			}
			value, err := strconv.Atoi(expression[0])
			if err != nil {
				return 0, true, mode, fmt.Errorf("Invalid expression")
			}
			return value, true, mode, nil
		}
	case Calculation:
		{
			in, err := PreCalc(expression, variables)
			if err != nil {
				return 0, true, mode, fmt.Errorf("Invalid expression")
			}
			result, err := calc(in)
			if err != nil {
				return 0, true, mode, fmt.Errorf("Invalid expression")
			}
			return result, true, mode, nil
		}
	case Assignment:
		{
			err := Assign(expression, variables)
			return 0, true, mode, err
		}
	case CalculationAssignment:
		{
			in, err := PreCalc(expression[2:], variables)
			if err != nil {
				return 0, true, mode, fmt.Errorf("Invalid expression")
			}
			result, err := calc(in)
			if err != nil {
				return 0, true, mode, fmt.Errorf("Invalid expression")
			}
			assignIn := []string{expression[0], "=", strconv.Itoa(result)}
			err = Assign(assignIn, variables)
			return 0, true, mode, err
		}
	}
	return 0, true, Error, nil
}

func main() {
	variables := make(map[string]int)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		expression := scanner.Text()
		if len(expression) == 0 || expression == "" {
			continue
		}
		response, shouldContinue, mode, err := run(expression, &variables)
		if !shouldContinue {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		if mode != Assignment && mode != CalculationAssignment {
			fmt.Println(response)
		}
	}
	fmt.Println("Bye!")
}
