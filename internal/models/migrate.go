package models

// Models devuelve todos los modelos para AutoMigrate.
func Models() []any {
	return []any{
		// Geografía
		&Country{}, &Region{}, &City{}, &Commune{},

		// Dirección / usuarios
		&Address{},
		&User{}, &Contact{}, &UserPhone{},

		// RBAC
		&Role{}, &Permission{},

		// 🏭 BODEGA / WAREHOUSE
		&Space{},
		&SpaceFloor{},
		&Warehouse{},
		&WarehouseRack{},
	}
}
