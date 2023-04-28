package calculator

import (
	"context"
	"fmt"

	"github.com/mnogu/go-calculator"
)

const (
	pluginName         = `Calculator`
	pluginDesc         = `A calculator, capable of performing mathematical calculations, where the input is a description of a mathematical expression and the return is the result of the calculation. For example: the input is: one plus two, the return is three.`
	pluginInputExample = `1+2`
)

type Calculator struct{}

func NewCalculator() *Calculator {

	return &Calculator{}
}

func (c Calculator) GetInputExample() string {
	return pluginInputExample
}

func (Calculator) Do(ctx context.Context, query string) (answer string, err error) {

	result, err := calculator.Calculate(query)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

func (c Calculator) GetName() string {
	return pluginName
}

func (c Calculator) GetDesc() string {
	return pluginDesc
}
