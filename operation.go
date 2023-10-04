package fscli

type Operation interface {
	OperationType() string
	Collection() string
}

type Operator string

const (
	OP_EQ     Operator = "=="
	OP_NOT_EQ Operator = "!="
	OP_IN     Operator = "in"
)

type Filter interface {
	FieldName() string
	Operator() Operator
	Value() any
}

type BaseFilter struct {
	field    string
	operator Operator
}

func (f *BaseFilter) FieldName() string {
	return f.field
}
func (f *BaseFilter) Operator() Operator {
	return f.operator
}

type IntFilter struct {
	BaseFilter
	value int
}

func NewIntFilter(field string, operator Operator, value int) *IntFilter {
	return &IntFilter{BaseFilter{field, operator}, value}
}

func (f *IntFilter) Value() any {
	return f.value
}

type FloatFilter struct {
	BaseFilter
	value float64
}

func NewFloatFilter(field string, operator Operator, value float64) *FloatFilter {
	return &FloatFilter{BaseFilter{field, operator}, value}
}

func (f *FloatFilter) Value() any {
	return f.value
}

type StringFilter struct {
	BaseFilter
	value string
}

func NewStringFilter(field string, operator Operator, value string) *StringFilter {
	return &StringFilter{BaseFilter{field, operator}, value}
}

func (f *StringFilter) Value() any {
	return f.value
}

type ArrayFilter struct {
	BaseFilter
	value []any
}

func NewArrayFilter(field string, operator Operator, value []any) *ArrayFilter {
	return &ArrayFilter{BaseFilter{field, operator}, value}
}

func (f *ArrayFilter) Value() any {
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
