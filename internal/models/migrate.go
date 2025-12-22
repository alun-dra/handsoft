package models

// Models devuelve todos los modelos para AutoMigrate.
func Models() []any {
	return []any{
		&Country{}, &Region{}, &City{}, &Commune{},
		&Address{},
		&User{}, &Contact{}, &UserPhone{},
		&Role{}, &Permission{},
	}
}
