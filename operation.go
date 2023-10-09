package fscli

type OperationType string

const (
	OPERATION_TYPE_QUERY OperationType = "QUERY"
	OPERATION_TYPE_GET   OperationType = "GET"
)

type Operation interface {
	OperationType() OperationType
	Collection() string
}

type Operator string

const (
	OPERATOR_EQ                 Operator = "=="
	OPERATOR_NOT_EQ             Operator = "!="
	OPERATOR_IN                 Operator = "in"
	OPERATOR_ARRAY_CONTAINS     Operator = "array-contains"
	OPERATOR_ARRAY_CONTAINS_ANY Operator = "array-contains-any"
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
	collection string
	selects    []string
	filters    []Filter
}

func NewQueryOperation(collection string, selects []string, filters []Filter) *QueryOperation {
	return &QueryOperation{collection, selects, filters}
}

func (op *QueryOperation) OperationType() OperationType {
	return OPERATION_TYPE_QUERY
}

func (op *QueryOperation) Collection() string {
	return op.collection
}

type GetOperation struct {
	collection string
	docId      string
}

func NewGetOperation(collection string, docId string) *GetOperation {
	return &GetOperation{collection, docId}
}

func (op *GetOperation) OperationType() OperationType {
	return OPERATION_TYPE_GET
}

func (op *GetOperation) Collection() string {
	return op.collection
}

func (op *GetOperation) DocId() string {
	return op.docId
}
