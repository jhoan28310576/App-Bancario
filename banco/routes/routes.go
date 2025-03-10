package routes

import (
	"banco/db" // Asegúrate de importar tu paquete de base de datos
	"banco/model"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Configurar sesiones seguras
	store := cookie.NewStore([]byte("secreto-muy-seguro-32-caracteres"))
	store.Options(sessions.Options{
		MaxAge:   86400 * 7, // 7 días
		HttpOnly: true,
		Secure:   true, // Solo HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	r.Use(sessions.Sessions("banco-session", store))

	// Cargar las plantillas HTML
	r.LoadHTMLGlob("@templates/*") // Asegúrate de que esta ruta sea correcta

	r.GET("/", showLoginPage)
	r.POST("/login", login)
	r.GET("/register", showRegisterPage)
	r.POST("/register", register)

	r.GET("/cuenta", authMiddleware(), showCuentaPage)

	r.GET("/usuario", authMiddleware(), showUsuarioPage) // Añadir esta línea

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

	return r
}

func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title":   "Login",
		"mensaje": "¡Registro exitoso! Por favor inicia sesión",
		"ShowNav": false, // No mostrar la navegación
	})
}

func showRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func login(c *gin.Context) {
	var json struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBind(&json); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Credenciales inválidas"})
		return
	}

	var cliente model.Cliente
	err := db.DB.QueryRow("SELECT id, password FROM Clientes WHERE email = ?", json.Email).Scan(&cliente.ID, &cliente.Password)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(cliente.Password), []byte(json.Password)) != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Credenciales incorrectas"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", cliente.ID)
	session.Save()

	c.Redirect(http.StatusFound, "/cuenta")
}

func register(c *gin.Context) {
	var json struct {
		Nombre          string    `form:"nombre" binding:"required"`
		Apellido        string    `form:"apellido" binding:"required"`
		Email           string    `form:"email" binding:"required,email"`
		Password        string    `form:"password" binding:"required,min=8"`
		FechaNacimiento time.Time `form:"fecha_nacimiento" binding:"required" time_format:"2006-01-02"`
	}

	if err := c.ShouldBindWith(&json, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Convertir el string de fecha a time.Time
	fecha, err := time.Parse("2006-01-02", c.PostForm("fecha_nacimiento"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Use YYYY-MM-DD"})
		return
	}
	json.FechaNacimiento = fecha

	// Insertar el nuevo cliente en la base de datos
	query := "INSERT INTO Clientes (nombre, apellido, email, password, fecha_nacimiento) VALUES (?, ?, ?, ?, ?)"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(json.Password),
		bcrypt.DefaultCost, // Produce hash de 60 caracteres
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el usuario"})
		return
	}

	_, err = db.DB.Exec(query, json.Nombre, json.Apellido, json.Email, hashedPassword, json.FechaNacimiento)
	if err != nil {
		fmt.Printf("ERROR SQL: %v\nQUERY: %s\nPARAMS: %s, %s, %s, [hashed], %s\n",
			err, query, json.Nombre, json.Apellido, json.Email, json.FechaNacimiento)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error en el servidor. Por favor intenta más tarde",
		})
		return
	}
	// Cambiar la respuesta JSON por redirección
	c.Redirect(http.StatusSeeOther, "/login?registro=exitoso")
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func showCuentaPage(c *gin.Context) {
	// Supongamos que obtienes estos datos de la base de datos
	clienteNombre := "Juan Pérez"
	saldo := 1500.00
	transacciones := []struct {
		Fecha       string
		Descripcion string
		Monto       float64
	}{
		{"2023-10-01", "Depósito", 500.00},
		{"2023-10-05", "Pago de Servicios", -100.00},
		{"2023-10-10", "Transferencia", -200.00},
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"Title":         "Cuenta",
		"ShowNav":       true,
		"ClienteNombre": clienteNombre,
		"Saldo":         saldo,
		"Transacciones": transacciones,
	})
}

func showUsuarioPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var cliente struct {
		Nombre   string
		Apellido string
		Email    string
		Telefono sql.NullString // Mantener como sql.NullString para manejar NULL
	}

	err := db.DB.QueryRow("SELECT nombre, apellido, email, telefono FROM Clientes WHERE id = ?", userID).Scan(&cliente.Nombre, &cliente.Apellido, &cliente.Email, &cliente.Telefono)

	if err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v\n", err) // Log del error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los datos del cliente"})
		return
	}

	c.HTML(http.StatusOK, "usuario.html", cliente)
}
