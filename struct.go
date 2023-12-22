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

var typeCastingCache = struct {
	sync.RWMutex
	types map[reflect.Type]map[reflect.Type]bool
}{
	types: make(map[reflect.Type]map[reflect.Type]bool),
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
func (b *Struct) Call(f any, args ...any) []reflect.Value {
	typeInfo := b.typeInfo
	fv := reflect.ValueOf(f)
	fullName := runtime.FuncForPC(fv.Pointer()).Name()

	typeInfo.RLock()
	m, ok := typeInfo.fns[fullName]
	typeInfo.RUnlock()

	if !ok {
		tokens := strings.Split(fullName, ".")
		name := strings.Replace(tokens[len(tokens)-1], "-fm", "", 1)
		m, ok = typeInfo.realType.MethodByName(name)
		if !ok {
			m, ok = reflect.PtrTo(typeInfo.realType).MethodByName(name)
			if !ok {
				panic(fmt.Errorf("method not found %s in type %s", name, typeInfo.realType.String()))
			}
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

		// check typeCastingCache

		typeCastingCache.RLock()
		if v, ok1 := typeCastingCache.types[realTp]; ok1 {
			if v2, ok2 := v[dstTp]; ok2 && v2 {
				typeCastingCache.RUnlock()
				if v2 {
					goto success
				} else {
					return nil
				}
			}
		}
		typeCastingCache.RUnlock()

		// traverse the inheritance tree

		f := realTp.Field(0).Type
		for f != dstTp && f.Kind() == reflect.Struct {
			f = f.Field(0).Type
		}
		ok := f == dstTp

		// save to typeCastingCache

		typeCastingCache.Lock()
		if typeCastingCache.types[realTp] == nil {
			typeCastingCache.types[realTp] = make(map[reflect.Type]bool)
		}
		typeCastingCache.types[realTp][dstTp] = ok
		typeCastingCache.Unlock()
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
