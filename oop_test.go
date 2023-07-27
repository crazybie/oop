package oop_test

import (
	"oop"
	"testing"
)

type Base struct {
	oop.Rtti
	B int
}

type Sub struct {
	Base
	S string
}

type SubSub struct {
	Sub
	F float32
}

func TestOOP(t *testing.T) {
	s := &Sub{
		Base: Base{B: 1},
		S:    "s",
	}
	oop.InitRtti(s)

	b, ok := oop.Cast[Base](s)
	if !ok || b.B != 1 {
		t.Errorf("up Cast failed")
	}

	s.B = 2
	if b.B != 2 {
		t.Errorf("b.B!=2")
	}

	s, ok = oop.Cast[Sub](b)
	if !ok || s.S != "s" || s.B != 2 {
		t.Errorf("down Cast failed")
	}

	ss := &SubSub{
		Sub: *s,
		F:   22,
	}
	oop.InitRtti(ss)

	b, ok = oop.Cast[Base](ss)
	if !ok || b.B != 2 {
		t.Errorf("up Cast failed")
	}
	s, ok = oop.Cast[Sub](b)
	if !ok || s.S != "s" || s.B != 2 {
		t.Errorf("down Cast failed")
	}

	ss, ok = oop.Cast[SubSub](b)
	if !ok || ss.F != 22 || ss.S != "s" || ss.B != 2 {
		t.Errorf("down Cast failed")
	}

	s, ok = oop.Cast[Sub](b)
	if !ok || s.S != "s" {
		t.Errorf("down Cast failed")
	}

	ss, ok = oop.Cast[SubSub](s)
	if !ok || ss.F != 22 || ss.S != "s" || ss.B != 2 {
		t.Errorf("down Cast failed")
	}
}
