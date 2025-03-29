package routes

import (
	"banco/db" // Asegúrate de importar tu paquete de base de datos
	"banco/model"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	//"strconv"
	"strconv"
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

	r.GET("/admin", authMiddleware(), showAdminPage)
	r.POST("/registerAdmin", authMiddleware(), registerAdmin)
	r.GET("/registerAdmin", authMiddleware(), showRegisterAdminPage)
	r.POST("/recargar", authMiddleware(), recargarSaldo)
	r.GET("/historial", authMiddleware(), showHistorialRecargas)
	r.GET("/cuenta", authMiddleware(), showCuentaPage)
	r.POST("/abrirCuenta", authMiddleware(), abrirCuenta)
	r.GET("/usuario", authMiddleware(), showUsuarioPage)
	r.GET("/transaccion", authMiddleware(), showTransaccionPage)

	// **Asegúrate de agregar esta línea para la transferencia**
	r.POST("/transferir", authMiddleware(), transferir) // Ruta para manejar la transferencia

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

	// Nuevo endpoint para obtener datos en tiempo real
	r.GET("/obtener-destinatario", func(c *gin.Context) {
		codigo := c.Query("codigo")
		var nombre, apellido string
		err := db.DB.QueryRow(`
			SELECT nombre, apellido 
			FROM Clientes 
			WHERE codigo_cuenta = ?`, codigo).Scan(&nombre, &apellido)

		if err != nil {
			fmt.Println("Error en consulta:", err)
			c.JSON(http.StatusOK, gin.H{"nombre": "No encontrado", "apellido": "No encontrado"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"nombre": nombre, "apellido": apellido})
	})

	return r
}

// funciones del admin

func showRegisterAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "registerAdmin.html", nil)
}

func showAdminPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	// Verificar si el usuario es un administrador
	if userID == nil || userID != 1 { // Asegúrate de que el ID 1 sea el administrador
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener el mensaje de la sesión
	mensaje := session.Get("mensaje")
	session.Delete("mensaje") // Eliminar el mensaje después de mostrarlo

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"mensaje": mensaje, // Pasar el mensaje a la plantilla
		"ShowNav": true,    // Mostrar la barra de navegación
		"IsAdmin": true,    // Indicar que es un administrador
	}) // Cargar la plantilla de administración
}

func registerAdmin(c *gin.Context) {
	var json struct {
		Nombre   string `form:"nombre" binding:"required"`
		Apellido string `form:"apellido" binding:"required"`
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al hashear la contraseña"})
		return
	}

	// Insertar el nuevo administrador en la base de datos
	query := "INSERT INTO Administradores (nombre, apellido, email, password) VALUES (?, ?, ?, ?)"
	_, err = db.DB.Exec(query, json.Nombre, json.Apellido, json.Email, hashedPassword)
	if err != nil {
		fmt.Printf("Error al registrar el administrador: %v\n", err) // Imprimir el error en el servidor
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el administrador"})
		return
	}

	// Establecer la sesión para el nuevo administrador
	session := sessions.Default(c)
	session.Set("user_id", 1) // Asumiendo que el ID del nuevo administrador es 1
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin") // Redirigir a la página de administración
}

