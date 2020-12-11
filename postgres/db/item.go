package db

import "database/sql"

func (db Database) AddItem(item *Item) error {
	var id int
	var createdAt string
	query := `INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, item.Name, item.Description).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	item.ID = id
	item.CreatedAt = createdAt
	return nil
}

func (db Database) GetItemById(itemId int) (Item, error) {
	item := Item{}
	query := `SELECT * FROM items WHERE id = $1;`
	row := db.Conn.QueryRow(query, itemId)
	switch err := row.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt); err {
	case sql.ErrNoRows:
		return item, ErrNotFound
	default:
		return item, err
	}
}

func (db Database) DeleteItem(itemId int) error {
	query := `DELETE FROM items WHERE id = $1;`
	_, err := db.Conn.Exec(query, itemId)
	switch err {
	case sql.ErrNoRows:
		return ErrNotFound
	default:
		return err
	}
}

func (db Database) UpdateItem(itemId int, itemData Item) (Item, error) {
	item := Item{}
	query := `UPDATE items SET name=$1, description=$2 WHERE id=$3 RETURNING *;`
	err := db.Conn.QueryRow(query, itemData.Name, itemData.Description, itemId).Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, ErrNotFound
		}
		return item, err
	}
	return item, nil
}

func (db Database) Get1000Items() (*ItemList, error) {
	list := &ItemList{}
	rows, err := db.Conn.Query("SELECT * FROM items ORDER BY created_at DESC LIMIT 1000")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Items = append(list.Items, item)
	}
	return list, nil
}
