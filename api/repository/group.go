package repository

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger"
	ddpportal "gitlab.adeo.com/ddp-portal-api"
)

type GroupRepository interface {
	Insert(group ddpportal.Group) error
	Update(group ddpportal.Group) error
	Delete(group ddpportal.Group) error
	GetAll() ([]*ddpportal.Group, error)
}

type BadgerGroupRepository struct {
	db *badger.DB
}

func (b BadgerGroupRepository) Insert(group ddpportal.Group) error {
	key := fmt.Sprintf("group_%s", group.ID)
	jsonGroup, err := json.Marshal(group)
	if err != nil {
		return err
	}
	err = b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), jsonGroup)
		return err
	})
	return err
}

func (b BadgerGroupRepository) Update(group ddpportal.Group) error {
	key := fmt.Sprintf("group_%s", group.ID)
	jsonGroup, err := json.Marshal(group)
	if err != nil {
		return err
	}
	err = b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), jsonGroup)
		return err
	})
	return err
}

func (b BadgerGroupRepository) Delete(group ddpportal.Group) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("group_%s", group.ID)
		return txn.Delete([]byte(key))
	})
	return err
}

func (b BadgerGroupRepository) GetAll() ([]*ddpportal.Group, error) {
	groups := make([]*ddpportal.Group, 0)

	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("group_")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var group ddpportal.Group

			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}

			err = json.Unmarshal([]byte(v), &group)
			if err != nil {
				return err
			}

			groups = append(groups, &group)
		}
		return nil
	})

	return groups, err
}

func NewBadgerGroupRepository(db *badger.DB) BadgerGroupRepository {
	return BadgerGroupRepository{db: db}
}
