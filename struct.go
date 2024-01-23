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

/*
Usage:
	1. embed inherit.Struct as value and the first member in the base class.
	2. call inherit.Init on the subclass instance.

Then you can:
	1. call Cast to cast upper/down.
	2. use inherit.InvokeX_X to call methods of subclass from base.
*/

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
	call(f any, args ...any) []reflect.Value
}

// Init the real type of obj for down casting.
func Init[T iStruct](imp T) T {
	imp.setType(reflect.TypeOf(imp).Elem())
	return imp
}

type TypeInfo struct {
	sync.RWMutex
	realType reflect.Type
	fns      map[string]reflect.Method
}

var typeInfoCache = struct {
	sync.RWMutex
	types map[reflect.Type]*TypeInfo
}{
	types: map[reflect.Type]*TypeInfo{},
}

// Struct is the root type of the inheritance tree.
type Struct struct {
	typeInfo *TypeInfo
}

func (b *Struct) getType() reflect.Type {
	if b.typeInfo == nil {
		panic(&ErrUnInit{reflect.TypeOf(b)})
	}
	return b.typeInfo.realType
}

func (b *Struct) setType(t reflect.Type) {
	typeInfoCache.RLock()
	info, ok := typeInfoCache.types[t]
	typeInfoCache.RUnlock()

	if ok {
		b.typeInfo = info
	} else {
		info = &TypeInfo{
			realType: t,
			fns:      map[string]reflect.Method{},
		}
		typeInfoCache.Lock()
		typeInfoCache.types[t] = info
		typeInfoCache.Unlock()
	}
	b.typeInfo = info
}

// Call can invoke the methods of the concrete type.
func (b *Struct) call(f any, args ...any) []reflect.Value {
	typeInfo := b.typeInfo
	fv := reflect.ValueOf(f)
	fullName := runtime.FuncForPC(fv.Pointer()).Name()

	typeInfo.RLock()
	m, ok := typeInfo.fns[fullName]
	typeInfo.RUnlock()

	if !ok {
		tokens := strings.Split(fullName, ".")
		name := strings.Replace(tokens[len(tokens)-1], "-fm", "", 1)
		m, ok = reflect.PtrTo(typeInfo.realType).MethodByName(name)
		if !ok {
			panic(fmt.Errorf("method not found %s in type %s", name, typeInfo.realType.String()))
		}

		typeInfo.Lock()
		typeInfo.fns[fullName] = m
		typeInfo.Unlock()
	}

	vArgs := make([]reflect.Value, 0, len(args))
	vArgs = append(vArgs, reflect.NewAt(typeInfo.realType, unsafe.Pointer(b)))
	for _, i := range args {
		vArgs = append(vArgs, reflect.ValueOf(i))
	}
	return m.Func.Call(vArgs)
}

// To support both down & up casting.
func To[Dst any](s iStruct) (d *Dst) {
	dstTp := reflect.TypeOf(d).Elem()
	realTp := s.getType()

	if dstTp != realTp {
		// traverse the inheritance tree
		f := realTp.Field(0).Type
		for f != dstTp && f.Kind() == reflect.Struct {
			f = f.Field(0).Type
		}
		if f != dstTp {
			return nil
		}
	}

	return UnsafeCast[Dst](s)
}

func UnsafeCast[To any](s iStruct) (d *To) {
	d = (*To)((*[2]unsafe.Pointer)(unsafe.Pointer(&s))[1])
	return
}

func checkedCast[T any](src reflect.Value) (r T) {
	if v := src.Interface(); v != nil {
		return v.(T)
	}
	return
}

// region type safe call wrappers

// nolint:lll
func Invoke0_0(f func(), obj iStruct) {
	obj.call(f)
	return
}

