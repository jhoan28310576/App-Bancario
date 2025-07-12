# 🏦 App Bancaria + Simulador Económico
*(Inspirada en "La Riqueza de las Naciones" de Adam Smith)*

Un sistema bancario completo desarrollado en Go utilizando el framework Gin para el manejo de rutas y MySQL como base de datos. Implementa operaciones bancarias esenciales con seguridad mejorada y principios económicos fundamentales.

## 📋 Tabla de Contenidos
- [Características Principales](#características-principales)
- [Tecnologías Utilizadas](#tecnologías-utilizadas)
- [Arquitectura del Sistema](#arquitectura-del-sistema)
- [Instalación y Configuración](#instalación-y-configuración)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Uso del Sistema](#uso-del-sistema)
- [API Endpoints](#api-endpoints)
- [Contribución](#contribución)
- [Licencia](#licencia)

## ✨ Características Principales

### 🔐 Sistema de Autenticación
- **Registro de usuarios** con validación completa de campos
- **Inicio de sesión seguro** con sesiones persistentes
- **Almacenamiento seguro** de contraseñas usando bcrypt
- **Gestión de sesiones** con cookies cifradas
- **Panel administrativo** para gestión de usuarios

### 💰 Operaciones Bancarias
- **Depósitos**: Acreditación de fondos a cuentas con validación
- **Consultas de Saldo**: Visualización de saldo actual en tiempo real
- **Historial de Transacciones**: Registro detallado de movimientos
- **Transferencias**: Envío de fondos entre cuentas con validación de saldo
- **Códigos de Seguridad**: Generación de códigos únicos por transacción

### 🏛️ Módulo de Depósitos
- **Interfaz administrativa** para gestionar depósitos
- **Validación de efectivo físico** con verificación
- **Actualización instantánea** de saldos
- **Registro de transacciones** con timestamps

### 🔄 Sistema de Transferencias
- **Transferencias entre cuentas** con verificación de saldo
- **Generación de códigos únicos** por transacción
- **Validación en tiempo real** de códigos de seguridad
- **Historial completo** de operaciones
- **Notificaciones** de transacciones exitosas

## 🛠️ Tecnologías Utilizadas

### Backend
- **Lenguaje Principal**: Go (Golang) 1.23.3+
- **Framework Web**: Gin (para enrutamiento y middleware)
- **Gestión de Sesiones**: Gin Sessions con almacenamiento en cookies
- **Base de Datos**: MySQL 8.0+
- **Driver de DB**: `github.com/go-sql-driver/mysql`
- **Hashing de Contraseñas**: Bcrypt (`golang.org/x/crypto/bcrypt`)
- **Manejo de Formularios**: Gin Binding

### Frontend
- **Plantillas HTML**: Gin HTML rendering
- **Estilos**: CSS personalizado
- **Interfaz Responsiva**: Diseño adaptable a diferentes dispositivos

### Funcionalidades Técnicas
- **Generación de códigos únicos** para transferencias
- **Validación de datos** de entrada en formularios
- **Manejo de errores** y respuestas HTTP adecuadas
- **Conexión segura** a base de datos MySQL
- **Sesiones persistentes** con cookies cifradas

## 🏗️ Arquitectura del Sistema

```
App-Bancario/
├── banco/                    # Directorio principal de la aplicación
│   ├── main.go              # Punto de entrada de la aplicación
│   ├── go.mod               # Dependencias de Go
│   ├── go.sum               # Checksums de dependencias
│   ├── db/                  # Configuración de base de datos
│   │   └── db.go           # Conexión y configuración MySQL
│   ├── model/               # Modelos de datos
│   │   └── models.go       # Estructuras de Cliente, Cuenta, Transaccion
│   ├── routes/              # Controladores y rutas
│   │   └── routes.go       # Definición de endpoints y lógica de negocio
│   ├── templates/           # Plantillas HTML
│   │   ├── base.html       # Plantilla base
│   │   ├── login.html      # Página de inicio de sesión
│   │   ├── register.html   # Página de registro
│   │   ├── cuenta.html     # Dashboard de cuenta
│   │   ├── transaccion.html # Interfaz de transferencias
│   │   ├── historial.html  # Historial de transacciones
│   │   ├── admin.html      # Panel administrativo
│   │   └── nav.html        # Navegación
│   ├── assets/             # Recursos estáticos
│   │   ├── css/           # Estilos CSS
│   │   ├── img/           # Imágenes
│   │   ├── fonts/         # Fuentes
│   │   └── icon/          # Iconos
│   ├── test/              # Pruebas unitarias
│   └── doc/               # Documentación adicional
├── .gitignore             # Archivos ignorados por Git
├── .gitattributes         # Configuración de Git
└── README.md              # Este archivo
```

## 🚀 Instalación y Configuración

### Requisitos Previos
- **Go** 1.23.3 o superior
- **MySQL** 8.0 o superior
- **Git**

### Pasos de Instalación

1. **Clonar el repositorio**:
   ```bash
   git clone https://github.com/jhoan28310576/App-Bancario.git
   cd App-Bancario/banco
   ```

2. **Instalar dependencias**:
   ```bash
   go mod download
   ```

3. **Configurar la base de datos MySQL**:
OPCION 1)  importar a la db de MYSQL contiene todo los usuarios y datos ya establecidos
<img width="182" height="153" alt="image" src="https://github.com/user-attachments/assets/dece8e0c-36a7-4d81-82de-b6c14408f01f" />

en notas.txt pueden encontrar las claves de los usuarios y admin del sistema 
<img width="410" height="340" alt="image" src="https://github.com/user-attachments/assets/ca9e4f8b-8b33-4b1e-8554-d943804e810b" />

OPCION 2) crear la db y sus tablas con sus campos y datos 

   ```sql
   CREATE DATABASE Banco;
   USE Banco;
   
   -- Tabla de clientes
   CREATE TABLE clientes (
       id INT AUTO_INCREMENT PRIMARY KEY,
       nombre VARCHAR(100) NOT NULL,
       apellido VARCHAR(100) NOT NULL,
       fecha_nacimiento DATE NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       telefono VARCHAR(20),
       password VARCHAR(255) NOT NULL,
       codigo_cuenta VARCHAR(20) UNIQUE NOT NULL
   );
   
   -- Tabla de cuentas
   CREATE TABLE cuentas (
       id INT AUTO_INCREMENT PRIMARY KEY,
       id_cliente INT NOT NULL,
       tipo_cuenta VARCHAR(50) NOT NULL,
       saldo DECIMAL(15,2) DEFAULT 0.00,
       fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY (id_cliente) REFERENCES clientes(id)
   );
   
   -- Tabla de transacciones
   CREATE TABLE transacciones (
       id INT AUTO_INCREMENT PRIMARY KEY,
       id_cuenta INT NOT NULL,
       tipo_transaccion VARCHAR(50) NOT NULL,
       monto DECIMAL(15,2) NOT NULL,
       fecha_transaccion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY (id_cuenta) REFERENCES cuentas(id)
   );
   
   -- Tabla de administradores
   CREATE TABLE administradores (
       id INT AUTO_INCREMENT PRIMARY KEY,
       nombre VARCHAR(100) NOT NULL,
       apellido VARCHAR(100) NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       password VARCHAR(255) NOT NULL
   );
   ```

5. **Configurar conexión a la base de datos** (en `db/db.go`):
   ```go
   // Modificar la línea de conexión según tus credenciales
   DB, err = sql.Open("mysql", "usuario:contraseña@tcp(127.0.0.1:3306)/Banco")
   ```

6. **Iniciar el servidor**:
   ```bash
   go run main.go
   ```

7. **Acceder a la aplicación**:
   - Abrir navegador en: `http://localhost:8080`

## 📁 Estructura del Proyecto

### Modelos de Datos (`model/models.go`)
- **Cliente**: Información personal del usuario
- **Cuenta**: Datos de la cuenta bancaria
- **Transaccion**: Registro de movimientos financieros
- **Administrador**: Gestión de usuarios del sistema

### Rutas y Controladores (`routes/routes.go`)
- **Autenticación**: Login, registro, logout
- **Operaciones bancarias**: Depósitos, transferencias, consultas
- **Panel administrativo**: Gestión de usuarios y transacciones
- **Validaciones**: Verificación de datos y permisos

### Plantillas HTML (`templates/`)
- **Interfaz de usuario**: Formularios y dashboards
- **Navegación**: Menús y estructura de páginas
- **Responsive design**: Adaptable a diferentes dispositivos

## 💻 Uso del Sistema

### Para Usuarios Regulares
1. **Registro**: Crear una nueva cuenta con datos personales
2. **Login**: Iniciar sesión con email y contraseña
3. **Dashboard**: Ver saldo actual y opciones disponibles
4. **Transferencias**: Enviar dinero a otras cuentas
5. **Historial**: Consultar movimientos anteriores

### Para Administradores
1. **Panel Admin**: Acceso a gestión de usuarios
2. **Depósitos**: Procesar depósitos físicos
3. **Monitoreo**: Supervisar transacciones del sistema
4. **Usuarios**: Gestionar cuentas de clientes

## 🔌 API Endpoints

### Autenticación
- `GET /` - Página principal
- `GET /login` - Formulario de login
- `POST /login` - Procesar login
- `GET /register` - Formulario de registro
- `POST /register` - Procesar registro
- `GET /logout` - Cerrar sesión

### Operaciones Bancarias
- `GET /cuenta` - Dashboard de cuenta
- `POST /deposito` - Realizar depósito
- `GET /transaccion` - Formulario de transferencia
- `POST /transaccion` - Procesar transferencia
- `GET /historial` - Historial de transacciones

### Administración
- `GET /admin` - Panel administrativo
- `GET /admin/register` - Registro de administrador
- `POST /admin/register` - Crear administrador

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🙏 Agradecimientos

- Inspirado en los principios económicos de "La Riqueza de las Naciones" de Adam Smith
- Desarrollado con tecnologías modernas y mejores prácticas de seguridad
- Diseñado para ser escalable y mantenible

---

**Desarrollado con ❤️ usando Go y Gin Framework**

