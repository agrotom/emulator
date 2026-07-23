package wialonips

import "fmt"

type Option interface {
	Format() string
}

type IntOption struct {
	Name  string
	Value int
}

func (opt *IntOption) Format() string {
	return fmt.Sprintf("%s:1:%d", opt.Name, opt.Value)
}

type DoubleOption struct {
	Name  string
	Value float64
}

func (opt *DoubleOption) Format() string {
	return fmt.Sprintf("%s:2:%f", opt.Name, opt.Value)
}

type StringOption struct {
	Name  string
	Value string
}

func (opt *StringOption) Format() string {
	return fmt.Sprintf("%s:3:%s", opt.Name, opt.Value)
}
