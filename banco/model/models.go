package model

import (
	"database/sql"
	"time"
)

type Cliente struct {
	ID              int       `db:"id"`
	Nombre          string    `db:"nombre"`
	Apellido        string    `db:"apellido"`
	FechaNacimiento time.Time `db:"fecha_nacimiento"`
	Email           string    `db:"email"`
	Telefono        sql.NullString
	Password        string `db:"password"`
	CodigoCuenta    string `db:"codigo_cuenta"`
}

type Cuenta struct {
	ID            int
	IDCliente     int
	TipoCuenta    string
	Saldo         float64
	FechaCreacion time.Time
}

type Transaccion struct {
	ID               int
	IDCuenta         int
	TipoTransaccion  string
	Monto            float64
	FechaTransaccion time.Time
}

type Administrador struct {
	ID       int    `db:"id"`
	Nombre   string `db:"nombre"`
	Apellido string `db:"apellido"`
	Email    string `db:"email"`
	Password string `db:"password"`
}
