package option

type Option interface {
	Return(value interface{}) Option
	Bind(func(interface{}) Option) Option
}

type Some struct {
	Value interface{}
}

type None struct{}

func (s Some) Return(value interface{}) Option {
	return Some{value}
}

func (s Some) Bind(fn func(interface{}) Option) Option {
	return fn(s.Value)
}

func (n None) Return(value interface{}) Option {
	return None{}
}

func (n None) Bind(fn func(interface{}) Option) Option {
	return None{}
}
