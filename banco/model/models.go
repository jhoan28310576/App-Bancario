package model

import (
	"time"
)

type Cliente struct {
	ID              int       `db:"id"`
	Nombre          string    `db:"nombre"`
	Apellido        string    `db:"apellido"`
	FechaNacimiento time.Time `db:"fecha_nacimiento"`
	Email           string    `db:"email"`
	Telefono        string
	Password        string `db:"password"`
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
