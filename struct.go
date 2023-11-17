//
// Single inheritance simulation for Golang.
//
// original created by lisong@legoutech.com
//

package inherit

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var cache = struct {
	sync.RWMutex
	casting map[reflect.Type]map[reflect.Type]bool
}{
	casting: make(map[reflect.Type]map[reflect.Type]bool),
}

type ErrUnInit struct {
	tp reflect.Type
}

func (er ErrUnInit) Error() string {
	return fmt.Sprintf("missing call inherit.Init for %s", er.tp.String())
}

type iStruct interface {
	getType() reflect.Type
	setType(reflect.Type)
}

// Init the real type of obj for down casting.
func Init[T iStruct](imp T) T {
	imp.setType(reflect.TypeOf(imp).Elem())
	return imp
}

// Struct is the root type of the inheritance tree.
type Struct struct {
	realType reflect.Type // for down casting
}

func (b *Struct) getType() reflect.Type {
	return b.realType
}

func (b *Struct) setType(t reflect.Type) {
	b.realType = t
}

// Cast support both down & up casting.
func Cast[To any](s iStruct) (d *To) {
	dstTp := reflect.TypeOf(d).Elem()
	realTp := s.getType()

	if realTp == nil {
		panic(&ErrUnInit{reflect.TypeOf(s)})
	}

	if dstTp != realTp {

		// check cache

		cache.RLock()
		if v, ok1 := cache.casting[realTp]; ok1 {
			if v2, ok2 := v[dstTp]; ok2 && v2 {
				cache.RUnlock()
				if v2 {
					goto success
				} else {
					return nil
				}
			}
		}
		cache.RUnlock()

		// traverse the inheritance tree

		f := realTp.Field(0).Type
		for f != dstTp && f.Kind() == reflect.Struct {
			f = f.Field(0).Type
		}
		ok := f == dstTp

		// save to cache

		cache.Lock()
		if cache.casting[realTp] == nil {
			cache.casting[realTp] = make(map[reflect.Type]bool)
		}
		cache.casting[realTp][dstTp] = ok
		cache.Unlock()
		if !ok {
			return nil
		}
	}

success:
	return UnsafeCast[To](s)
}

func UnsafeCast[To any](s iStruct) (d *To) {
	d = (*To)((*[2]unsafe.Pointer)(unsafe.Pointer(&s))[1])
	return
}
