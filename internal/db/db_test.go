package db

import (
	"database/sql"
	"flag"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"pawrest/internal/yamlconfig"
)

var database *Database

type quote struct {
	ID      int64
	Quote   string
	Ranking int64
	FK      int64
}

func runMain(m *testing.M) (int, error) {
	flag.Parse()

	if testing.Short() {
		log.Println("Skipping Database testing in short mode because we don't have connection to the database")
		return 0, nil
	}

	cfg := &yamlconfig.Config{
		DBUser: "user_test",
		DBPass: "testpass",
		DBName: "paw_test",
		DBHost: "127.0.0.1",
		DBPort: "3306",
	}

	var err error
	database, err = ConnectToDB(cfg)
	if err != nil {
		return 0, err
	}
	defer database.CloseDB()

	if err := setupTestDatabase(database.Pool()); err != nil {
		return 0, err
	}

	defer database.Pool().Exec("DROP TABLE IF EXISTS test_fk")
	defer database.Pool().Exec("DROP TABLE IF EXISTS test_table")

	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestQueryID(t *testing.T) {
	tests := map[string]struct {
		giveID    int64
		wantQuote string
		wantRank  int64
		wantFK    int64
		wantErrIs error
	}{
		"Success": {
			giveID:    1,
			wantQuote: "Lorem ipsum dolor sit amet",
			wantRank:  3,
			wantFK:    1,
		},
		"ErrNotFound": {
			giveID:    1000,
			wantErrIs: ErrNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query := `
				SELECT
					id, quote, ranking, fk
				FROM test_table
				WHERE id = ?`

			quoteFunc := func(q *quote, row *sql.Row) error {
				return row.Scan(&q.ID, &q.Quote, &q.Ranking, &q.FK)
			}

			if tt.wantErrIs == nil {
				q, err := queryID[quote](database, query, tt.giveID, quoteFunc)
				assert.NoError(t, err)

				assert.Equal(t, q.Quote, tt.wantQuote)
				assert.Equal(t, q.Ranking, tt.wantRank)
				assert.Equal(t, q.FK, tt.wantFK)
			} else {
				_, err := queryID[quote](database, query, tt.giveID, quoteFunc)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErrIs)
			}
		})
	}
}

func TestQueryWithParams(t *testing.T) {
	tests := map[string]struct {
		giveParams url.Values
		wantErrIs  error
	}{
		"SuccessNoParams": {
			giveParams: url.Values{},
		},
		"SuccessIDParam": {
			giveParams: url.Values{"id": {"1"}},
		},
		"SuccessQuoteParam": {
			giveParams: url.Values{"quote": {"Lorem"}},
		},
		"SuccessRankingParam": {
			giveParams: url.Values{"ranking": {"1"}},
		},
		"SuccessIDParamSortDefault": {
			giveParams: url.Values{"sort_by": {"id"}},
		},
		"SuccessIDParamSortDescending": {
			giveParams: url.Values{"sort_by": {"-id"}},
		},
		"SuccessIDParamEq": {
			giveParams: url.Values{"id.eq": {"1"}},
		},
		"SuccessIDParamNotEq": {
			giveParams: url.Values{"id.neq": {"1"}},
		},
		"SuccessIDParamGt": {
			giveParams: url.Values{"id.gt": {"1"}},
		},
		"SuccessIDParamGte": {
			giveParams: url.Values{"id.gte": {"1"}},
		},
		"SuccessIDParamLt": {
			giveParams: url.Values{"id.lt": {"3"}},
		},
		"SuccessIDParamLte": {
			giveParams: url.Values{"id.lte": {"3"}},
		},
		"SuccessRankingParamNeqNull": {
			giveParams: url.Values{"ranking.neq": {"null"}},
		},
		"SuccessLimit": {
			giveParams: url.Values{"limit": {"10"}},
		},
		"SuccessOffset(Limit)": {
			giveParams: url.Values{"offset": {"2"}, "limit": {"10"}},
		},
		"ErrorUnknownParam": {
			giveParams: url.Values{"foo": {"bar"}},
			wantErrIs:  ErrParam,
		},
		"ErrorEmptyParam": {
			giveParams: url.Values{"id": {""}},
			wantErrIs:  ErrParam,
		},
		"ErrorMultipleSameParam": {
			giveParams: url.Values{"id": {"1", "2"}},
			wantErrIs:  ErrParam,
		},
		"ErrorGtOnNull": {
			giveParams: url.Values{"ranking.gt": {"null"}},
			wantErrIs:  ErrParam,
		},
		"ErrorEmptySortingParam": {
			giveParams: url.Values{"sort_by": {""}},
			wantErrIs:  ErrParam,
		},
		"ErrorUnknownSortingParam": {
			giveParams: url.Values{"sort_by": {"foo"}},
			wantErrIs:  ErrParam,
		},
		"ErrorMultipleSortingParam": {
			giveParams: url.Values{"sort_by": {"id", "quote"}},
			wantErrIs:  ErrParam,
		},
		"ErrorOffsetNoLimit": {
			giveParams: url.Values{"offset": {"2"}},
			wantErrIs:  ErrParam,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query := `
				SELECT
					id, quote, ranking
				FROM test_table`

			allowedParams := map[string]string{
				"id":      "id",
				"quote":   "quote",
				"ranking": "ranking",
			}

			quoteFunc := func(q *quote, rows *sql.Rows) error {
				return rows.Scan(&q.ID, &q.Quote, &q.Ranking)
			}

			if tt.wantErrIs == nil {
				qs, err := queryWithParams[quote](database, query, tt.giveParams, allowedParams, quoteFunc)
				assert.NoError(t, err)

				assert.NotEmpty(t, qs)
			} else {
				_, err := queryWithParams[quote](database, query, tt.giveParams, allowedParams, quoteFunc)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErrIs)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	tests := map[string]struct {
		giveQuote any
		giveRank  any
		giveFK    int64
		wantErr   bool
		wantErrIs error
	}{
		"Success": {
			giveQuote: "Test quote",
			giveRank:  int64(5),
			giveFK:    1,
		},
		"ErrStringRank": {
			giveQuote: "Test quote",
			giveRank:  "Should be a number",
			giveFK:    1,
			wantErr:   true,
		},
		"ErrForeignKey": {
			giveQuote: "Test quote",
			giveRank:  int64(5),
			giveFK:    1000,
			wantErr:   true,
			wantErrIs: ErrForeignKey,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query := `
				INSERT INTO test_table
				(quote, ranking, fk)
				VALUES (?, ?, ?)`

			if !tt.wantErr {
				id, err := database.insert(query, tt.giveQuote, tt.giveRank, tt.giveFK)
				assert.NoError(t, err)
				assert.NotEmpty(t, id)
			} else {
				_, err := database.insert(query, tt.giveQuote, tt.giveRank, tt.giveFK)
				assert.Error(t, err)

				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			}
		})
	}
}

func TestUpdateWholeID(t *testing.T) {
	tests := map[string]struct {
		giveQuote string
		giveRank  any
		giveFK    int64
		giveID    any
		wantErr   bool
		wantErrIs error
	}{
		"Success": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFK:    1,
			giveID:    int64(2),
		},
		"ErrStringRank": {
			giveQuote: "Update test quote",
			giveRank:  "string",
			giveFK:    1,
			giveID:    int64(2),
			wantErr:   true,
		},
		"ErrStringID": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFK:    1,
			giveID:    "string",
			wantErr:   true,
		},
		"ErrNotFound": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFK:    1,
			giveID:    int64(1000),
			wantErr:   true,
			wantErrIs: ErrNotFound,
		},
		"ErrForeignKey": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFK:    1000,
			giveID:    int64(1),
			wantErr:   true,
			wantErrIs: ErrForeignKey,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query := `
				UPDATE test_table
				SET
					quote = ?,
					ranking = ?,
					fk = ?
				WHERE id = ?`

			err := database.updateWholeID(query, tt.giveQuote, tt.giveRank, tt.giveFK, tt.giveID)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			}
		})
	}
}

