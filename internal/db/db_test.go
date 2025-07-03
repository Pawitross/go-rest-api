package db

import (
	"database/sql"
	"flag"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var database *Database

type quote struct {
	Id      int64
	Quote   string
	Ranking int64
	Fk      int64
}

func runMain(m *testing.M) (int, error) {
	flag.Parse()

	if testing.Short() {
		log.Println("Skipping Database testing in short mode because we don't have connection to the database")
		return 0, nil
	}

	os.Setenv("DBUSER", "root")
	os.Setenv("DBPASS", "")
	os.Setenv("DBNAME", "paw_test")

	var err error
	database, err = ConnectToDB()
	if err != nil {
		return 0, err
	}
	defer database.CloseDB()

	if err := setupTestDatabase(database.Pool()); err != nil {
		return 0, err
	}

	return m.Run(), nil
}

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestQueryId(t *testing.T) {
	tests := map[string]struct {
		giveId    int64
		wantQuote string
		wantRank  int64
		wantFk    int64
		wantErrIs error
	}{
		"Success": {
			giveId:    1,
			wantQuote: "Lorem ipsum dolor sit amet",
			wantRank:  3,
			wantFk:    1,
		},
		"ErrNotFound": {
			giveId:    1000,
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
				return row.Scan(&q.Id, &q.Quote, &q.Ranking, &q.Fk)
			}

			if tt.wantErrIs == nil {
				q, err := queryId[quote](database, query, tt.giveId, quoteFunc)
				assert.NoError(t, err)

				assert.Equal(t, q.Quote, tt.wantQuote)
				assert.Equal(t, q.Ranking, tt.wantRank)
				assert.Equal(t, q.Fk, tt.wantFk)
			} else {
				_, err := queryId[quote](database, query, tt.giveId, quoteFunc)
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
		"SuccessIdParam": {
			giveParams: url.Values{"id": {"1"}},
		},
		"SuccessQuoteParam": {
			giveParams: url.Values{"quote": {"Lorem"}},
		},
		"SuccessRankingParam": {
			giveParams: url.Values{"ranking": {"1"}},
		},
		"SuccessIdParamSortDefault": {
			giveParams: url.Values{"sort_by": {"id"}},
		},
		"SuccessIdParamSortDescending": {
			giveParams: url.Values{"sort_by": {"-id"}},
		},
		"SuccessIdParamEq": {
			giveParams: url.Values{"id.eq": {"1"}},
		},
		"SuccessIdParamNotEq": {
			giveParams: url.Values{"id.neq": {"1"}},
		},
		"SuccessIdParamGt": {
			giveParams: url.Values{"id.gt": {"1"}},
		},
		"SuccessIdParamGte": {
			giveParams: url.Values{"id.gte": {"1"}},
		},
		"SuccessIdParamLt": {
			giveParams: url.Values{"id.lt": {"3"}},
		},
		"SuccessIdParamLte": {
			giveParams: url.Values{"id.lte": {"3"}},
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
		"ErrorEmptySortingParam": {
			giveParams: url.Values{"sort_by": {}},
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
				return rows.Scan(&q.Id, &q.Quote, &q.Ranking)
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
		giveFk    int64
		wantErr   bool
		wantErrIs error
	}{
		"Success": {
			giveQuote: "Test quote",
			giveRank:  int64(5),
			giveFk:    1,
		},
		"ErrStringRank": {
			giveQuote: "Test quote",
			giveRank:  "Should be a number",
			giveFk:    1,
			wantErr:   true,
		},
		"ErrForeignKey": {
			giveQuote: "Test quote",
			giveRank:  int64(5),
			giveFk:    1000,
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
				id, err := database.insert(query, tt.giveQuote, tt.giveRank, tt.giveFk)
				assert.NoError(t, err)
				assert.NotEmpty(t, id)
			} else {
				_, err := database.insert(query, tt.giveQuote, tt.giveRank, tt.giveFk)
				assert.Error(t, err)

				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			}
		})
	}
}

func TestUpdateWholeId(t *testing.T) {
	tests := map[string]struct {
		giveQuote string
		giveRank  any
		giveFk    int64
		giveId    any
		wantErr   bool
		wantErrIs error
	}{
		"Success": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFk:    1,
			giveId:    int64(2),
		},
		"ErrStringRank": {
			giveQuote: "Update test quote",
			giveRank:  "string",
			giveFk:    1,
			giveId:    int64(2),
			wantErr:   true,
		},
		"ErrStringId": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFk:    1,
			giveId:    "string",
			wantErr:   true,
		},
		"ErrNotFound": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFk:    1,
			giveId:    int64(1000),
			wantErr:   true,
			wantErrIs: ErrNotFound,
		},
		"ErrForeignKey": {
			giveQuote: "Update test quote",
			giveRank:  int64(10),
			giveFk:    1000,
			giveId:    int64(1),
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

			err := database.updateWholeId(query, tt.giveQuote, tt.giveRank, tt.giveFk, tt.giveId)
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

func TestUpdatePartId(t *testing.T) {
	tests := map[string]struct {
		giveQuote quote
		giveId    int64
		wantErr   bool
		wantErrIs error
	}{
		"SuccessQuote": {
			giveQuote: quote{
				Quote: "Updating quote",
			},
			giveId: 1,
		},
		"SuccessRanking": {
			giveQuote: quote{
				Ranking: 100,
			},
			giveId: 1,
		},
		"SuccessFk": {
			giveQuote: quote{
				Fk: 1,
			},
			giveId: 1,
		},
		"SuccessWholeObj": {
			giveQuote: quote{
				Quote:   "Updating quote",
				Ranking: 100,
				Fk:      1,
			},
			giveId: 1,
		},
		"ErrEmptyObj": {
			giveQuote: quote{},
			giveId:    1,
			wantErr:   true,
		},
		"ErrNotFound": {
			giveQuote: quote{
				Quote: "Updating quote",
			},
			giveId:    1000,
			wantErr:   true,
			wantErrIs: ErrNotFound,
		},
		"ErrForeignKey": {
			giveQuote: quote{
				Fk: 1000,
			},
			giveId:    1,
			wantErr:   true,
			wantErrIs: ErrForeignKey,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fieldToDb := map[string]string{
				"Quote":   "quote",
				"Ranking": "ranking",
				"Fk":      "fk",
			}

			err := database.updatePartId(tt.giveQuote, "test_table", tt.giveId, fieldToDb)
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

			err := database.deleteId(query, tt.id)
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
