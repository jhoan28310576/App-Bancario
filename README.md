App Bancaria + Simulador Económico
(Inspirada en "La Riqueza de las Naciones" de Adam Smith)

https://img.shields.io/badge/Go-1.18%252B-blue.svg
https://img.shields.io/badge/License-MIT-yellow.svg
https://img.shields.io/badge/Gin_Framework-1.8%252B-lightblue

Un sistema bancario completo con simulador económico inspirado en los principios de Adam Smith. Desarrollado en Go con Gin para el enrutamiento y MySQL como base de datos, implementa operaciones bancarias seguras y un modelo económico virtual.

https://github.com/user-attachments/assets/f4783aa3-de0e-4062-8070-190a59c188ac

🌟 Características Principales
🔐 Sistema de Autenticación Avanzado
Registro de usuarios con validación de campos

Inicio de sesión seguro con sesiones persistentes

Almacenamiento de contraseñas con Bcrypt

Gestión de sesiones con cookies cifradas

💰 Operaciones Bancarias Esenciales
Depósitos: Acreditación de fondos con verificación

Consultas de saldo en tiempo real

Historial de transacciones detallado

Transferencias entre cuentas con validación

⚙️ Módulo de Depósitos
Interfaz administrativa para gestión de depósitos

Validación de efectivo físico

Actualización instantánea de saldos

🔄 Sistema de Transferencias Seguras
Verificación de saldo antes de operaciones

Generación de códigos únicos por transacción

Validación en tiempo real de códigos de seguridad

Historial completo de operaciones auditables

🛠️ Tecnologías Utilizadas
Backend
Tecnología	Uso
Go (Golang) 1.18+	Lenguaje principal del sistema
Gin Framework	Enrutamiento y manejo de middlewares
Gin Sessions	Gestión de sesiones de usuario
MySQL	Almacenamiento persistente de datos
go-sql-driver/mysql	Conexión a base de datos MySQL
Bcrypt	Cifrado seguro de contraseñas
Gin Binding	Validación de formularios y datos
Frontend
Tecnología	Uso
Gin HTML Templates	Renderizado de vistas del servidor
CSS/HTML	Estructura y estilos básicos
(Extensible a Bootstrap)	Para futuras mejoras de UI
🚀 Instalación y Configuración
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
  

## 🏗️ Arquitectura 
![Arquitectura](https://github.com/user-attachments/assets/f4783aa3-de0e-4062-8070-190a59c188ac)
