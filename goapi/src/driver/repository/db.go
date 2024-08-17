package repository

import (
	"context"

	"gorm.io/gorm"
)

type DbRepo[T any] struct {
	db *gorm.DB
}

func NewDbRepo[T any](db *gorm.DB) DbRepo[T] {
	return DbRepo[T]{
		db: db,
	}
}

func (d DbRepo[T]) Select(ctx context.Context, id string) ([]T, error) {
	var values []T

	err := d.db.Find(&values).Where(id).Error
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (d DbRepo[T]) Insert(ctx context.Context, value T) error {
	err := d.db.Create(&value).Error
	if err != nil {
		return err
	}

	return nil
}

func (d DbRepo[T]) Update(ctx context.Context, value T) error {
	err := d.db.Save(&value).Error
	if err != nil {
		return err
	}

	return nil
}

func (d DbRepo[T]) Delete(ctx context.Context, value T) error {
	err := d.db.Delete(&value).Error
	if err != nil {
		return err
	}

	return nil
}
