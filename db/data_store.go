package db

import (
	"context"
	"database/sql"
	"time"
)

type DataStore struct {
	Db *sql.DB
}

func NewDataStore(db *sql.DB) *DataStore {
	return &DataStore{
		Db: db,
	}
}

func (s *DataStore) InsertData(data string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	statement := `insert into data $1`

	err := s.Db.QueryRowContext(ctx, statement, data).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *DataStore) GetAllData() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data []string
	query := `select * from data`

	rows, err := s.Db.QueryContext(ctx, query)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var datum string
		err := rows.Scan(&datum)
		if err != nil {
			return data, err
		}
		data = append(data, datum)
	}

	if err = rows.Err(); err != nil {
		return data, err
	}
	return data, nil
}
