package services

type Go struct {}

type Params struct {
 A int `json:"a"`
 B int `json:"b"`
}

type Result = int

func (*Go) Sub(params *Params, result *Result) error {
 a := params.A - params.B
 *result = interface{}(a).(Result)
 return nil
}

func (*Go) Sub1(params *Params, result *Result) error {
	a := params.A + params.B
	*result = interface{}(a).(Result)
	return nil
}