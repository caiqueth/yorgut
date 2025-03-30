package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func (r *Repository[T]) Create(entity T) error {
	result := r.db.Create(&entity)
	if result.Error != nil {
		return fmt.Errorf("error creating entity: %v", result.Error)
	}
	return nil
}

func (r *Repository[T]) Delete(id uint) error {
	var ent T
	result := r.db.Where("id = ?", id).Delete(&ent)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository[T]) FindById(id uint) (T, error) {
	var ent T
	result := r.db.Find(&ent, id)
	if result != nil {
		return ent, nil
	} else {
		return ent, errors.New("entity not found")
	}
}

func (r *Repository[T]) FindMany() ([]T, error) {
	var ents []T
	results := r.db.Find(&ents)
	if results.Error != nil {
		return nil, results.Error
	}
	return ents, nil
}

func (r *Repository[T]) Update(ent T) error {
	updateEntity := ent
	result := r.db.Model(&ent).Updates(updateEntity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
