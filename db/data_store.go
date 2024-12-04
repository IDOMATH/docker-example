package db

import (
	"context"
	"database/sql"
	"time"
)

type DataStore struct {
	Db *sql.DB
}

type Entry struct {
	Id   int
	Data string
}

func NewDataStore(db *sql.DB) *DataStore {
	return &DataStore{
		Db: db,
	}
}

func (s *DataStore) InsertData(entry Entry) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	statement := `insert into data $1`

	err := s.Db.QueryRowContext(ctx, statement, entry.Data).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *DataStore) GetAllData() ([]Entry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data []Entry
	query := `select * from data`

	rows, err := s.Db.QueryContext(ctx, query)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var datum Entry
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

func (s *DataStore) GetDataById(id int) (Entry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data Entry
	query := `select * from data where id = $1`

	err := s.Db.QueryRowContext(ctx, query, id).Scan(data)
	if err != nil {
		return data, err
	}

	return data, nil
}
