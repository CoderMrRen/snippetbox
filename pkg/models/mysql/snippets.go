package mysql

import (
	"database/sql"

	"helloren.cn/snippetbox/pkg/models"
)

//SnippetModel fdfsdf
type SnippetModel struct {
	DB *sql.DB
}

//Insert xxx
func (m *SnippetModel) Insert(title, content, expries string) (int, error) {
	stmt := "INSERT INTO snippets (title,content,created,expires) value (?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))"
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, title, content, expries)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return int(id), tx.Commit()
}

//Get xx
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := "SELECT id,title,content,created,expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"
	s := &models.Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

//Latest xx
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	//检测rows.Next() 是否有错误
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