func TestUpdatePartID(t *testing.T) {
	tests := map[string]struct {
		giveQuote quote
		giveID    int64
		wantErr   bool
		wantErrIs error
	}{
		"SuccessQuote": {
			giveQuote: quote{
				Quote: "Updating quote",
			},
			giveID: 1,
		},
		"SuccessRanking": {
			giveQuote: quote{
				Ranking: 100,
			},
			giveID: 1,
		},
		"SuccessFK": {
			giveQuote: quote{
				FK: 1,
			},
			giveID: 1,
		},
		"SuccessWholeObj": {
			giveQuote: quote{
				Quote:   "Updating quote",
				Ranking: 100,
				FK:      1,
			},
			giveID: 1,
		},
		"ErrEmptyObj": {
			giveQuote: quote{},
			giveID:    1,
			wantErr:   true,
		},
		"ErrNotFound": {
			giveQuote: quote{
				Quote: "Updating quote",
			},
			giveID:    1000,
			wantErr:   true,
			wantErrIs: ErrNotFound,
		},
		"ErrForeignKey": {
			giveQuote: quote{
				FK: 1000,
			},
			giveID:    1,
			wantErr:   true,
			wantErrIs: ErrForeignKey,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fieldToDB := map[string]string{
				"Quote":   "quote",
				"Ranking": "ranking",
				"FK":      "fk",
			}

			err := database.updatePartID(tt.giveQuote, "test_table", tt.giveID, fieldToDB)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := map[string]struct {
		id      int64
		wantErr error
	}{
		"Success": {
			id: 1,
		},
		"ErrNotFound": {
			id:      1000,
			wantErr: ErrNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query := "DELETE FROM test_table WHERE id = ?"

			err := database.deleteID(query, tt.id)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrNotFound)
			}
		})
	}
}

func setupTestDatabase(db *sql.DB) error {
	if _, err := db.Exec("DROP TABLE IF EXISTS test_table"); err != nil {
		return err
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS test_fk"); err != nil {
		return err
	}

	if _, err := db.Exec(`
		CREATE TABLE test_fk(
			id  INT AUTO_INCREMENT PRIMARY KEY,
			val VARCHAR(64)
		)
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`
		CREATE TABLE test_table(
			id      INT AUTO_INCREMENT,
			quote   VARCHAR(1024) NOT NULL,
			ranking INT NOT NULL,
			fk      INT NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY (fk) REFERENCES test_fk(id)
		)
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`
		INSERT INTO test_fk (val) VALUES
			("Val 1"), ("Val 1"), ("Val 3")
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`
		INSERT INTO test_table (quote, ranking, fk) VALUES
			("Lorem ipsum dolor sit amet", 3, 1),
			("Lorem", 2, 2),
			("ipsum", 1, 3)
	`); err != nil {
		return err
	}

	return nil
}
