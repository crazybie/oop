# Single inheritance simulation for Golang.

## Why need this?

### The Biggest pros over interface
1. no interface declaration for dynamic method call.
   - this will save a lot of tedious interface declaration code
   - save the effort to maintain signature match between interface and implementation.
2. objects can cast in inherit tree smoothly without interface.
   - because only interface has the inheritance trait in Golang, 
casting pointers must be performed through interfaces, 
and one field must have one pair of matching get&set methods in the corresponding interface, 
which cause a lot of boring interface code, even we just need a simple up-casting.

### Cons
1. overload checking is performed at runtime.
   - some IDE(e.g. Goland) can show the function shadowing, makes jumping between base and subclass possible.
   - the shadowing icon can be used to verify the overloading correction.
2. not performant as interface.
   - for invoking: method name lookup cost, arguments boxing overhead.
   - for casting: cost to inheritance tree walking.

## Usage
1. embed inherit.Struct as value and make it the first member in the base class.
2. call inherit.Init on the subclass instance.

#### Then you can
1. call `inherit.To[DestType]` to cast from base to subclass.
2. use `inherit.InvokeX_X` to call methods of subclass from base. (methods must be *Public* accessible)

## Examples

- struct casting.
```go
type Base struct {
  inherit.Struct //----------- step 1
}

type Sub struct {
  Base //----------- step 2

  val int	
}

func NewSub() *Sub {
  r := &Sub{val: 11}
  inherit.Init(r) //----------- step 3
  return r
}

func main() {
  var base *Base = NewSub()	
	
  // casting
  var sub *Sub = inherit.To[Sub](base)
  // sub.val == 11
}
```

- dynamic method call
```go
type Base struct {
  inherit.Struct //----------- step 1 
}

func (b *Base) Foo() int {
  panic("overwrite me")
}

type Sub struct {
  Base //----------- step 2 
}

func NewSub() *Sub {
  r := &Sub{}
  inherit.Init(r) //----------- step 3
  return r
}

func (s *Sub) Foo() int {
  return 11
}

func main(){
  var base *Base = NewSub()	
  
  // call the subclass version
  ret := inherit.Invoke0_1(base.Foo, base)
  // ret == 11
}
```