// nolint:lll
func Invoke0_1[R0 any](f func() R0, obj iStruct) R0 {
	r := obj.call(f)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke0_2[R0, R1 any](f func() (R0, R1), obj iStruct) (R0, R1) {
	r := obj.call(f)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke0_3[R0, R1, R2 any](f func() (R0, R1, R2), obj iStruct) (R0, R1, R2) {
	r := obj.call(f)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke0_4[R0, R1, R2, R3 any](f func() (R0, R1, R2, R3), obj iStruct) (R0, R1, R2, R3) {
	r := obj.call(f)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// nolint:lll
func Invoke1_0[A0 any](f func(A0), obj iStruct, a0 A0) {
	obj.call(f, a0)
	return
}

// nolint:lll
func Invoke1_1[R0, A0 any](f func(A0) R0, obj iStruct, a0 A0) R0 {
	r := obj.call(f, a0)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke1_2[R0, R1, A0 any](f func(A0) (R0, R1), obj iStruct, a0 A0) (R0, R1) {
	r := obj.call(f, a0)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke1_3[R0, R1, R2, A0 any](f func(A0) (R0, R1, R2), obj iStruct, a0 A0) (R0, R1, R2) {
	r := obj.call(f, a0)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke1_4[R0, R1, R2, R3, A0 any](f func(A0) (R0, R1, R2, R3), obj iStruct, a0 A0) (R0, R1, R2, R3) {
	r := obj.call(f, a0)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// nolint:lll
func Invoke2_0[A0, A1 any](f func(A0, A1), obj iStruct, a0 A0, a1 A1) {
	obj.call(f, a0, a1)
	return
}

// nolint:lll
func Invoke2_1[R0, A0, A1 any](f func(A0, A1) R0, obj iStruct, a0 A0, a1 A1) R0 {
	r := obj.call(f, a0, a1)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke2_2[R0, R1, A0, A1 any](f func(A0, A1) (R0, R1), obj iStruct, a0 A0, a1 A1) (R0, R1) {
	r := obj.call(f, a0, a1)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke2_3[R0, R1, R2, A0, A1 any](f func(A0, A1) (R0, R1, R2), obj iStruct, a0 A0, a1 A1) (R0, R1, R2) {
	r := obj.call(f, a0, a1)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke2_4[R0, R1, R2, R3, A0, A1 any](f func(A0, A1) (R0, R1, R2, R3), obj iStruct, a0 A0, a1 A1) (R0, R1, R2, R3) {
	r := obj.call(f, a0, a1)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// nolint:lll
func Invoke3_0[A0, A1, A2 any](f func(A0, A1, A2), obj iStruct, a0 A0, a1 A1, a2 A2) {
	obj.call(f, a0, a1, a2)
	return
}

// nolint:lll
func Invoke3_1[R0, A0, A1, A2 any](f func(A0, A1, A2) R0, obj iStruct, a0 A0, a1 A1, a2 A2) R0 {
	r := obj.call(f, a0, a1, a2)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke3_2[R0, R1, A0, A1, A2 any](f func(A0, A1, A2) (R0, R1), obj iStruct, a0 A0, a1 A1, a2 A2) (R0, R1) {
	r := obj.call(f, a0, a1, a2)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke3_3[R0, R1, R2, A0, A1, A2 any](f func(A0, A1, A2) (R0, R1, R2), obj iStruct, a0 A0, a1 A1, a2 A2) (R0, R1, R2) {
	r := obj.call(f, a0, a1, a2)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke3_4[R0, R1, R2, R3, A0, A1, A2 any](f func(A0, A1, A2) (R0, R1, R2, R3), obj iStruct, a0 A0, a1 A1, a2 A2) (R0, R1, R2, R3) {
	r := obj.call(f, a0, a1, a2)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// nolint:lll
func Invoke4_0[A0, A1, A2, A3 any](f func(A0, A1, A2, A3), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3) {
	obj.call(f, a0, a1, a2, a3)
	return
}

// nolint:lll
func Invoke4_1[R0, A0, A1, A2, A3 any](f func(A0, A1, A2, A3) R0, obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3) R0 {
	r := obj.call(f, a0, a1, a2, a3)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke4_2[R0, R1, A0, A1, A2, A3 any](f func(A0, A1, A2, A3) (R0, R1), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3) (R0, R1) {
	r := obj.call(f, a0, a1, a2, a3)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke4_3[R0, R1, R2, A0, A1, A2, A3 any](f func(A0, A1, A2, A3) (R0, R1, R2), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3) (R0, R1, R2) {
	r := obj.call(f, a0, a1, a2, a3)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke4_4[R0, R1, R2, R3, A0, A1, A2, A3 any](f func(A0, A1, A2, A3) (R0, R1, R2, R3), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3) (R0, R1, R2, R3) {
	r := obj.call(f, a0, a1, a2, a3)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// nolint:lll
func Invoke5_0[A0, A1, A2, A3, A4 any](f func(A0, A1, A2, A3, A4), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4) {
	obj.call(f, a0, a1, a2, a3, a4)
	return
}

// nolint:lll
func Invoke5_1[R0, A0, A1, A2, A3, A4 any](f func(A0, A1, A2, A3, A4) R0, obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4) R0 {
	r := obj.call(f, a0, a1, a2, a3, a4)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke5_2[R0, R1, A0, A1, A2, A3, A4 any](f func(A0, A1, A2, A3, A4) (R0, R1), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4) (R0, R1) {
	r := obj.call(f, a0, a1, a2, a3, a4)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke5_3[R0, R1, R2, A0, A1, A2, A3, A4 any](f func(A0, A1, A2, A3, A4) (R0, R1, R2), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4) (R0, R1, R2) {
	r := obj.call(f, a0, a1, a2, a3, a4)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke5_4[R0, R1, R2, R3, A0, A1, A2, A3, A4 any](f func(A0, A1, A2, A3, A4) (R0, R1, R2, R3), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4) (R0, R1, R2, R3) {
	r := obj.call(f, a0, a1, a2, a3, a4)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// nolint:lll
func Invoke6_0[A0, A1, A2, A3, A4, A5 any](f func(A0, A1, A2, A3, A4, A5), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5) {
	obj.call(f, a0, a1, a2, a3, a4, a5)
	return
}

// nolint:lll
func Invoke6_1[R0, A0, A1, A2, A3, A4, A5 any](f func(A0, A1, A2, A3, A4, A5) R0, obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5) R0 {
	r := obj.call(f, a0, a1, a2, a3, a4, a5)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke6_2[R0, R1, A0, A1, A2, A3, A4, A5 any](f func(A0, A1, A2, A3, A4, A5) (R0, R1), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5) (R0, R1) {
	r := obj.call(f, a0, a1, a2, a3, a4, a5)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke6_3[R0, R1, R2, A0, A1, A2, A3, A4, A5 any](f func(A0, A1, A2, A3, A4, A5) (R0, R1, R2), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5) (R0, R1, R2) {
	r := obj.call(f, a0, a1, a2, a3, a4, a5)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke6_4[R0, R1, R2, R3, A0, A1, A2, A3, A4, A5 any](f func(A0, A1, A2, A3, A4, A5) (R0, R1, R2, R3), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5) (R0, R1, R2, R3) {
	r := obj.call(f, a0, a1, a2, a3, a4, a5)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// nolint:lll
func Invoke7_0[A0, A1, A2, A3, A4, A5, A6 any](f func(A0, A1, A2, A3, A4, A5, A6), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6) {
	obj.call(f, a0, a1, a2, a3, a4, a5, a6)
	return
}

// nolint:lll
func Invoke7_1[R0, A0, A1, A2, A3, A4, A5, A6 any](f func(A0, A1, A2, A3, A4, A5, A6) R0, obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6) R0 {
	r := obj.call(f, a0, a1, a2, a3, a4, a5, a6)
	return checkedCast[R0](r[0])
}

// nolint:lll
func Invoke7_2[R0, R1, A0, A1, A2, A3, A4, A5, A6 any](f func(A0, A1, A2, A3, A4, A5, A6) (R0, R1), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6) (R0, R1) {
	r := obj.call(f, a0, a1, a2, a3, a4, a5, a6)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1])
}

// nolint:lll
func Invoke7_3[R0, R1, R2, A0, A1, A2, A3, A4, A5, A6 any](f func(A0, A1, A2, A3, A4, A5, A6) (R0, R1, R2), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6) (R0, R1, R2) {
	r := obj.call(f, a0, a1, a2, a3, a4, a5, a6)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2])
}

// nolint:lll
func Invoke7_4[R0, R1, R2, R3, A0, A1, A2, A3, A4, A5, A6 any](f func(A0, A1, A2, A3, A4, A5, A6) (R0, R1, R2, R3), obj iStruct, a0 A0, a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6) (R0, R1, R2, R3) {
	r := obj.call(f, a0, a1, a2, a3, a4, a5, a6)
	return checkedCast[R0](r[0]), checkedCast[R1](r[1]), checkedCast[R2](r[2]), checkedCast[R3](r[3])
}

// endregion
