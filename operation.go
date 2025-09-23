package fscli

import (
	"time"

	"cloud.google.com/go/firestore"
)

type OperationType string

const (
	OPERATION_TYPE_QUERY OperationType = "QUERY"
	OPERATION_TYPE_GET   OperationType = "GET"
	OPERATION_TYPE_COUNT OperationType = "COUNT"
)

type Operation interface {
	Type() string
	OperationType() OperationType
	Collection() string
}

type BaseOperation struct {
}

func (op *BaseOperation) Type() string {
	return "Operation"
}

type Operator string

const (
	OPERATOR_EQ                 Operator = "=="
	OPERATOR_NOT_EQ             Operator = "!="
	OPERATOR_GT                 Operator = ">"
	OPERATOR_GTE                Operator = ">="
	OPERATOR_LT                 Operator = "<"
	OPERATOR_LTE                Operator = "<="
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

type TimestampFilter struct {
	BaseFilter
	value time.Time
}

func NewTimestampFilter(field string, operator Operator, value time.Time) *TimestampFilter {
	return &TimestampFilter{BaseFilter{field, operator}, value}
}

func (f *TimestampFilter) Value() any {
	return f.value
}

type OrderBy struct {
	field     string
	direction firestore.Direction
}

type QueryOperation struct {
	BaseOperation
	collection string
	selects    []string
	filters    []Filter
	orderBys   []OrderBy
	limit      int
}

func NewQueryOperation(collection string, selects []string, filters []Filter, orderBys []OrderBy, limit int) *QueryOperation {
	return &QueryOperation{collection: collection, selects: selects, filters: filters, orderBys: orderBys, limit: limit}
}

func (op *QueryOperation) OperationType() OperationType {
	return OPERATION_TYPE_QUERY
}

func (op *QueryOperation) Collection() string {
	return op.collection
}

type GetOperation struct {
	BaseOperation
	collection string
	docId      string
}

func NewGetOperation(collection string, docId string) *GetOperation {
	return &GetOperation{collection: collection, docId: docId}
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

type CountOperation struct {
	BaseOperation
	collection string
	filters    []Filter
}

func NewCountOperation(collection string, filters []Filter) *CountOperation {
	return &CountOperation{collection: collection, filters: filters}
}

func (op *CountOperation) OperationType() OperationType {
	return OPERATION_TYPE_COUNT
}

func (op *CountOperation) Collection() string {
	return op.collection
}