func recargarSaldo(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	// Verificar si el usuario es un administrador
	if userID == nil || userID != 1 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	codigoCuenta := c.PostForm("codigoCuenta")
	monto, err := strconv.ParseFloat(c.PostForm("monto"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Monto inválido"})
		return
	}

	// Obtener el ID del cliente usando el codigo_cuenta
	var idCliente int
	err = db.DB.QueryRow("SELECT id FROM Clientes WHERE codigo_cuenta = ?", codigoCuenta).Scan(&idCliente)
	if err != nil {
		fmt.Printf("Error al obtener el ID del cliente: %v\n", err) // Imprimir el error en el servidor
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código de cuenta no encontrado"})
		return
	}

	// Actualizar el saldo en la base de datos
	result, err := db.DB.Exec("UPDATE Cuentas SET saldo = saldo + ? WHERE id_cliente = ?", monto, idCliente)
	if err != nil {
		fmt.Printf("Error al recargar el saldo: %v\n", err) // Imprimir el error en el servidor
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recargar el saldo"})
		return
	}

	// Verificar si se actualizó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error al obtener filas afectadas: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar la actualización"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código de cuenta no encontrado"})
		return
	}

	// Insertar en el historial de recargas
	_, err = db.DB.Exec("INSERT INTO historial_recargas (id_admin, id_cliente, monto, codigo_cuenta) VALUES (?, ?, ?, ?)", userID, idCliente, monto, codigoCuenta)
	if err != nil {
		fmt.Printf("Error al registrar la recarga en el historial: %v\n", err)
	}

	// Establecer el mensaje de éxito en la sesión
	session.Set("mensaje", "Recarga realizada con éxito.")
	session.Save()

	c.Redirect(http.StatusFound, "/admin") // Redirigir a la página de administración
}

func showHistorialRecargas(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	// Verificar si el usuario es un administrador
	if userID == nil || userID != 1 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Modificar la consulta para obtener nombres y apellidos, ordenado por fecha descendente
	rows, err := db.DB.Query(`
        SELECT
            a.nombre as admin_nombre,
            a.apellido as admin_apellido,
            c.nombre as cliente_nombre,
            c.apellido as cliente_apellido,
            h.monto,
            DATE_FORMAT(h.fecha_recarga, '%d/%m/%Y %H:%i') as fecha_recarga,
            h.codigo_cuenta
        FROM historial_recargas h
        JOIN Administradores a ON h.id_admin = a.id
        JOIN Clientes c ON h.id_cliente = c.id
        ORDER BY h.fecha_recarga DESC`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el historial"})
		return
	}
	defer rows.Close()

	var historial []struct {
		AdminNombre     string
		AdminApellido   string
		ClienteNombre   string
		ClienteApellido string
		Monto           float64
		FechaRecarga    string
		CodigoCuenta    string
	}

	for rows.Next() {
		var recarga struct {
			AdminNombre     string
			AdminApellido   string
			ClienteNombre   string
			ClienteApellido string
			Monto           float64
			FechaRecarga    string
			CodigoCuenta    string
		}
		if err := rows.Scan(
			&recarga.AdminNombre,
			&recarga.AdminApellido,
			&recarga.ClienteNombre,
			&recarga.ClienteApellido,
			&recarga.Monto,
			&recarga.FechaRecarga,
			&recarga.CodigoCuenta); err != nil {
			fmt.Printf("Error al escanear el historial: %v\n", err)
			continue
		}
		historial = append(historial, recarga)
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"historial": historial,
		"ShowNav":   true,
		"IsAdmin":   true,
		"Title":     "Historial de Recargas",
		"Content":   "historial.html",
	})
}

// fin funciones del admin

func abrirCuenta(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Insertar una nueva cuenta para el cliente
	query := "INSERT INTO Cuentas (id_cliente, tipo_cuenta, saldo) VALUES (?, ?, ?)"
	_, err := db.DB.Exec(query, userID, "Cuenta de Ahorros", 0) // Puedes cambiar el tipo de cuenta según sea necesario
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al abrir la cuenta"})
		return
	}

	c.Redirect(http.StatusFound, "/cuenta") // Redirigir a la página de cuenta
}

