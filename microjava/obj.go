package microjava

type ObjectType int

const (
	objConstant ObjectType = iota
	objVariable
	objType
	objMethod
	objProgram
)

type Obj struct {
	kind    ObjectType
	name    string
	mjType  MJStruct
	val     int
	addr    int
	level   int
	nParams int
	locals  *Obj
	next    *Obj
}

func NewObj(kind ObjectType, name string, mjType MJStruct) *Obj {
	return &Obj{
		kind:   kind,
		name:   name,
		mjType: mjType,
	}
}
