package kll

import (
	"fmt"
	"math"
)

type KllValue interface {
	on_sum(KllValue) KllValue
	on_sub(KllValue) KllValue
	on_mul(KllValue) KllValue
	on_div(KllValue) KllValue
	re_string() string
	re_number() float64
}

type KllNumber struct {
	value float64
}

func (this KllNumber) on_sum(value KllValue) KllValue {
	return KllNumber{value: this.re_number() + value.re_number()}
}
func (this KllNumber) on_sub(value KllValue) KllValue {
	return KllNumber{value: this.re_number() - value.re_number()}
}
func (this KllNumber) on_div(value KllValue) KllValue {
	return KllNumber{value: this.re_number() / value.re_number()}
}
func (this KllNumber) on_mul(value KllValue) KllValue {
	return KllNumber{value: this.re_number() * value.re_number()}
}
func (this KllNumber) re_string() string {
	re1, re2 := math.Modf(this.value)
	if (re2) == 0 {
		return fmt.Sprint(int(re1))
	}
	return fmt.Sprint(this.value)
}
func Sum(value1 KllValue, value2 KllValue) KllValue {
	return value1.on_sum(value2)
}
func Sub(value1 KllValue, value2 KllValue) KllValue {
	return value1.on_sub(value2)
}
func Mul(value1 KllValue, value2 KllValue) KllValue {
	return value1.on_mul(value2)
}
func Div(value1 KllValue, value2 KllValue) KllValue {
	return value1.on_div(value2)
}
func Create_Number(value float64) KllValue {
	return KllNumber{value: value}
}
func (this KllNumber) re_number() float64 {
	return this.value
}
func main() {
	fmt.Println(mul(Create_Number(50), Create_Number((50))).re_string())
}
