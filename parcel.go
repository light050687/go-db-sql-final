package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("INSERT INTO parcel(client, status, address, created_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	res, err := stmt.Exec(p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	// Получение ID последней вставленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("SELECT client, status, address, created_at FROM parcel WHERE Number = ?")
	if err != nil {
		return Parcel{}, err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	p := Parcel{}
	err = stmt.QueryRow(number).Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	p.Number = number

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	rows, err := stmt.Query(client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	// Проверка на ошибки при выполнении запроса
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("UPDATE parcel SET status = ? WHERE number = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	_, err = stmt.Exec(status, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Проверка текущего статуса parcel
	p, err := s.Get(number)
	if err != nil {
		return err
	}

	if p.Status != ParcelStatusRegistered {
		return fmt.Errorf("Невозможно изменить адрес parcel, который не находится в статусе 'registered'")
	}

	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("UPDATE parcel SET address = ? WHERE number = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	_, err = stmt.Exec(address, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// Проверка текущего статуса parcel
	p, err := s.Get(number)
	if err != nil {
		return err
	}

	if p.Status != ParcelStatusRegistered {
		return fmt.Errorf("Невозможно удалить parcel, которая не находится в статусе 'registered'")
	}

	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("DELETE FROM parcel WHERE number = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	_, err = stmt.Exec(number)
	if err != nil {
		return err
	}

	return nil
}
