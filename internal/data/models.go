package data

import (
	"database/sql"
)

type Models struct {
	Movies     MovieModel
	Users      UserModel
	Token      TokenModel
	Permissons PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:     MovieModel{DB: db},
		Users:      UserModel{DB: db},
		Token:      TokenModel{DB: db},
		Permissons: PermissionModel{DB: db},
	}
}
