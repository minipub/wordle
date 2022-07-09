package internal

import (
	"fmt"
	"math/big"
)

type Int big.Int

func NewXBigInt(x int64) *Int {
	y := big.NewInt(x)
	return (*Int)(y)
}

func (i *Int) String() string {
	return fmt.Sprint(bigInt(i))
}

// bigInt *xbig.Int to *big.Int
func bigInt(i *Int) *big.Int {
	return (*big.Int)(i)
}

func (i *Int) Copy() *Int {
	x := big.NewInt(0)
	src := bigInt(i).Bits()
	dst := make([]big.Word, len(src))
	copy(dst, src)
	x.SetBits(dst)
	return (*Int)(x)
}

// Operate with xbig.Int
// First operand will be copied before operated
// The result will be a new pointer

func Add(x *Int, y *Int) *Int {
	z := x.Copy()
	z.Add(y)
	return z
}

func Sub(x *Int, y *Int) *Int {
	z := x.Copy()
	z.Sub(y)
	return z
}

func Mul(x *Int, y *Int) *Int {
	z := x.Copy()
	z.Mul(y)
	return z
}

func Div(x *Int, y *Int) *Int {
	z := x.Copy()
	z.Div(y)
	return z
}

func Pow(x *Int, y *Int) *Int {
	z := NewXBigInt(0)
	bigInt(z).Exp(bigInt(x), bigInt(y), nil)
	return z
}

// Operate with xbig.Int
// The result apply on recv pointer

func (i *Int) Add(x *Int) *Int {
	bigInt(i).Add(bigInt(i), bigInt(x))
	return i
}

func (i *Int) Sub(x *Int) *Int {
	bigInt(i).Sub(bigInt(i), bigInt(x))
	return i
}

func (i *Int) Mul(x *Int) *Int {
	bigInt(i).Mul(bigInt(i), bigInt(x))
	return i
}

func (i *Int) Div(x *Int) *Int {
	bigInt(i).Div(bigInt(i), bigInt(x))
	return i
}

// Compare with xbig.Int

func (i *Int) Gt(x *Int) bool {
	switch bigInt(i).Cmp(bigInt(x)) {
	case 1:
		return true
	default:
		return false
	}
}

func (i *Int) Ge(x *Int) bool {
	switch bigInt(i).Cmp(bigInt(x)) {
	case 1, 0:
		return true
	default:
		return false
	}
}

func (i *Int) Eq(x *Int) bool {
	switch bigInt(i).Cmp(bigInt(x)) {
	case 0:
		return true
	default:
		return false
	}
}

func (i *Int) Lt(x *Int) bool {
	switch bigInt(i).Cmp(bigInt(x)) {
	case -1:
		return true
	default:
		return false
	}
}

func (i *Int) Le(x *Int) bool {
	switch bigInt(i).Cmp(bigInt(x)) {
	case -1, 0:
		return true
	default:
		return false
	}
}
