package pgerrors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Полный список: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	// Class 22 Data Exception
	SQLStateInvalidTextRepresentation = "22P02" // invalid input syntax for type

	// Class 23 Integrity Constraint Violation
	SQLStateNotNullViolation    = "23502"
	SQLStateForeignKeyViolation = "23503"
	SQLStateUniqueViolation     = "23505"
	SQLStateCheckViolation      = "23514"
	SQLStateExclusionViolation  = "23P01"

	// Class 40 Transaction Rollback
	SQLStateSerializationFailure = "40001"
	SQLStateDeadlockDetected     = "40P01"

	// Class 42 Syntax Error or Access Rule Violation
	SQLStateUndefinedTable  = "42P01"
	SQLStateUndefinedColumn = "42703"
)

var (
	ErrUniqueViolation           = errors.New("db: unique violation")
	ErrForeignKeyViolation       = errors.New("db: foreign key violation")
	ErrNotNullViolation          = errors.New("db: not null violation")
	ErrCheckViolation            = errors.New("db: check violation")
	ErrExclusionViolation        = errors.New("db: exclusion violation")
	ErrSerializationFailure      = errors.New("db: serialization failure")
	ErrDeadlockDetected          = errors.New("db: deadlock detected")
	ErrInvalidTextRepresentation = errors.New("db: invalid text representation")
	ErrUndefinedTable            = errors.New("db: undefined table")
	ErrUndefinedColumn           = errors.New("db: undefined column")
)

var codeToError = map[string]error{
	SQLStateUniqueViolation:           ErrUniqueViolation,
	SQLStateForeignKeyViolation:       ErrForeignKeyViolation,
	SQLStateNotNullViolation:          ErrNotNullViolation,
	SQLStateCheckViolation:            ErrCheckViolation,
	SQLStateExclusionViolation:        ErrExclusionViolation,
	SQLStateSerializationFailure:      ErrSerializationFailure,
	SQLStateDeadlockDetected:          ErrDeadlockDetected,
	SQLStateInvalidTextRepresentation: ErrInvalidTextRepresentation,
	SQLStateUndefinedTable:            ErrUndefinedTable,
	SQLStateUndefinedColumn:           ErrUndefinedColumn,
}

type Details struct {
	SQLState      string
	Severity      string
	Message       string
	Schema        string
	Table         string
	Column        string
	Constraint    string
	DataTypeName  string
	Detail        string
	Hint          string
	Where         string
	InternalQuery string
}

// Extract возвращает *pgconn.PgError, если err — это pgx-ошибка.
func Extract(err error) (*pgconn.PgError, bool) {
	var pe *pgconn.PgError

	if errors.As(err, &pe) {
		return pe, true
	}

	return nil, false
}

// Code извлекает SQLSTATE.
func Code(err error) (string, bool) {
	pe, ok := Extract(err)
	if ok {
		return pe.SQLState(), true
	}

	return "", false
}

// Class возвращает класс SQLSTATE (первые 2 символа), например "23".
func Class(err error) (string, bool) {
	code, ok := Code(err)
	if ok && len(code) >= 2 {
		return code[:2], true
	}

	return "", false
}

// Map маппит pgx-ошибку в одну из sentinel-ошибок.
// Если код неизвестен или это не pgx-ошибка — вернёт исходную err.
func Map(err error) error {
	if err == nil {
		return nil
	}

	code, ok := Code(err)
	if ok {
		if mapped, ok := codeToError[code]; ok {
			return mapped
		}
	}

	return err
}

// DetailsOf возвращает структурированные детали для логов.
func DetailsOf(err error) (Details, bool) {
	pe, ok := Extract(err)
	if !ok {
		return Details{}, false
	}

	return Details{
		SQLState:      pe.SQLState(),
		Severity:      pe.Severity,
		Message:       pe.Message,
		Schema:        pe.SchemaName,
		Table:         pe.TableName,
		Column:        pe.ColumnName,
		Constraint:    pe.ConstraintName,
		DataTypeName:  pe.DataTypeName,
		Detail:        pe.Detail,
		Hint:          pe.Hint,
		Where:         pe.Where,
		InternalQuery: pe.InternalQuery,
	}, true
}

func IsUniqueViolation(err error) bool      { return errors.Is(Map(err), ErrUniqueViolation) }
func IsForeignKeyViolation(err error) bool  { return errors.Is(Map(err), ErrForeignKeyViolation) }
func IsNotNullViolation(err error) bool     { return errors.Is(Map(err), ErrNotNullViolation) }
func IsCheckViolation(err error) bool       { return errors.Is(Map(err), ErrCheckViolation) }
func IsExclusionViolation(err error) bool   { return errors.Is(Map(err), ErrExclusionViolation) }
func IsSerializationFailure(err error) bool { return errors.Is(Map(err), ErrSerializationFailure) }
func IsDeadlockDetected(err error) bool     { return errors.Is(Map(err), ErrDeadlockDetected) }
func IsInvalidText(err error) bool          { return errors.Is(Map(err), ErrInvalidTextRepresentation) }
func IsUndefinedTable(err error) bool       { return errors.Is(Map(err), ErrUndefinedTable) }
func IsUndefinedColumn(err error) bool      { return errors.Is(Map(err), ErrUndefinedColumn) }
