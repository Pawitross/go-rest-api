package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"pawrest/internal/yamlconfig"
)

var (
	ErrNotFound   = errors.New("No resource found")
	ErrForeignKey = errors.New("Foreign key constraint error")
	ErrParam      = errors.New("Parameter error")
)

type DatabaseInterface interface {
	BookDatabaseInterface
	AuthorDatabaseInterface
	GenreDatabaseInterface
	LanguageDatabaseInterface
}

type Database struct {
	pool *sql.DB
}

var _ DatabaseInterface = (*Database)(nil)

func ConnectToDB(cfg *yamlconfig.Config) (*Database, error) {
	dbCfg := mysql.NewConfig()

	dbCfg.User = cfg.DBUser
	dbCfg.Passwd = cfg.DBPass
	dbCfg.Net = "tcp"
	dbCfg.Addr = cfg.DBHost + ":" + cfg.DBPort
	dbCfg.DBName = cfg.DBName
	dbCfg.ClientFoundRows = true

	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetConnMaxLifetime(4 * time.Minute)
	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(150)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{pool: db}, nil
}

func (d *Database) CloseDB() {
	if d.pool != nil {
		d.pool.Close()
	}
}

func (d *Database) Pool() *sql.DB {
	return d.pool
}

func queryWithParams[T any](
	d *Database,
	query string,
	params url.Values,
	allowPar map[string]string,
	scanFunc func(*T, *sql.Rows) error,
) ([]T, error) {
	var args []any

	if len(params) > 0 {
		filter, argsOut, err := AssembleFilter(params, allowPar)
		if err != nil {
			return nil, err
		}

		query += filter
		args = argsOut
	}

	records := []T{}

	rows, err := d.pool.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Query error (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r T

		if err := scanFunc(&r, rows); err != nil {
			return nil, fmt.Errorf("Scan error (%v)", err)
		}

		records = append(records, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Rows error (%v)", err)
	}

	return records, nil
}

func queryID[T any](
	d *Database,
	query string,
	id int64,
	scanFunc func(*T, *sql.Row) error,
) (T, error) {
	var r T

	row := d.pool.QueryRow(query, id)
	if err := scanFunc(&r, row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r, fmt.Errorf("%w with id %v", ErrNotFound, id)
		}

		return r, fmt.Errorf("Scan error (%v)", err)
	}

	return r, nil
}

func (d *Database) insert(query string, args ...any) (int64, error) {
	res, err := d.pool.Exec(query, args...)
	if err != nil {
		if isErrForeignKey(err) {
			return 0, ErrForeignKey
		}

		return 0, fmt.Errorf("Failed to insert record (%v)", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve id (%v)", err)
	}

	return id, nil
}

func (d *Database) updateWholeID(query string, args ...any) error {
	res, err := d.pool.Exec(query, args...)
	if err != nil {
		if isErrForeignKey(err) {
			return ErrForeignKey
		}

		return fmt.Errorf("Failed to update (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Rows affected error (%v)", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (d *Database) updatePartID(r any, table string, id int64, fToDB map[string]string) error {
	var (
		updates []string
		args    []any
	)

	valOfR := reflect.ValueOf(r)

	fields := reflect.VisibleFields(reflect.TypeOf(r))
	for i, f := range fields {
		fValue := valOfR.Field(i)
		if fValue.IsZero() {
			continue
		}

		columnName, ok := fToDB[f.Name]
		if !ok {
			continue
		}

		updates = append(updates, columnName+" = ?")
		args = append(args, fValue.Interface())
	}

	if len(updates) == 0 {
		return fmt.Errorf("No columns to update")
	}

	query := "UPDATE " + table + " SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	res, err := d.pool.Exec(query, args...)
	if err != nil {
		if isErrForeignKey(err) {
			return ErrForeignKey
		}

		return fmt.Errorf("Failed to update (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Rows affected error (%v)", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (d *Database) deleteID(query string, id int64) error {
	res, err := d.pool.Exec(query, id)
	if err != nil {
		if isErrForeignKey(err) {
			return ErrForeignKey
		}

		return fmt.Errorf("Failed to delete (%v)", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Rows affected error (%v)", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w with id %v", ErrNotFound, id)
	}

	return nil
}

func isErrForeignKey(err error) bool {
	var mysqlerr *mysql.MySQLError

	if errors.As(err, &mysqlerr) {
		return mysqlerr.Number == 1452 ||
			mysqlerr.Number == 1451
	}

	return false
}

func AssembleFilter(params url.Values, allowedParams map[string]string) (string, []any, error) {
	var (
		conditions []string
		args       []any
	)

	operators := map[string]string{
		".eq":  "=",
		".gt":  ">",
		".lt":  "<",
		".gte": ">=",
		".lte": "<=",
		".neq": "<>",
	}

	limit, hasLimit := params["limit"]
	offset, hasOffset := params["offset"]

	if !hasLimit && hasOffset {
		return "", nil, fmt.Errorf("%w: a limit must be provided when using an offset", ErrParam)
	}

	for key, valSlice := range params {
		if key == "limit" || key == "offset" || key == "sort_by" || key == "extend" {
			continue
		}

		operator := "="

		for suffix, sqlOp := range operators {
			if before, found := strings.CutSuffix(key, suffix); found {
				key = before
				operator = sqlOp
				break
			}
		}

		columnName, allowed := allowedParams[key]
		if !allowed {
			return "", nil, fmt.Errorf("%w: an unknown parameter was provided", ErrParam)
		}

		if len(valSlice) == 0 || valSlice[0] == "" {
			return "", nil, fmt.Errorf("%w: provided parameter is empty", ErrParam)
		}

		if len(valSlice) > 1 {
			return "", nil, fmt.Errorf("%w: too many parameters were provided for a single column", ErrParam)
		}

		if strings.ToLower(valSlice[0]) != "null" {
			conditions = append(conditions, columnName+" "+operator+" ?")
			args = append(args, valSlice[0])
		} else {
			if operator != "=" && operator != "<>" {
				return "", nil, fmt.Errorf("%w: cannot use other operations than equal or not equal on null", ErrParam)
			}

			if operator == "=" {
				operator = "IS"
			} else {
				operator = "IS NOT"
			}

			conditions = append(conditions, columnName+" "+operator+" NULL")
		}
	}

	filter := ""

	if len(conditions) > 0 {
		filter += " WHERE " + strings.Join(conditions, " AND ")
	}

	if sort, hasSort := params["sort_by"]; hasSort {
		if len(sort) == 0 || sort[0] == "" {
			return "", nil, fmt.Errorf("%w: provided column is empty", ErrParam)
		}

		if len(sort) > 1 {
			return "", nil, fmt.Errorf("%w: too many columns were provided for sorting", ErrParam)
		}

		order := ""

		arg := sort[0]
		if after, found := strings.CutPrefix(arg, "-"); found {
			order = " DESC"
			arg = after
		}

		columnName, allowed := allowedParams[arg]
		if !allowed {
			return "", nil, fmt.Errorf("%w: an unknown column was provided", ErrParam)
		}

		filter = filter + " ORDER BY " + columnName + order
	}

	if hasLimit {
		filter += " LIMIT ?"
		args = append(args, limit[0])
	}

	if hasOffset {
		filter += " OFFSET ?"
		args = append(args, offset[0])
	}

	return filter, args, nil
}
