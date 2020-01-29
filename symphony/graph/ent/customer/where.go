// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package customer

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.Customer {
	return predicate.Customer(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.In(s.C(FieldID), v...))
	},
	)
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	},
	)
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	},
	)
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	},
	)
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	},
	)
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	},
	)
}

// ExternalID applies equality check predicate on the "external_id" field. It's identical to ExternalIDEQ.
func ExternalID(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExternalID), v))
	},
	)
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreateTime), v...))
	},
	)
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreateTime), v...))
	},
	)
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	},
	)
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdateTime), v...))
	},
	)
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdateTime), v...))
	},
	)
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	},
	)
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	},
	)
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	},
	)
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	},
	)
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	},
	)
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	},
	)
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	},
	)
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	},
	)
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	},
	)
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	},
	)
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	},
	)
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	},
	)
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	},
	)
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	},
	)
}

// ExternalIDEQ applies the EQ predicate on the "external_id" field.
func ExternalIDEQ(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDNEQ applies the NEQ predicate on the "external_id" field.
func ExternalIDNEQ(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDIn applies the In predicate on the "external_id" field.
func ExternalIDIn(vs ...string) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldExternalID), v...))
	},
	)
}

// ExternalIDNotIn applies the NotIn predicate on the "external_id" field.
func ExternalIDNotIn(vs ...string) predicate.Customer {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Customer(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldExternalID), v...))
	},
	)
}

// ExternalIDGT applies the GT predicate on the "external_id" field.
func ExternalIDGT(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDGTE applies the GTE predicate on the "external_id" field.
func ExternalIDGTE(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDLT applies the LT predicate on the "external_id" field.
func ExternalIDLT(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDLTE applies the LTE predicate on the "external_id" field.
func ExternalIDLTE(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDContains applies the Contains predicate on the "external_id" field.
func ExternalIDContains(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDHasPrefix applies the HasPrefix predicate on the "external_id" field.
func ExternalIDHasPrefix(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDHasSuffix applies the HasSuffix predicate on the "external_id" field.
func ExternalIDHasSuffix(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDIsNil applies the IsNil predicate on the "external_id" field.
func ExternalIDIsNil() predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldExternalID)))
	},
	)
}

// ExternalIDNotNil applies the NotNil predicate on the "external_id" field.
func ExternalIDNotNil() predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldExternalID)))
	},
	)
}

// ExternalIDEqualFold applies the EqualFold predicate on the "external_id" field.
func ExternalIDEqualFold(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldExternalID), v))
	},
	)
}

// ExternalIDContainsFold applies the ContainsFold predicate on the "external_id" field.
func ExternalIDContainsFold(v string) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldExternalID), v))
	},
	)
}

// HasServices applies the HasEdge predicate on the "services" edge.
func HasServices() predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServicesTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ServicesTable, ServicesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	},
	)
}

// HasServicesWith applies the HasEdge predicate on the "services" edge with a given conditions (other predicates).
func HasServicesWith(preds ...predicate.Service) predicate.Customer {
	return predicate.Customer(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServicesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ServicesTable, ServicesPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Customer) predicate.Customer {
	return predicate.Customer(
		func(s *sql.Selector) {
			s1 := s.Clone().SetP(nil)
			for _, p := range predicates {
				p(s1)
			}
			s.Where(s1.P())
		},
	)
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Customer) predicate.Customer {
	return predicate.Customer(
		func(s *sql.Selector) {
			s1 := s.Clone().SetP(nil)
			for i, p := range predicates {
				if i > 0 {
					s1.Or()
				}
				p(s1)
			}
			s.Where(s1.P())
		},
	)
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Customer) predicate.Customer {
	return predicate.Customer(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
