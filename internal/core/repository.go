package core

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type (
	NotFoundError struct {
		Err   error
		Field string
		Value string
	}

	DuplicateKeyError struct {
		Err error
		Key string
	}

	ConstraintError struct {
		Err        error
		Constraint string
	}

	InvalidArgumentError struct {
		Err      error
		Argument string
	}
)

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("entity with %s %s not found: %v", e.Field, e.Value, e.Err)
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("duplicate key violation for %s: %v", e.Key, e.Err)
}

func (e *ConstraintError) Error() string {
	return fmt.Sprintf("constraint violation for %s: %v", e.Constraint, e.Err)
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument %s: %v", e.Argument, e.Err)
}

func (e *NotFoundError) Unwrap() error        { return e.Err }
func (e *DuplicateKeyError) Unwrap() error    { return e.Err }
func (e *ConstraintError) Unwrap() error      { return e.Err }
func (e *InvalidArgumentError) Unwrap() error { return e.Err }

func HandleError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return &DuplicateKeyError{Key: pgErr.ConstraintName, Err: err}
		case "23503": // foreign_key_violation
			return &ConstraintError{Constraint: pgErr.ConstraintName, Err: err}
		}
	}

	return err
}
