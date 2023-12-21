//
// Single inheritance simulation for Golang.
//
// original created by lisong@legoutech.com
//

package inherit

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

// region global cache

var typeCache = struct {
	sync.RWMutex
	casting map[reflect.Type]map[reflect.Type]bool
}{
	casting: make(map[reflect.Type]map[reflect.Type]bool),
}

var fnCache = struct {
	sync.RWMutex
	fns map[string]reflect.Method
}{
	fns: map[string]reflect.Method{},
}

// endregion

// region Error

type ErrUnInit struct {
	tp reflect.Type
}

func (er ErrUnInit) Error() string {
	return fmt.Sprintf("missing call inherit.Init for %s", er.tp.String())
}

// endregion

type iStruct interface {
	getType() reflect.Type
	setType(reflect.Type)
	Call(f any, args ...any) []reflect.Value
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

// Call can invoke the methods of the concrete type.
func (b *Struct) Call(f any, args ...any) []reflect.Value {
	fv := reflect.ValueOf(f)
	fullName := runtime.FuncForPC(fv.Pointer()).Name()

	fnCache.RLock()
	m, ok := fnCache.fns[fullName]
	if !ok {
		fnCache.RUnlock()

		tokens := strings.Split(fullName, ".")
		name := strings.Replace(tokens[len(tokens)-1], "-fm", "", 1)
		m, ok = b.realType.MethodByName(name)
		if !ok {
			m, ok = reflect.PtrTo(b.realType).MethodByName(name)
			if !ok {
				panic(fmt.Errorf("method not found %s in type %s", name, b.realType.String()))
			}
		}

		fnCache.Lock()
		fnCache.fns[fullName] = m
		fnCache.Unlock()
	}

	vArgs := make([]reflect.Value, 0, len(args))
	vArgs = append(vArgs, reflect.NewAt(b.realType, unsafe.Pointer(b)))
	for _, i := range args {
		vArgs = append(vArgs, reflect.ValueOf(i))
	}
	return m.Func.Call(vArgs)
}

// To support both down & up casting.
func To[Dst any](s iStruct) (d *Dst) {
	dstTp := reflect.TypeOf(d).Elem()
	realTp := s.getType()

	if realTp == nil {
		panic(&ErrUnInit{reflect.TypeOf(s)})
	}

	if dstTp != realTp {

		// check typeCache

		typeCache.RLock()
		if v, ok1 := typeCache.casting[realTp]; ok1 {
			if v2, ok2 := v[dstTp]; ok2 && v2 {
				typeCache.RUnlock()
				if v2 {
					goto success
				} else {
					return nil
				}
			}
		}
		typeCache.RUnlock()

		// traverse the inheritance tree

		f := realTp.Field(0).Type
		for f != dstTp && f.Kind() == reflect.Struct {
			f = f.Field(0).Type
		}
		ok := f == dstTp

		// save to typeCache

		typeCache.Lock()
		if typeCache.casting[realTp] == nil {
			typeCache.casting[realTp] = make(map[reflect.Type]bool)
		}
		typeCache.casting[realTp][dstTp] = ok
		typeCache.Unlock()
		if !ok {
			return nil
		}
	}

success:
	return UnsafeCast[Dst](s)
}

func UnsafeCast[To any](s iStruct) (d *To) {
	d = (*To)((*[2]unsafe.Pointer)(unsafe.Pointer(&s))[1])
	return
}

// region type safe call wrappers

func CallArg0Ret0(obj iStruct, f func()) {
	obj.Call(f)
	return
}
func CallArg0Ret1[R0 any](obj iStruct, f func() R0) R0 {
	out := obj.Call(f)
	return out[0].Interface().(R0)
}
func CallArg0Ret2[R0, R1 any](obj iStruct, f func() (R0, R1)) (R0, R1) {
	out := obj.Call(f)
	return out[0].Interface().(R0), out[1].Interface().(R1)
}
func CallArg1Ret0[A0 any](obj iStruct, f func(A0), a0 A0) {
	obj.Call(f, a0)
	return
}
func CallArg1Ret1[R0, A0 any](obj iStruct, f func(A0) R0, a0 A0) R0 {
	out := obj.Call(f, a0)
	return out[0].Interface().(R0)
}
func CallArg1Ret2[R0, R1, A0 any](obj iStruct, f func(A0) (R0, R1), a0 A0) (R0, R1) {
	out := obj.Call(f, a0)
	return out[0].Interface().(R0), out[1].Interface().(R1)
}
func CallArg2Ret0[A0, A1 any](obj iStruct, f func(A0, A1), a0 A0, a1 A1) {
	obj.Call(f, a0, a1)
	return
}
func CallArg2Ret1[R0, A0, A1 any](obj iStruct, f func(A0, A1) R0, a0 A0, a1 A1) R0 {
	out := obj.Call(f, a0, a1)
	return out[0].Interface().(R0)
}
func CallArg2Ret2[R0, R1, A0, A1 any](obj iStruct, f func(A0, A1) (R0, R1), a0 A0, a1 A1) (R0, R1) {
	out := obj.Call(f, a0, a1)
	return out[0].Interface().(R0), out[1].Interface().(R1)
}
func CallArg3Ret0[A0, A1, A2 any](obj iStruct, f func(A0, A1, A2), a0 A0, a1 A1, a2 A2) {
	obj.Call(f, a0, a1, a2)
	return
}
func CallArg3Ret1[R0, A0, A1, A2 any](obj iStruct, f func(A0, A1, A2) R0, a0 A0, a1 A1, a2 A2) R0 {
	out := obj.Call(f, a0, a1, a2)
	return out[0].Interface().(R0)
}
func CallArg3Ret2[R0, R1, A0, A1, A2 any](obj iStruct, f func(A0, A1, A2) (R0, R1), a0 A0, a1 A1, a2 A2) (R0, R1) {
	out := obj.Call(f, a0, a1, a2)
	return out[0].Interface().(R0), out[1].Interface().(R1)
}
func CallArg4Ret0[A0, A1, A2, A3 any](obj iStruct, f func(A0, A1, A2, A3), a0 A0, a1 A1, a2 A2, a3 A3) {
	obj.Call(f, a0, a1, a2, a3)
	return
}
func CallArg4Ret1[R0, A0, A1, A2, A3 any](obj iStruct, f func(A0, A1, A2, A3) R0, a0 A0, a1 A1, a2 A2, a3 A3) R0 {
	out := obj.Call(f, a0, a1, a2, a3)
	return out[0].Interface().(R0)
}
func CallArg4Ret2[R0, R1, A0, A1, A2, A3 any](obj iStruct, f func(A0, A1, A2, A3) (R0, R1), a0 A0, a1 A1, a2 A2, a3 A3) (R0, R1) {
	out := obj.Call(f, a0, a1, a2, a3)
	return out[0].Interface().(R0), out[1].Interface().(R1)
}

// endregion
