package types

import (
	"database/sql/driver"
	"errors"
	"math/big"
)

type Wei big.Int

type BigIntPtr = *big.Int

var bTen = big.NewInt(10)
var bEighteen = big.NewInt(18)
var bTenPowTen = new(big.Int).Exp(bTen, bTen, nil)           // 10**10
var bTenPowEighteen = new(big.Int).Exp(bTen, bEighteen, nil) // 10**18

func NewWei(x int64) *Wei {
	w := new(Wei)
	w.AsBigInt().SetInt64(x)
	return w
}

func NewUWei(x uint64) *Wei {
	w := new(Wei)
	w.AsBigInt().SetUint64(x)
	return w
}

func NewBigWei(x *big.Int) *Wei {
	w := new(Wei)
	w.AsBigInt().Set(x)
	return w
}

func SatoshiToWei(x uint64) *Wei {
	sat := new(big.Int).SetUint64(x)
	w := new(Wei)
	w.AsBigInt().Mul(sat, bTenPowTen)
	return w
}

func (w *Wei) Copy() *Wei {
	return NewBigWei(w.AsBigInt())
}

func (w *Wei) Cmp(y *Wei) int {
	return w.AsBigInt().Cmp(y.AsBigInt())
}

func (w *Wei) AsBigInt() BigIntPtr {
	return BigIntPtr(w)
}

func (w *Wei) Uint64() uint64 {
	return w.AsBigInt().Uint64()
}

func (w *Wei) ToRbtc() *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(w.AsBigInt()), new(big.Float).SetInt(bTenPowEighteen))
}

func (w *Wei) ToSatoshi() *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(w.AsBigInt()), new(big.Float).SetInt(bTenPowTen))
}

func (w *Wei) String() string {
	return w.AsBigInt().String()
}

func (w *Wei) Value() (driver.Value, error) {
	if w == nil {
		return "", errors.New("cannot retrieve value from <nil>")
	}
	return w.AsBigInt().String(), nil
}

func (w *Wei) Scan(src interface{}) error {
	switch src.(type) {
	case string:
		_, ok := w.AsBigInt().SetString(src.(string), 10)
		if !ok {
			return errors.New("cannot scan invalid value")
		}
		return nil
	case nil:
		return errors.New("cannot scan <nil> value")
	default:
		return errors.New("cannot scan invalid type of value")
	}
}

func (w *Wei) MarshalJSON() ([]byte, error) {
	return w.AsBigInt().MarshalJSON()
}

func (w *Wei) UnmarshalJSON(bytes []byte) error {
	return w.AsBigInt().UnmarshalJSON(bytes)
}

func (w *Wei) Add(x, y *Wei) *Wei {
	w.AsBigInt().Add(x.AsBigInt(), y.AsBigInt())
	return w
}

func (w *Wei) Sub(x, y *Wei) *Wei {
	w.AsBigInt().Sub(x.AsBigInt(), y.AsBigInt())
	return w
}

func (w *Wei) Mul(x, y *Wei) *Wei {
	w.AsBigInt().Mul(x.AsBigInt(), y.AsBigInt())
	return w
}
