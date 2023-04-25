package calculator

type Calculator struct {
	Name string
	Desc string
}

func (Calculator) Do(query string) (answer string, err error) {

	answer = "结算结果"
	return
}

func (c Calculator) GetName() string {
	return c.Name
}

func (c Calculator) GetDescription() string {
	return c.Desc
}