// funciones de login
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

	// Verificar si el usuario es un cliente
	var cliente model.Cliente
	err := db.DB.QueryRow("SELECT id, password FROM Clientes WHERE email = ?", json.Email).Scan(&cliente.ID, &cliente.Password)

	if err != nil {
		// Si no se encuentra un cliente, verificar si es un administrador
		var admin model.Administrador
		err = db.DB.QueryRow("SELECT id, password FROM Administradores WHERE email = ?", json.Email).Scan(&admin.ID, &admin.Password)

		if err != nil || bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(json.Password)) != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Credenciales incorrectas"})
			return
		}

		// Si es un administrador, establecer la sesión
		session := sessions.Default(c)
		session.Set("user_id", admin.ID) // Almacenar el ID del administrador en la sesión
		session.Save()

		c.Redirect(http.StatusFound, "/admin") // Redirigir a la página de administración
		return
	}

	// Si es un cliente, verificar la contraseña
	if bcrypt.CompareHashAndPassword([]byte(cliente.Password), []byte(json.Password)) != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Si es un cliente, establecer la sesión
	session := sessions.Default(c)
	session.Set("user_id", cliente.ID) // Almacenar el ID del cliente en la sesión
	session.Save()

	c.Redirect(http.StatusFound, "/cuenta") // Redirigir a la página de cuenta
}

// Función para generar un código de cuenta único
func generarCodigoCuenta() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("CUENTA-%d", rand.Intn(1000000)) // Genera un código aleatorio
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

	// Generar el código de cuenta
	codigoCuenta := generarCodigoCuenta()

	// Insertar el nuevo cliente en la base de datos
	query := "INSERT INTO Clientes (nombre, apellido, email, password, fecha_nacimiento, codigo_cuenta) VALUES (?, ?, ?, ?, ?, ?)"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(json.Password),
		bcrypt.DefaultCost, // Produce hash de 60 caracteres
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el usuario"})
		return
	}

	_, err = db.DB.Exec(query, json.Nombre, json.Apellido, json.Email, hashedPassword, json.FechaNacimiento, codigoCuenta)
	if err != nil {
		fmt.Printf("ERROR SQL: %v\nQUERY: %s\nPARAMS: %s, %s, %s, [hashed], %s, %s\n",
			err, query, json.Nombre, json.Apellido, json.Email, json.FechaNacimiento, codigoCuenta)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error en el servidor. Por favor intenta más tarde",
		})
		return
	}
	// Cambiar la respuesta JSON por redirección
	c.Redirect(http.StatusSeeOther, "/")
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

// fin funciones de login

// funciones de pagina cuenta

func showCuentaPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Estructura para almacenar los datos del cliente
	var cliente struct {
		Nombre   string
		Apellido string
		Saldo    float64
	}

	// Verificar si el usuario tiene una cuenta
	var tieneCuenta bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Cuentas WHERE id_cliente = ?)", userID).Scan(&tieneCuenta)
	if err != nil {
		fmt.Printf("Error al verificar cuenta: %v\n", err)
	}

	// Obtener nombre, apellido y saldo del cliente
	err = db.DB.QueryRow(`
		SELECT 
			c.nombre, 
			c.apellido, 
			COALESCE(cu.saldo, 0) as saldo
		FROM Clientes c
		LEFT JOIN Cuentas cu ON c.id = cu.id_cliente
		WHERE c.id = ?`, userID).Scan(&cliente.Nombre, &cliente.Apellido, &cliente.Saldo)

	if err != nil {
		fmt.Printf("Error al obtener datos del cliente: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos del cliente"})
		return
	}

	// Formatear el saldo con el símbolo de Bolívares
	saldoFormateado := fmt.Sprintf("%.2f Bs", cliente.Saldo)

	// Obtener el historial de transacciones
	transacciones, err := historialtransaccioneusuarios(userID)
	if err != nil {
		fmt.Printf("Error al obtener historial: %v\n", err)
	}

	c.HTML(http.StatusOK, "cuenta.html", gin.H{
		"Title":         "Cuenta",
		"ShowNav":       true,
		"IsUser":        true,
		"Nombre":        cliente.Nombre,
		"Apellido":      cliente.Apellido,
		"Saldo":         saldoFormateado,
		"Content":       "cuenta.html",
		"Transacciones": transacciones,
		"TieneCuenta":   tieneCuenta, // Pasar el estado de la cuenta a la plantilla
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

	// Obtener los datos del cliente desde la base de datos
	err := db.DB.QueryRow("SELECT nombre, apellido, email, telefono FROM Clientes WHERE id = ?", userID).Scan(&cliente.Nombre, &cliente.Apellido, &cliente.Email, &cliente.Telefono)

	if err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v\n", err) // Log del error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los datos del cliente"})
		return
	}

	// Pasar los datos a la plantilla
	c.HTML(http.StatusOK, "usuario.html", gin.H{
		"ShowNav":  true,
		"IsUser":   true, // Indicar que es un usuario
		"Nombre":   cliente.Nombre,
		"Apellido": cliente.Apellido,
		"Email":    cliente.Email,
		"Telefono": cliente.Telefono, // Pasar el campo Telefono
	})
}

func showTransaccionPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var saldo float64
	// Obtener el saldo del cliente desde la base de datos
	err := db.DB.QueryRow("SELECT saldo FROM cuentas WHERE id_cliente = ?", userID).Scan(&saldo)
	if err != nil {
		fmt.Printf("Error al obtener el saldo: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el saldo"})
		return
	}

	// Formatear el saldo a un formato legible en bolívares
	saldoFormateado := fmt.Sprintf("%.2f Bs", saldo) // Formato a dos decimales

	// Obtener el mensaje de la sesión y limpiarlo inmediatamente
	mensaje := session.Get("mensaje")
	session.Delete("mensaje")
	session.Save()

	c.HTML(http.StatusOK, "transaccion.html", gin.H{
		"Saldo":   saldoFormateado, // Pasar el saldo formateado a la plantilla
		"IsUser":  true,            // Indicar que es un usuario
		"mensaje": mensaje,         // Pasar el mensaje a la plantilla
	})
}

func transferir(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener el código de cuenta y el monto a transferir
	codigoCuentaDestino := c.PostForm("codigoCuenta")
	monto, err := strconv.ParseFloat(c.PostForm("monto"), 64)
	if err != nil {
		session.Set("mensaje", "Monto inválido")
		session.Save()
		c.Redirect(http.StatusFound, "/transaccion")
		return
	}

	// Verificar que la cuenta destino existe y obtener los datos del destinatario
	var nombreDestino, apellidoDestino string
	var idClienteDestino int
	err = db.DB.QueryRow("SELECT id, nombre, apellido FROM Clientes WHERE codigo_cuenta = ?", codigoCuentaDestino).Scan(&idClienteDestino, &nombreDestino, &apellidoDestino)
	if err != nil {
		session.Set("mensaje", "Código de cuenta destino no encontrado")
		session.Save()
		c.Redirect(http.StatusFound, "/transaccion")
		return
	}

	// Obtener el saldo del usuario que está realizando la transferencia
	var saldo float64
	err = db.DB.QueryRow("SELECT saldo FROM Cuentas WHERE id_cliente = ?", userID).Scan(&saldo)
	if err != nil {
		session.Set("mensaje", "Error al obtener el saldo")
		session.Save()
		c.Redirect(http.StatusFound, "/transaccion")
		return
	}

	// Verificar que el saldo es suficiente
	if saldo < monto {
		session.Set("mensaje", "Saldo insuficiente para realizar la transferencia")
		session.Save()
		c.Redirect(http.StatusFound, "/transaccion")
		return
	}

	// Realizar la transferencia
	// Restar del saldo del usuario
	_, err = db.DB.Exec("UPDATE Cuentas SET saldo = saldo - ? WHERE id_cliente = ?", monto, userID)
	if err != nil {
		session.Set("mensaje", "Error al realizar la transferencia")
		session.Save()
		c.Redirect(http.StatusFound, "/transaccion")
		return
	}

	// Sumar al saldo del destinatario usando el id_cliente
	_, err = db.DB.Exec("UPDATE Cuentas SET saldo = saldo + ? WHERE id_cliente = ?", monto, idClienteDestino)
	if err != nil {
		session.Set("mensaje", "Error al actualizar el saldo del destinatario")
		session.Save()
		c.Redirect(http.StatusFound, "/transaccion")
		return
	}

	// Después de realizar la transferencia exitosamente, registrar en el historial
	descripcion := fmt.Sprintf("Transferencia a %s %s", nombreDestino, apellidoDestino)

	// Obtener el código de cuenta del emisor
	var codigoCuentaEmisor string
	err = db.DB.QueryRow("SELECT codigo_cuenta FROM Clientes WHERE id = ?", userID).Scan(&codigoCuentaEmisor)
	if err != nil {
		fmt.Printf("Error al obtener código de cuenta emisor: %v\n", err)
	}

	// Registrar la transacción en el historial
	_, err = db.DB.Exec(`
		INSERT INTO historial_transacciones_usuarios 
		(id_emisor, id_receptor, monto, descripcion, codigo_cuenta_emisor, codigo_cuenta_receptor) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		userID, idClienteDestino, monto, descripcion, codigoCuentaEmisor, codigoCuentaDestino)

	if err != nil {
		fmt.Printf("Error al registrar en historial: %v\n", err)
	}

	// Mostrar un mensaje de éxito
	session.Set("mensaje", fmt.Sprintf("Transferencia de %.2f Bs a %s %s (Código: %s) realizada con éxito.", monto, nombreDestino, apellidoDestino, codigoCuentaDestino))
	session.Save()

	c.Redirect(http.StatusFound, "/transaccion")
}

func historialtransaccioneusuarios(userID interface{}) ([]struct {
	Fecha         string
	Descripcion   string
	Monto         float64
	NombreOtro    string
	ApellidoOtro  string
	CodigoCuenta  string
	TipoOperacion string
}, error) {
	var transacciones []struct {
		Fecha         string
		Descripcion   string
		Monto         float64
		NombreOtro    string
		ApellidoOtro  string
		CodigoCuenta  string
		TipoOperacion string
	}

	// Consulta para obtener tanto las transacciones enviadas como recibidas
	rows, err := db.DB.Query(`
		SELECT 
			DATE_FORMAT(h.fecha_transaccion, '%d/%m/%Y %H:%i') as fecha,
			h.descripcion,
			h.monto,
			CASE 
				WHEN h.id_emisor = ? THEN c_receptor.nombre
				ELSE c_emisor.nombre
			END as nombre_otro,
			CASE 
				WHEN h.id_emisor = ? THEN c_receptor.apellido
				ELSE c_emisor.apellido
			END as apellido_otro,
			CASE 
				WHEN h.id_emisor = ? THEN h.codigo_cuenta_receptor
				ELSE h.codigo_cuenta_emisor
			END as codigo_cuenta,
			CASE 
				WHEN h.id_emisor = ? THEN 'Enviado'
				ELSE 'Recibido'
			END as tipo_operacion
		FROM historial_transacciones_usuarios h
		JOIN Clientes c_emisor ON h.id_emisor = c_emisor.id
		JOIN Clientes c_receptor ON h.id_receptor = c_receptor.id
		WHERE h.id_emisor = ? OR h.id_receptor = ?
		ORDER BY h.fecha_transaccion DESC`,
		userID, userID, userID, userID, userID, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t struct {
			Fecha         string
			Descripcion   string
			Monto         float64
			NombreOtro    string
			ApellidoOtro  string
			CodigoCuenta  string
			TipoOperacion string
		}
		err := rows.Scan(&t.Fecha, &t.Descripcion, &t.Monto, &t.NombreOtro,
			&t.ApellidoOtro, &t.CodigoCuenta, &t.TipoOperacion)
		if err != nil {
			return nil, err
		}
		transacciones = append(transacciones, t)
	}

	return transacciones, nil
}
