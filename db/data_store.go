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

func (s *DataStore) Seed() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `
	CREATE TABLE data (
    id serial PRIMARY KEY,
    value VARCHAR(50));
	`

	_, err := s.Db.ExecContext(ctx, statement)
	if err != nil {
		return err
	}

	return nil
}

func (s *DataStore) Drop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `
	DROP table data
	`
}

func (s *DataStore) InsertData(entry string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	statement := `insert into data (value) values $1 returning id`

	err := s.Db.QueryRowContext(ctx, statement, entry).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *DataStore) UpdateData(entry Entry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `update data set value = $1 where id = $2`

	_, err := s.Db.ExecContext(ctx, statement, entry.Data, entry.Id)

	return err
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

func (s *DataStore) DeleteData(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data Entry
	query := `delete from data where id = $1`

	err := s.Db.QueryRowContext(ctx, query, id).Scan(data)
	return err

}
