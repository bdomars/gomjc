package microjava

type Scope struct {
	outer  *Scope
	locals *Obj
	nVars  int
	level  int
}
