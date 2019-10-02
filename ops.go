package corezoid

type Op interface {
	Raw() map[string]interface{}
}

type MapOp map[string]interface{}

type Ops struct {
	Ops []Op
}

func (ops *Ops) Add(op Op) {
	ops.Ops = append(ops.Ops, op)
}

func (ops *Ops) Raw() map[string]interface{} {
	var rawOps []map[string]interface{}

	for _, op := range ops.Ops {
		rawOps = append(rawOps, op.Raw())
	}

	return map[string]interface{}{
		"ops": rawOps,
	}
}

func (op MapOp) Raw() map[string]interface{} {
	return op
}
