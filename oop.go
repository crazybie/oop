package oop

import (
	"reflect"
	"sync"
	"unsafe"
)

var tpLookup = struct {
	sync.RWMutex
	types map[reflect.Type]map[reflect.Type]bool
}{
	types: make(map[reflect.Type]map[reflect.Type]bool),
}

type IRtti interface {
	getType() reflect.Type
	setType(reflect.Type)
}

type Rtti struct {
	tp reflect.Type
}

func InitRtti(imp IRtti) {
	imp.setType(reflect.TypeOf(imp).Elem())
}

func (b *Rtti) getType() reflect.Type {
	return b.tp
}

func (b *Rtti) setType(t reflect.Type) {
	b.tp = t
}

func Cast[To any](s IRtti) (d *To, ok bool) {
	dstTp := reflect.TypeOf(d).Elem()
	rTp := s.getType()
	if dstTp != rTp {

		// check cache

		tpLookup.RLock()
		if v, ok1 := tpLookup.types[rTp]; ok1 {
			if v2, ok2 := v[dstTp]; ok2 && v2 {
				tpLookup.RUnlock()
				if v2 {
					goto success
				} else {
					return d, false
				}
			}
		}
		tpLookup.RUnlock()

		// save to cache

		f := rTp.Field(0).Type
		for f != dstTp && f.Kind() == reflect.Struct {
			f = f.Field(0).Type
		}
		ok = f == dstTp

		tpLookup.Lock()
		if tpLookup.types[rTp] == nil {
			tpLookup.types[rTp] = make(map[reflect.Type]bool)
		}
		tpLookup.types[rTp][dstTp] = ok
		tpLookup.Unlock()
		if !ok {
			return d, false
		}
	}
success:
	d = (*To)((*[2]unsafe.Pointer)(unsafe.Pointer(&s))[1])
	return d, true
}
