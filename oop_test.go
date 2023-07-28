package oop_test

import (
	"oop"
	"testing"
)

type Base struct {
	oop.Rtti
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

func TestOOP(t *testing.T) {
	sub := &Sub{
		Base: Base{iVal: 1},
		sVal: "sub",
	}
	oop.InitRtti(sub)

	base, ok := oop.Cast[Base](sub)
	if !ok || base.iVal != 1 {
		t.Errorf("up Cast failed")
	}

	sub.iVal = 2
	if base.iVal != 2 {
		t.Errorf("base.iVal!=2")
	}

	sub, ok = oop.Cast[Sub](base)
	if !ok || sub.sVal != "sub" || sub.iVal != 2 {
		t.Errorf("down Cast failed")
	}

	subSub := &SubSub{
		Sub:  *sub,
		fVal: 22,
	}
	oop.InitRtti(subSub)

	base, ok = oop.Cast[Base](subSub)
	if !ok || base.iVal != 2 {
		t.Errorf("up Cast failed")
	}
	sub, ok = oop.Cast[Sub](base)
	if !ok || sub.sVal != "sub" || sub.iVal != 2 {
		t.Errorf("down Cast failed")
	}

	subSub, ok = oop.Cast[SubSub](base)
	if !ok || subSub.fVal != 22 || subSub.sVal != "sub" || subSub.iVal != 2 {
		t.Errorf("down Cast failed")
	}

	sub, ok = oop.Cast[Sub](base)
	if !ok || sub.sVal != "sub" {
		t.Errorf("down Cast failed")
	}

	subSub, ok = oop.Cast[SubSub](sub)
	if !ok || subSub.fVal != 22 || subSub.sVal != "sub" || subSub.iVal != 2 {
		t.Errorf("down Cast failed")
	}
}
