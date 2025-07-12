# App Bancaria + Simulador Económico 
*(Inspirada en "La Riqueza de las Naciones" de Adam Smith)*  

## 🏗️ Arquitectura 
![Arquitectura](https://github.com/user-attachments/assets/f4783aa3-de0e-4062-8070-190a59c188ac)

Un sistema bancario completo desarrollado en Go utilizando el framework Gin para el manejo de rutas y MySQL como base de datos. Implementa operaciones bancarias esenciales con seguridad mejorada.

Características Principales:

Sistema de Autenticación:
Registro de nuevos usuarios con validación de campos

Inicio de sesión seguro con sesiones persistentes

Almacenamiento seguro de contraseñas usando bcrypt

Gestión de sesiones con cookies cifradas

Operaciones Bancarias:

Depósitos: Acreditación de fondos a cuentas

Consultas de Saldo: Visualización de saldo actual

Historial de Transacciones: Registro detallado de movimientos

Transferencias: Envío de fondos entre cuentas con validación

Módulo de Depósitos:

Interfaz administrativa para gestionar depósitos

Validación de efectivo físico

Actualización instantánea de saldos

Sistema de Transferencias:

Transferencias entre cuentas con verificación de saldo

Generación de códigos únicos por transacción

Validación en tiempo real de códigos de seguridad

Historial completo de operaciones

Tecnologías Utilizadas:

Backend

Lenguaje Principal: Go (Golang) 1.18+

Framework Web: Gin (para enrutamiento y middleware)

Gestión de Sesiones: Gin Sessions con almacenamiento en cookies

Base de Datos: MySQL

Driver de DB: github.com/go-sql-driver/mysql

Hashing de Contraseñas: Bcrypt (golang.org/x/crypto/bcrypt)

Manejo de Formularios: Gin Binding

Frontend:

Plantillas HTML: Gin HTML rendering

Estilos: CSS (puede extenderse con Bootstrap u otros frameworks)

Funcionalidades Técnicas:

Generación de códigos únicos para transferencias

Validación de datos de entrada en formularios

Manejo de errores y respuestas HTTP adecuadas

Conexión segura a base de datos MySQL

Instalación y Configuración:

Requisitos Previos

Go 1.18+

MySQL 5.7+

Git



Clonar el repositorio:

git clone https://github.com/jhoan28310576/App-Bancario.git

cd App-Bancario/banco

Instalar dependencias:

go mod download

Configurar conexión a DB (en db/db.go):

const (
    host     = "localhost"
    port     = 5432
    user     = "tu_usuario"
    password = "tu_contraseña"
    dbname   = "app_bancario"
)

Iniciar el servidor:

go run main.go

