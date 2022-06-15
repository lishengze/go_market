package modelext

import (
	"fmt"
	"strings"
	"time"
)

const (
	formatTime = "2006-01-02 15:04:05"
)

func NewOrderBy(field string, order Order) *OrderBy {
	return &OrderBy{
		field: field,
		order: order,
	}
}

type OrderBy struct {
	field string
	order Order
}

func (o *OrderBy) SqlContent() string {
	switch o.order {
	case ASC:
		return fmt.Sprintf("%s asc", o.field)
	case DESC:
		return fmt.Sprintf("%s desc", o.field)
	default:
		return fmt.Sprintf("%s asc", o.field)
	}
}

func (o *OrderBy) CacheKeyContent() string {
	switch o.order {
	case ASC:
		return fmt.Sprintf("%s asc", o.field)
	case DESC:
		return fmt.Sprintf("%s desc", o.field)
	default:
		return fmt.Sprintf("%s asc", o.field)
	}
}

type Order string

const (
	DESC Order = "desc"
	ASC  Order = "asc"
)

type WhereValue interface {
	Where() string
	Value() []interface{}
	CacheKeyContent() string
	Fields() []string
}

func NewEqualValue(field string, value interface{}) WhereValue {
	return &EqualValue{
		field: field,
		value: value,
	}
}

type EqualValue struct {
	field string
	value interface{}
}

func (o *EqualValue) Where() string {
	return fmt.Sprintf("`%s` = ?", o.field)
}

func (o *EqualValue) Value() []interface{} {
	return []interface{}{o.value}
}

func (o *EqualValue) CacheKeyContent() string {
	return fmt.Sprintf("%s=%v", o.field, o.value)
}

func (o *EqualValue) Fields() []string {
	return []string{o.field}
}

func NewNotEqualValue(field string, value interface{}) WhereValue {
	return &NotEqualValue{
		field: field,
		value: value,
	}
}

type NotEqualValue struct {
	field string
	value interface{}
}

func (o *NotEqualValue) Where() string {
	return fmt.Sprintf("`%s` != ?", o.field)
}

func (o *NotEqualValue) Value() []interface{} {
	return []interface{}{o.value}
}

func (o *NotEqualValue) CacheKeyContent() string {
	return fmt.Sprintf("%s!=%v", o.field, o.value)
}

func (o *NotEqualValue) Fields() []string {
	return []string{o.field}
}

func NewContainsValue(field string, value interface{}) WhereValue {
	return &ContainsValue{
		field: field,
		value: value,
	}
}

type ContainsValue struct {
	field string
	value interface{}
}

func (o *ContainsValue) Where() string {
	return fmt.Sprintf("`%s` like ?", o.field)
}

func (o *ContainsValue) Value() []interface{} {
	return []interface{}{fmt.Sprintf("%%%v%%", o.value)}
	//return []interface{}{o.value}
}

func (o *ContainsValue) CacheKeyContent() string {
	return fmt.Sprintf("%s contains %v", o.field, o.value)
}

func (o *ContainsValue) Fields() []string {
	return []string{o.field}
}

func NewInValue(field string, values []interface{}) WhereValue {
	return &InValue{
		field:  field,
		values: values,
	}
}

type InValue struct {
	field  string
	values []interface{}
}

func (o *InValue) Where() string {
	var placeholders []string

	if len(o.values) == 0 {
		placeholders = append(placeholders, "?")
	}

	for range o.values {
		placeholders = append(placeholders, "?")
	}
	return fmt.Sprintf("`%s` in (%s)", o.field, strings.Join(placeholders, ","))
}

func (o *InValue) Value() []interface{} {
	if len(o.values) == 0 {
		return []interface{}{""}
	}
	return o.values
}

