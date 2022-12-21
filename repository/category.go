package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	res := []entity.Category{}
	resp := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&res).Error
	if resp != nil {
		return []entity.Category{}, resp
	}
	return res, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	store := r.db.WithContext(ctx).Create(&category).Error
	if store != nil {
		return 0, store
	}
	return category.ID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	store := r.db.WithContext(ctx).Create(&categories).Error
	if store != nil {
		return store
	}
	return nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	res := entity.Category{}
	resp := r.db.WithContext(ctx).First(&res, id).Error
	if resp != nil {
		return res, resp
	}
	return res, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	err := r.db.WithContext(ctx).Model(&entity.Category{}).Where("id =?", category.ID).Updates(&category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
