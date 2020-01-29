// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"strings"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
)

// Order applies an ordering on either graph traversal or sql selector.
type Order func(*sql.Selector)

// Asc applies the given fields in ASC order.
func Asc(fields ...string) Order {
	return func(s *sql.Selector) {
		for _, f := range fields {
			s.OrderBy(sql.Asc(f))
		}
	}
}

// Desc applies the given fields in DESC order.
func Desc(fields ...string) Order {
	return func(s *sql.Selector) {
		for _, f := range fields {
			s.OrderBy(sql.Desc(f))
		}
	}
}

// Aggregate applies an aggregation step on the group-by traversal/selector.
type Aggregate func(*sql.Selector) string

// As is a pseudo aggregation function for renaming another other functions with custom names. For example:
//
//	GroupBy(field1, field2).
//	Aggregate(ent.As(ent.Sum(field1), "sum_field1"), (ent.As(ent.Sum(field2), "sum_field2")).
//	Scan(ctx, &v)
//
func As(fn Aggregate, end string) Aggregate {
	return func(s *sql.Selector) string {
		return sql.As(fn(s), end)
	}
}

// Count applies the "count" aggregation function on each group.
func Count() Aggregate {
	return func(s *sql.Selector) string {
		return sql.Count("*")
	}
}

// Max applies the "max" aggregation function on the given field of each group.
func Max(field string) Aggregate {
	return func(s *sql.Selector) string {
		return sql.Max(s.C(field))
	}
}

// Mean applies the "mean" aggregation function on the given field of each group.
func Mean(field string) Aggregate {
	return func(s *sql.Selector) string {
		return sql.Avg(s.C(field))
	}
}

// Min applies the "min" aggregation function on the given field of each group.
func Min(field string) Aggregate {
	return func(s *sql.Selector) string {
		return sql.Min(s.C(field))
	}
}

// Sum applies the "sum" aggregation function on the given field of each group.
func Sum(field string) Aggregate {
	return func(s *sql.Selector) string {
		return sql.Sum(s.C(field))
	}
}

// NotFoundError returns when trying to fetch a specific entity and it was not found in the database.
type NotFoundError struct {
	label string
}

// Error implements the error interface.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("ent: %s not found", e.label)
}

// IsNotFound returns a boolean indicating whether the error is a not found error.
func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// MaskNotFound masks nor found error.
func MaskNotFound(err error) error {
	if IsNotFound(err) {
		return nil
	}
	return err
}

// NotSingularError returns when trying to fetch a singular entity and more then one was found in the database.
type NotSingularError struct {
	label string
}

// Error implements the error interface.
func (e *NotSingularError) Error() string {
	return fmt.Sprintf("ent: %s not singular", e.label)
}

// IsNotSingular returns a boolean indicating whether the error is a not singular error.
func IsNotSingular(err error) bool {
	_, ok := err.(*NotSingularError)
	return ok
}

// NotLoadedError returns when trying to get a node that was not loaded by the query.
type NotLoadedError struct {
	label string
}

// Error implements the error interface.
func (e *NotLoadedError) Error() string {
	return fmt.Sprintf("ent: %s not loaded", e.label)
}

// IsNotLoaded returns a boolean indicating whether the error is a not loaded error.
func IsNotLoaded(err error) bool {
	_, ok := err.(*NotLoadedError)
	return ok
}

// ConstraintError returns when trying to create/update one or more entities and
// one or more of their constraints failed. For example, violation of edge or
// field uniqueness.
type ConstraintError struct {
	msg  string
	wrap error
}

// Error implements the error interface.
func (e ConstraintError) Error() string {
	return fmt.Sprintf("ent: constraint failed: %s", e.msg)
}

// Unwrap implements the errors.Wrapper interface.
func (e *ConstraintError) Unwrap() error {
	return e.wrap
}

// IsConstraintError returns a boolean indicating whether the error is a constraint failure.
func IsConstraintError(err error) bool {
	_, ok := err.(*ConstraintError)
	return ok
}

func isSQLConstraintError(err error) (*ConstraintError, bool) {
	var (
		msg = err.Error()
		// error format per dialect.
		errors = [...]string{
			"Error 1062",               // MySQL 1062 error (ER_DUP_ENTRY).
			"UNIQUE constraint failed", // SQLite.
			"duplicate key value violates unique constraint", // PostgreSQL.
		}
	)
	if _, ok := err.(*sqlgraph.ConstraintError); ok {
		return &ConstraintError{msg, err}, true
	}
	for i := range errors {
		if strings.Contains(msg, errors[i]) {
			return &ConstraintError{msg, err}, true
		}
	}
	return nil, false
}

// rollback calls to tx.Rollback and wraps the given error with the rollback error if occurred.
func rollback(tx dialect.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%s: %v", err.Error(), rerr)
	}
	if err, ok := isSQLConstraintError(err); ok {
		return err
	}
	return err
}

// insertLastID invokes the insert query on the transaction and returns the LastInsertID.
func insertLastID(ctx context.Context, tx dialect.Tx, insert *sql.InsertBuilder) (int64, error) {
	query, args := insert.Query()
	// PostgreSQL does not support the LastInsertId() method of sql.Result
	// on Exec, and should be extracted manually using the `RETURNING` clause.
	if insert.Dialect() == dialect.Postgres {
		rows := &sql.Rows{}
		if err := tx.Query(ctx, query, args, rows); err != nil {
			return 0, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, fmt.Errorf("no rows found for query: %v", query)
		}
		var id int64
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		return id, nil
	}
	// MySQL, SQLite, etc.
	var res sql.Result
	if err := tx.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// keys returns the keys/ids from the edge map.
func keys(m map[string]struct{}) []string {
	s := make([]string, 0, len(m))
	for id := range m {
		s = append(s, id)
	}
	return s
}
