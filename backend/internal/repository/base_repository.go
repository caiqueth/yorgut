package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	Db *gorm.DB
}

func (r *Repository[T]) Create(entity T) error {
	result := r.Db.Create(&entity)
	if result.Error != nil {
		return fmt.Errorf("error creating entity: %v", "asd")
	}
	return nil
}

func (r *Repository[T]) Delete(id int) {
	var ent T
	result := r.Db.Where("id = ?", id).Delete(&ent)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (r *Repository[T]) FindById(id int) (T, error) {
	var ent T
	result := r.Db.Find(&ent, id)
	if result != nil {
		return ent, nil
	} else {
		return ent, errors.New("entity not found")
	}
}

func (r *Repository[T]) FindMany() []T {
	var ents []T
	results := r.Db.Find(&ents)
	if results.Error != nil {
		panic(results.Error)
	}
	return ents
}

func (r *Repository[T]) Update(ent T) error {
	updateEntity := ent
	result := r.Db.Model(&ent).Updates(updateEntity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
