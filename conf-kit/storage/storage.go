package storage

import (
	"database/sql"
	"strings"

	"github.com/ironzhang/golang/easysql"
)

const tableName = "tb_config"

type Storage struct {
	db *sql.DB
}

func (s *Storage) List(path string) ([]string, error) {
	id, err := s.getid(path)
	if err != nil {
		return nil, err
	}
	return s.list(id)
}

func (s *Storage) Add(path string) error {
	return nil
}

func (s *Storage) Remove(path string) error {
	return nil
}

func (s *Storage) Set(path string, value []byte) error {
	return nil
}

func (s *Storage) Get(path string) ([]byte, error) {
	return nil, nil
}

func (s *Storage) getid(path string) (int64, error) {
	names := strings.Split(path, "/")
	return s.lookupID(names)
}

func (s *Storage) lookupID(names []string) (int64, error) {
	var id int64
	var err error
	for _, name := range names {
		if id, err = s.queryID(id, name); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (s *Storage) queryID(pid int64, name string) (int64, error) {
	var id int64
	esql := easysql.SelectFrom(tableName).Column("id", &id).Where("pid=? and name=?", pid, name)
	query, args, err := esql.Query()
	if err != nil {
		return 0, err
	}
	vars := esql.Vars()
	if err = s.db.QueryRow(query, args...).Scan(vars...); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) list(pid int64) ([]string, error) {
	var name string
	esql := easysql.SelectFrom(tableName).Column("name", &name).Where("pid=?", pid)
	query, args, err := esql.Query()
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	vars := esql.Vars()
	names := make([]string, 0)
	for rows.Next() {
		if err = rows.Scan(vars...); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}
