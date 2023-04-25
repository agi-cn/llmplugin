package calculator

import (
	"context"
	"fmt"

	"github.com/mnogu/go-calculator"
)

const (
	desc = `A calculator, capable of performing mathematical calculations, where the input is a description of a mathematical expression and the return is the result of the calculation. For example: the input is: one plus two, the return is three.`
)

type Calculator struct {
	Name         string
	Desc         string
	InputExample string
}

func NewCalculator(name, input string) *Calculator {

	return &Calculator{
		Name:         name,
		Desc:         desc,
		InputExample: input,
	}
}

func (c Calculator) GetInputExample() string {
	return c.InputExample
}

func (Calculator) Do(ctx context.Context, query string) (answer string, err error) {

	result, err := calculator.Calculate(query)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

func (c Calculator) GetName() string {
	return c.Name
}

func (c Calculator) GetDesc() string {
	return c.Desc
}
