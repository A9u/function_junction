package db

import (
	"context"
	"fmt"

	// "database/sql"
	"github.com/mongodb/mongo-go-driver/mongo"
	// "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/bson"
	// "time"
)

const (
	createCategoryQuery = `INSERT INTO categories (
        name, created_at, updated_at)
        VALUES($1, $2, $3)
        `
	listCategoriesQuery     = `SELECT * FROM categories`
	findCategoryByIDQuery   = `SELECT * FROM categories WHERE id = $1`
	deleteCategoryByIDQuery = `DELETE FROM categories WHERE id = $1`
	updateCategoryQuery     = `UPDATE categories SET name = $1, updated_at = $2 where id = $3`
)

type Category struct {
	// ID        string    `db:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	// CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time `db:"updated_at"`
}

func (s *store) CreateCategory(ctx context.Context, collection *mongo.Collection, category *Category) (err error) {
	_, err = collection.InsertOne(ctx, category)
	// id := res.InsertedID
	// now := time.Now()

	// return Transact(ctx, s.db, &sql.TxOptions{}, func(ctx context.Context) error {
	// 	_, err = s.db.Exec(
	// 		createCategoryQuery,
	// 		category.Name,
	// 		now,
	// 		now,
	// 	)
	return err
	// })
}

func (s *store) ListCategories(ctx context.Context, collection *mongo.Collection) (categories []*Category, err error) {
	// findOptions := options.Find()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Error in find: ", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem Category
		err = cur.Decode(&elem)
		categories = append(categories, &elem)
	}
	if err := cur.Err(); err != nil {
	}
	return categories, err
}

func (s *store) FindCategoryByID(ctx context.Context, id string) (category Category, err error) {
	// err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
	// 	return s.db.GetContext(ctx, &category, findCategoryByIDQuery, id)
	// })
	// if err == sql.ErrNoRows {
	// 	return category, ErrCategoryNotExist
	// }
	return
}

func (s *store) DeleteCategoryByID(ctx context.Context, id string) (err error) {
	// return Transact(ctx, s.db, &sql.TxOptions{}, func(ctx context.Context) error {
	// 	res, err := s.db.Exec(deleteCategoryByIDQuery, id)
	// 	cnt, err := res.RowsAffected()
	// 	if cnt == 0 {
	// 		return ErrCategoryNotExist
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}
	return err
	// })
}

func (s *store) UpdateCategory(ctx context.Context, collection *mongo.Collection, filter *Category, category *Category) (err error) {
	// fmt.Println(filter, category)
	// result := Category{}
	// filter.Name = category.Name
	// result,err = collection.FindOne(ctx, filter)

	_, err = collection.UpdateOne(ctx, bson.D{{"name", filter.Name}}, bson.D{{"$set", bson.D{{"name", category.Name}, {"type", category.Type}}}})

	// now := time.Now()

	// return Transact(ctx, s.db, &sql.TxOptions{}, func(ctx context.Context) error {
	// 	_, err = s.db.Exec(
	// 		updateCategoryQuery,
	// 		category.Name,
	// 		now,
	// 		category.ID,
	// 	)
	return err
	// })
}
