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

type BaseFilter struct {
	field    string
	operator string
}

func (f *BaseFilter) FieldName() string {
	return f.field
}
func (f *BaseFilter) Operator() string {
	return f.operator
}

type IntFilter struct {
	BaseFilter
	value int
}

func NewIntFilter(field, operator string, value int) *IntFilter {
	return &IntFilter{BaseFilter{field, operator}, value}
}

func (f *IntFilter) Value() any {
	return f.value
}

type FloatFilter struct {
	BaseFilter
	value float64
}

func NewFloatFilter(field, operator string, value float64) *FloatFilter {
	return &FloatFilter{BaseFilter{field, operator}, value}
}

func (f *FloatFilter) Value() any {
	return f.value
}

type StringFilter struct {
	BaseFilter
	value string
}

func NewStringFilter(field, operator, value string) *StringFilter {
	return &StringFilter{BaseFilter{field, operator}, value}
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