func (o *InValue) CacheKeyContent() string {
	var values []string
	for _, v := range o.values {
		values = append(values, fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("%s in (%s)", o.field, strings.Join(values, ","))
}

func (o *InValue) Fields() []string {
	return []string{o.field}
}

func NewTimeRangeValue(field string, start, end time.Time) WhereValue {
	return &TimeRangeValue{
		field: field,
		start: start,
		end:   end,
	}
}

type TimeRangeValue struct {
	field string
	start time.Time
	end   time.Time
}

func (o *TimeRangeValue) Where() string {
	return fmt.Sprintf("(`%s` >= ? and `%s` <= ?)", o.field, o.field)
}

func (o *TimeRangeValue) Value() []interface{} {
	return []interface{}{o.start, o.end}
}

func (o *TimeRangeValue) CacheKeyContent() string {
	return fmt.Sprintf("%s<=%s<=%s",
		o.start.In(o.start.Location()).Format(formatTime),
		o.field,
		o.end.In(o.start.Location()).Format(formatTime))
}

func (o *TimeRangeValue) Fields() []string {
	return []string{o.field}
}

func NewTimeGeValue(field string, t time.Time) WhereValue {
	return &TimeGeValue{
		field: field,
		t:     t,
	}
}

type TimeGeValue struct {
	field string
	t     time.Time
}

func (o *TimeGeValue) Where() string {
	return fmt.Sprintf("`%s` >= ?", o.field)
}

func (o *TimeGeValue) Value() []interface{} {
	return []interface{}{o.t}
}

func (o *TimeGeValue) CacheKeyContent() string {
	return fmt.Sprintf("%s>=%s",
		o.field,
		o.t.In(o.t.Location()).Format(formatTime))
}

func (o *TimeGeValue) Fields() []string {
	return []string{o.field}
}

func NewTimeLeValue(field string, t time.Time) WhereValue {
	return &TimeLeValue{
		field: field,
		t:     t,
	}
}

type TimeLeValue struct {
	field string
	t     time.Time
}

func (o *TimeLeValue) Where() string {
	return fmt.Sprintf("`%s` <=?", o.field)
}

func (o *TimeLeValue) Value() []interface{} {
	return []interface{}{o.t}
}

func (o *TimeLeValue) CacheKeyContent() string {
	return fmt.Sprintf("%s<=%s",
		o.field,
		o.t.In(o.t.Location()).Format(formatTime))
}

func (o *TimeLeValue) Fields() []string {
	return []string{o.field}
}

func NewOrValue(whereValues ...WhereValue) WhereValue {
	return &OrValue{
		Wheres: whereValues,
	}
}

type OrValue struct {
	Wheres
}

func (o OrValue) Where() string {
	var (
		s      string
		wheres []string
	)
	if len(o.Wheres) == 0 {
		return s
	}

	for _, w := range o.Wheres {
		wheres = append(wheres, w.Where())
	}

	return fmt.Sprintf("(%s)", strings.Join(wheres, " or "))
}

func (o *OrValue) Value() []interface{} {
	var res = make([]interface{}, 0)
	for _, w := range o.Wheres {
		res = append(res, w.Value()...)
	}
	return res
}

func (o *OrValue) CacheKeyContent() string {
	var (
		s      string
		wheres []string
	)
	if len(o.Wheres) == 0 {
		return s
	}

	for _, w := range o.Wheres {
		wheres = append(wheres, w.CacheKeyContent())
	}

	return fmt.Sprintf("(%s)", strings.Join(wheres, " or "))
}

func (o *OrValue) Fields() []string {
	var (
		fields []string
	)

	for _, w := range o.Wheres {
		fields = append(fields, w.Fields()...)
	}
	return fields
}

func NewPaginator(page, size int64) *Paginator {
	return &Paginator{
		page: page,
		size: size,
	}
}

type Paginator struct {
	page  int64
	size  int64
	count int64
}

func NewQuery() *Query {
	return &Query{
		hasPaginator: false,
		paginator:    NewPaginator(0, 0),
		Wheres:       make(Wheres, 0),
		orders:       make([]*OrderBy, 0),
	}
}

type Query struct {
	hasPaginator bool
	limit        int64
	paginator    *Paginator
	Wheres
	orders []*OrderBy
}

type Wheres []WhereValue

func (q *Query) Page() int64 {
	return q.paginator.page
}

func (q *Query) Size() int64 {
	return q.paginator.size
}

func (q *Query) SetCount(count int64) {
	q.paginator.count = count
}

func (q *Query) Limit(limit int64) *Query {
	q.limit = limit
	return q
}

func (q *Query) SetPaginator(page, size int64) *Query {
	q.hasPaginator = true
	q.paginator = NewPaginator(page, size)
	return q
}

func (q *Query) HasPaginator() bool {
	return q.hasPaginator
}

func (q *Query) Count() int64 {
	switch q.paginator {
	case nil:
		return 0
	default:
		return q.paginator.count
	}
}

func (q *Query) Where(whereValues ...WhereValue) *Query {
	q.Wheres.Where(whereValues...)
	return q
}

func (q *Query) Or(fn func(*OrValue)) *Query {
	q.Wheres.Or(fn)
	return q
}

func (q *Query) WhereFunc(fn func(wheres *Wheres)) *Query {
	q.Wheres.WhereFunc(fn)
	return q
}

func (q *Query) OrderBy(field string, order Order) *Query {
	q.orders = append(q.orders, NewOrderBy(field, order))
	return q
}

func (q *Query) Contains(field string, value interface{}) *Query {
	q.Wheres.Contains(field, value)
	return q
}

func (q *Query) In(field string, value []interface{}) *Query {
	q.Wheres.In(field, value)
	return q
}

func (q *Query) Equal(field string, value interface{}) *Query {
	q.Wheres.Equal(field, value)
	return q
}

func (q *Query) NotEqual(field string, value interface{}) *Query {
	q.Wheres.NotEqual(field, value)
	return q
}

func (q *Query) TimeRange(field string, start, end time.Time) *Query {
	q.Wheres.TimeRange(field, start, end)
	return q
}

func (q *Query) TimeGe(field string, t time.Time) *Query {
	q.Wheres.TimeGe(field, t)
	return q
}

func (q *Query) TimeLe(field string, t time.Time) *Query {
	q.Wheres.TimeLe(field, t)
	return q
}

func (w *Wheres) append(whereValues ...WhereValue) {
	*w = append(*w, whereValues...)
}

func (w *Wheres) Where(whereValues ...WhereValue) {
	w.append(whereValues...)
}

func (w *Wheres) Or(fn func(*OrValue)) {
	orValue := &OrValue{}
	fn(orValue)
	if len(orValue.Wheres) > 0 {
		w.append(orValue)
	}
}

func (w *Wheres) WhereFunc(fn func(wheres *Wheres)) {
	fn(w)
}

func (w *Wheres) Contains(field string, value interface{}) {
	w.append(NewContainsValue(field, value))
}

func (w *Wheres) In(field string, value []interface{}) {
	w.append(NewInValue(field, value))
}

func (w *Wheres) Equal(field string, value interface{}) {
	w.append(NewEqualValue(field, value))
}

func (w *Wheres) NotEqual(field string, value interface{}) {
	w.append(NewNotEqualValue(field, value))
}

func (w *Wheres) TimeGe(field string, t time.Time) {
	w.append(NewTimeGeValue(field, t))
}

func (w *Wheres) TimeLe(field string, t time.Time) {
	w.append(NewTimeLeValue(field, t))
}

func (w *Wheres) TimeRange(field string, start, end time.Time) {
	w.append(NewTimeRangeValue(field, start, end))
}

//func (q *Query) Sql(rows, table string) (string, string, []interface{}) {
//	var (
//		findManySql = fmt.Sprintf("select %s from %s $where$ $order$ $limit$offset$", rows, table)
//		countSql    = fmt.Sprintf("select count(*) from %s $where$", table)
//		wheres      []string
//		args        []interface{}
//		orders      []string
//	)
//
//	for _, where := range q.Wheres {
//		wheres = append(wheres, where.Where())
//		args = append(args, where.Value()...)
//	}
//	for _, order := range q.orders {
//		orders = append(orders, order.SqlContent())
//	}
//
//	if len(wheres) > 0 {
//		findManySql = strings.Replace(findManySql, "$where$", "where "+strings.Join(wheres, " and "), -1)
//		countSql = strings.Replace(countSql, "$where$", "where "+strings.Join(wheres, " and "), -1)
//	} else {
//		findManySql = strings.Replace(findManySql, "$where$", "", -1)
//		countSql = strings.Replace(countSql, "$where$", "", -1)
//	}
//	if len(orders) > 0 {
//		findManySql = strings.Replace(findManySql, "$order$", "order by "+strings.Join(orders, ", "), -1)
//	} else {
//		findManySql = strings.Replace(findManySql, "$order$", "", -1)
//	}
//
//	if q.Paginator != nil {
//		limit := q.Paginator.Size
//		offset := q.Paginator.Size * (q.Paginator.Page - 1)
//		findManySql = strings.Replace(findManySql, "$limit$offset$", fmt.Sprintf("limit %d offset %d", limit, offset), -1)
//	} else {
//		if q.limit > 0 {
//			findManySql = strings.Replace(findManySql, "$limit$offset$", fmt.Sprintf("limit %d", q.limit), -1)
//		} else {
//			findManySql = strings.Replace(findManySql, "$limit$offset$", "", -1)
//		}
//	}
//	return findManySql, countSql, args
//}

//func (q *Query) FindManySql(rows, table string) (string, string, []interface{}) {
//	var (
//		findManySql = fmt.Sprintf("select %s from %s $where$ $order$ $limit$offset$", rows, table)
//		countSql    = fmt.Sprintf("select count(*) from %s $where$", table)
//		wheres      []string
//		args        []interface{}
//		orders      []string
//	)
//
//	for _, where := range q.Wheres {
//		wheres = append(wheres, where.Where())
//		args = append(args, where.Value()...)
//	}
//	for _, order := range q.orders {
//		orders = append(orders, order.SqlContent())
//	}
//
//	if len(wheres) > 0 {
//		findManySql = strings.Replace(findManySql, "$where$", "where "+strings.Join(wheres, " and "), -1)
//		countSql = strings.Replace(countSql, "$where$", "where "+strings.Join(wheres, " and "), -1)
//	} else {
//		findManySql = strings.Replace(findManySql, "$where$", "", -1)
//		countSql = strings.Replace(countSql, "$where$", "", -1)
//	}
//	if len(orders) > 0 {
//		findManySql = strings.Replace(findManySql, "$order$", "order by "+strings.Join(orders, ", "), -1)
//	} else {
//		findManySql = strings.Replace(findManySql, "$order$", "", -1)
//	}
//
//	if q.Paginator != nil {
//		limit := q.Paginator.Size
//		offset := q.Paginator.Size * (q.Paginator.Page - 1)
//		findManySql = strings.Replace(findManySql, "$limit$offset$", fmt.Sprintf("limit %d offset %d", limit, offset), -1)
//	} else {
//		if q.limit > 0 {
//			findManySql = strings.Replace(findManySql, "$limit$offset$", fmt.Sprintf("limit %d", q.limit), -1)
//		} else {
//			findManySql = strings.Replace(findManySql, "$limit$offset$", "", -1)
//		}
//	}
//	return findManySql, countSql, args
//}

func (q *Query) FindPksSql(pk, table string) (string, string, []interface{}) {
	var (
		findManySql = fmt.Sprintf("select %s from %s $where$ $order$ $limit$offset$", pk, table)
		countSql    = fmt.Sprintf("select count(*) from %s $where$", table)
		wheres      []string
		args        []interface{}
		orders      []string
	)

	for _, where := range q.Wheres {
		wheres = append(wheres, where.Where())
		args = append(args, where.Value()...)
	}
	for _, order := range q.orders {
		orders = append(orders, order.SqlContent())
	}

	if len(wheres) > 0 {
		findManySql = strings.Replace(findManySql, "$where$", "where "+strings.Join(wheres, " and "), -1)
		countSql = strings.Replace(countSql, "$where$", "where "+strings.Join(wheres, " and "), -1)
	} else {
		findManySql = strings.Replace(findManySql, "$where$", "", -1)
		countSql = strings.Replace(countSql, "$where$", "", -1)
	}
	if len(orders) > 0 {
		findManySql = strings.Replace(findManySql, "$order$", "order by "+strings.Join(orders, ", "), -1)
	} else {
		findManySql = strings.Replace(findManySql, "$order$", "", -1)
	}

	if q.HasPaginator() {
		limit := q.Size()
		offset := q.Size() * (q.Page() - 1)
		findManySql = strings.Replace(findManySql, "$limit$offset$", fmt.Sprintf("limit %d offset %d", limit, offset), -1)
	} else {
		if q.limit > 0 {
			findManySql = strings.Replace(findManySql, "$limit$offset$", fmt.Sprintf("limit %d", q.limit), -1)
		} else {
			findManySql = strings.Replace(findManySql, "$limit$offset$", "", -1)
		}
	}
	return findManySql, countSql, args
}

func (q *Query) Validate(fieldMap map[string]struct{}) error {
	for _, where := range q.Wheres {
		for _, filed := range where.Fields() {
			if _, ok := fieldMap[filed]; !ok {
				return fmt.Errorf("field %s is invalid. ", filed)
			}
		}

	}
	for _, by := range q.orders {
		if _, ok := fieldMap[by.field]; !ok {
			return fmt.Errorf("field %s is invalid. ", by.field)
		}
	}
	return nil
}
