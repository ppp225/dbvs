package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ppp225/ndgo"
)

func (db Database) AddItem(item *Item) error {
	item.UID = "_:new"
	item.Type = []string{"Item"}
	t := time.Now().UTC()
	item.CreatedAt = &t
	txn := ndgo.NewTxn(context.TODO(), db.dg.NewTxn())
	resp, err := txn.Seti(item)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	item.UID = resp.GetUids()["new"]
	return nil
}

func (db Database) GetItemById(itemId string) (Item, error) {
	item := Item{}
	q := fmt.Sprintf(`{q(func: uid(%s)) @filter(eq(dgraph.type, Item)) { uid, dgraph.type name description createdAt }}`, itemId)
	resp, err := db.dg.NewTxn().Query(context.TODO(), q)
	if err != nil {
		return item, err
	}
	err = json.Unmarshal(ndgo.Unsafe{}.FlattenRespToObject(resp.GetJson()), &item)
	if err != nil {
		return item, err
	}
	if item.UID == "" {
		return item, ErrNotFound
	}
	return item, nil
}

func (db Database) Get1000Items() (*ItemList, error) {
	list := &ItemList{}
	q := `{items(func: type(Item), orderdesc: createdAt, first: 1000) { uid dgraph.type name description createdAt }}`
	resp, err := db.dg.NewTxn().Query(context.TODO(), q)
	if err != nil {
		return list, err
	}
	err = json.Unmarshal(resp.GetJson(), &list)
	if err != nil {
		return list, err
	}
	return list, nil
}
