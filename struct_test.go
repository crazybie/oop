package inherit

import (
	"errors"
	"fmt"
	"testing"
)

type Base struct {
	Struct
	iVal int
}

type Sub struct {
	Base
	sVal string
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
	Init(sub)

	// case 1

	// up casting

	base := Cast[Base](sub)
	if base == nil || base.iVal != 1 {
		t.Errorf("up Cast failed")
	}

	sub.iVal = 2
	if base.iVal != 2 {
		t.Errorf("base.iVal!=2")
	}

	// down casting

	sub = Cast[Sub](base)
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
	base = Cast[Base](subSub)
	if base == nil || base.iVal != 2 {
		t.Errorf("up Cast failed")
	}
	// down casting 1
	sub = Cast[Sub](base)
	if sub == nil || sub.sVal != "sub" || sub.iVal != 2 {
		t.Errorf("down Cast failed")
	}
	// down casting 2
	subSub = Cast[SubSub](base)
	if subSub == nil || subSub.fVal != 22 || subSub.sVal != "sub" || subSub.iVal != 2 {
		t.Errorf("down Cast failed")
	}
	// down casting 3
	sub = Cast[Sub](base)
	if sub == nil || sub.sVal != "sub" {
		t.Errorf("down Cast failed")
	}
	// down casting 4
	subSub = Cast[SubSub](sub)
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

	base := Cast[Base](sub)
	if base == nil || base.iVal != 1 {
		t.Errorf("up Cast failed")
	}
}
