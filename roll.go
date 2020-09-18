package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

func main() {
	input := strings.Join(os.Args[1:], " ")
	output := rollDice(input)

	result, err := evaluate(output)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s => %d\n", input, result)
}

var (
	rawRegex = `(\d*)d(\d*)`
	re       = regexp.MustCompile(rawRegex)
)

func rollDice(input string) string {
	matches := re.FindAllStringSubmatchIndex(input, -1)
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]

		rawNumDice := input[match[2]:match[3]]
		numDice := 1
		if rawNumDice != "" {
			var err error
			numDice, err = strconv.Atoi(rawNumDice)
			if err != nil {
				panic(fmt.Errorf("failed to parse number of dice: %q", rawNumDice))
			}
		}

		rawSides := input[match[4]:match[5]]
		sides, err := strconv.Atoi(rawSides)
		if err != nil {
			panic(fmt.Errorf("failed to parse dice sides: %q", rawSides))
		}

		rolls, err := rollDie(sides, numDice)
		joinedRolls := joinRolls("+", rolls)

		input = input[:match[0]] + joinedRolls + input[match[1]:]
	}
	return input
}

func rollDie(sides int, count int) ([]int, error) {
	if sides <= 0 {
		return nil, fmt.Errorf("sides must be >= 1")
	}
	if sides > 256 {
		return nil, fmt.Errorf("sides cannot be > 256")
	}
	if count <= 0 {
		return nil, fmt.Errorf("count must be >= 1")
	}

	rolls := make([]byte, count)
	_, err := rand.Read(rolls)
	if err != nil {
		return nil, fmt.Errorf("unable to get random value: %w", err)
	}
	result := []int{}
	for _, roll := range rolls {
		result = append(result, int(roll)%sides+1)
	}
	return result, nil
}

func joinRolls(sep string, rolls []int) string {
	str := &strings.Builder{}
	str.WriteString("(")
	for i := 0; i < len(rolls); i++ {
		str.WriteString(strconv.Itoa(rolls[i]))
		if i < len(rolls)-1 {
			str.WriteString(sep)
		}
	}
	str.WriteString(")")
	return str.String()
}

func evaluate(input string) (int, error) {
	expr, err := govaluate.NewEvaluableExpression(input)
	if err != nil {
		return 0, fmt.Errorf("failed to parse expression: %w", err)
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate expression: %w", err)
	}
	return int(result.(float64)), nil
}
