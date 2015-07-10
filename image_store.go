package main

import "database/sql"

const pageSize = 25

type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

var globalImageStore ImageStore

type DBImageStore struct {
	db *sql.DB
}

func NewDBImageStore() ImageStore {
	return &DBImageStore{
		db: globalMySQLDB,
	}
}

func (store *DBImageStore) Save(image *Image) error {
	_, err := store.db.Exec(
		`
		REPLACE INTO images
			(id, user_id,name,location,description,size,created_at)
		VALUES 
			(?,?,?,?,?,?,?)
		`,
		image.ID,
		image.UserID,
		image.Name,
		image.Location,
		image.Description,
		image.Size,
		image.CreateAt,
	)
	return err
}

func (store *DBImageStore) Find(id string) (*Image, error) {
	row := store.db.QueryRow(
		`
		select id,user_id,name,location,description,size,created_at
		from images
		where id = ? 
		`,
		id,
	)
	image := Image{}
	err := row.Scan(
		&image.ID,
		&image.UserID,
		&image.Name,
		&image.Location,
		&image.Description,
		&image.Size,
		&image.CreateAt,
	)
	return &image, err
}

func (store *DBImageStore) FindAll(offset int) ([]Image, error) {
	rows, err := store.db.Query(
		`
		select id,user_id,name,location,description,size,created_at
		from images
		order by created_at desc
		limit ?
		offset ?
		`,
		pageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		image := Image{}
		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.Name,
			&image.Location,
			&image.Description,
			&image.Size,
			&image.CreateAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (store *DBImageStore) FindAllByUser(user *User, offset int) ([]Image, error) {

	rows, err := store.db.Query(
		`
		select id,user_id,name,location,description,size,created_at
		from images
		where user_id = ?
		order by created_at desc
		limit ?
		offset ?
		`,
		user.ID,
		pageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		image := Image{}
		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.Name,
			&image.Location,
			&image.Description,
			&image.Size,
			&image.CreateAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}
