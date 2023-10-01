package fscli

type Operation interface {
	OperationType() string
	Collection() string
}

type Filter interface {
	FieldName() string
	Operator() string
	Value() any
}

type IntFilter struct {
	field    string
	operator string
	value    int
}

func (f *IntFilter) FieldName() string {
	return f.field
}
func (f *IntFilter) Operator() string {
	return f.operator
}
func (f *IntFilter) Value() any {
	return f.value
}

type StringFilter struct {
	field    string
	operator string
	value    string
}

func (f *StringFilter) FieldName() string {
	return f.field
}
func (f *StringFilter) Operator() string {
	return f.operator
}
func (f *StringFilter) Value() any {
	return f.value
}

type QueryOperation struct {
	opType     string
	collection string
	filters    []Filter
}

func NewQueryOperation(collection string, filters []Filter) *QueryOperation {
	return &QueryOperation{QUERY, collection, filters}
}

func (op *QueryOperation) OperationType() string {
	return op.opType
}

func (op *QueryOperation) Collection() string {
	return op.collection
}
