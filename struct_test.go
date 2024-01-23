package inherit

import (
	"errors"
	"fmt"
	"testing"
)

type Base struct {
	// ------------
	// 1 embed inherit.Struct as value in base class
	// ------------
	Struct

	iVal int
}

func (s *Base) VirtualMethod(v string) string {
	panic("not implemented")
}

type Sub struct {
	Base
	sVal string
}

func (s *Sub) VirtualMethod(v string) string {
	return v
}

type SubSub struct {
	Sub
	fVal float32
}

func TestSmoke(t *testing.T) {
	sub := &Sub{
		Base: Base{iVal: 1},
		sVal: "sub",
	}

	// ------------
	// 2 call inherit.Init on sub-class
	// ------------
	Init(sub)

	// case 1

	// up casting

	base := To[Base](sub)
	if base == nil || base.iVal != 1 {
		t.Errorf("up Cast failed")
	}

	sub.iVal = 2
	if base.iVal != 2 {
		t.Errorf("base.iVal!=2")
	}

	// down casting

	sub = To[Sub](base)
	if sub == nil || sub.sVal != "sub" || sub.iVal != 2 {
		t.Errorf("down Cast failed")
	}

	// case 2

	subSub := &SubSub{
		Sub:  *sub,
		fVal: 22,
	}
	Init(subSub)

	// up casting
	base = To[Base](subSub)
	if base == nil || base.iVal != 2 {
		t.Errorf("up Cast failed")
	}
	// down casting 1
	sub = To[Sub](base)
	if sub == nil || sub.sVal != "sub" || sub.iVal != 2 {
		t.Errorf("down Cast failed")
	}
	// down casting 2
	subSub = To[SubSub](base)
	if subSub == nil || subSub.fVal != 22 || subSub.sVal != "sub" || subSub.iVal != 2 {
		t.Errorf("down Cast failed")
	}
	// down casting 3
	sub = To[Sub](base)
	if sub == nil || sub.sVal != "sub" {
		t.Errorf("down Cast failed")
	}
	// down casting 4
	subSub = To[SubSub](sub)
	if subSub == nil || subSub.fVal != 22 || subSub.sVal != "sub" || subSub.iVal != 2 {
		t.Errorf("down Cast failed")
	}
}

func TestMissingInit(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("need panic")
		} else {
			var unInitErr *ErrUnInit
			if !errors.As(e.(error), &unInitErr) {
				t.Errorf("unexpect error %s", e)
			} else {
				fmt.Println(e.(error).Error())
			}
		}
	}()

	sub := &Sub{
		Base: Base{iVal: 1},
		sVal: "sub",
	}

	base := To[Base](sub)
	if base == nil || base.iVal != 1 {
		t.Errorf("up Cast failed")
	}
}

func TestCall(t *testing.T) {
	sub := &Sub{
		Base: Base{iVal: 1},
		sVal: "sub",
	}
	Init(sub)

	base := To[Base](sub)
	out := Invoke1_1(base.VirtualMethod, base, "1")
	if out != "1" {
		t.Fail()
	}

	out = Invoke1_1(base.VirtualMethod, base, "2")
	if out != "2" {
		t.Fail()
	}
}
