package dbController

import (
	"database/sql"
	"fmt"
	"log"
	"wbservice/config"

	_ "github.com/lib/pq"
)

const getQuery = "SELECT %s FROM %s WHERE %s = '%s'"
const addQuery = "INSERT INTO %s VALUES ('%s', '%s')"
const initQuery = "CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(255), %s JSON);"
const batchQuery = "SELECT %s, %s FROM %s LIMIT %d OFFSET %d;"
const deleteQuery = "DELETE FROM %s WHERE %s = '%s'"

// const initQuery = "CREATE TABLE	if not exists %s (%s varchar(255) not null unique, %s json not null);"

type DBController struct {
	db         *sql.DB
	tableName  string
	idColumn   string
	dataColumn string
}

func New(psqlConfig config.PostgresConfig) (DBController, error) {
	address := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", psqlConfig.User, psqlConfig.Password, psqlConfig.URL, psqlConfig.Port, psqlConfig.Db)
	result := new(DBController)
	var err error
	result.db, err = sql.Open("postgres", address)
	if err != nil {
		return *result, err
	}

	err = result.db.Ping()

	if err != nil {
		return *result, err
	}

	result.tableName = psqlConfig.TableName
	result.idColumn = psqlConfig.IDColumn
	result.dataColumn = psqlConfig.DataColumn

	result.db.Exec(fmt.Sprintf(initQuery, result.tableName, result.idColumn, result.dataColumn))

	return *result, err
}

func (dbc *DBController) LogToDB(order_uid string, object string) error {
	result, err := dbc.db.Exec(fmt.Sprintf(addQuery, dbc.tableName, order_uid, object))
	log.Println(result)
	return err
}

func (dbc *DBController) GetFromDB(order_uid string) (string, error) {
	query := fmt.Sprintf(getQuery, dbc.dataColumn, dbc.tableName, dbc.idColumn, order_uid)
	row := dbc.db.QueryRow(query)
	var res string
	err := row.Scan(&res)
	return res, err
}

func (dbc *DBController) LoadBatch(limit int, offstet int) (*[][2]string, bool, error) {
	rows, err := dbc.db.Query(fmt.Sprintf(batchQuery, dbc.idColumn, dbc.dataColumn, dbc.tableName, limit, offstet))

	if err != nil {
		return nil, false, err
	}

	var order_uid string
	var order_data string
	var result = make([][2]string, 0)
	i := 0
	for rows.Next() {
		rows.Scan(&order_uid, &order_data)
		result = append(result, [2]string{order_uid, order_data})
		i++
	}

	return &result, len(result) == 100, err
}

func (dbc *DBController) RemoveKey(order_uid string) (sql.Result, error) {
	return dbc.db.Exec(fmt.Sprintf(deleteQuery, dbc.tableName, dbc.idColumn, order_uid))
}

func (dbc *DBController) Close() error {
	return dbc.db.Close()
}
