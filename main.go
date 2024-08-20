package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Cargamos el driver de mysql (Ejecutar en consola: go get -u github.com/go-sql-driver/mysql)
	"log"                              // Información en terminal
	"net/http"                         // Para el sitio
	"text/template"                    // Cargar información en plantillas
)

// Obtenemos información de carpetas
var plantillas = template.Must(template.ParseGlob("plantillas/*"))

// Función para conectarnos a la base de datos
func conecionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "Sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}

	return conexion
}

type Orden struct {
	Lote         int
	Tamanio_lote int
	Producto     string
	Codigo       string
}

func main() {
	http.HandleFunc("/", inicio)     // Al ejecutar, entramos a la función 'inicio'
	http.HandleFunc("/crear", crear) // Al ejecutar, entramos a la función 'Crear'
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/borrar", borrar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/retirar", retirar)
	http.HandleFunc("/actualizar", actualizar)
	http.HandleFunc("/actualizarRetiro", actualizarRetiro)
	log.Printf("Servidor corriendo")
	http.ListenAndServe(":8080", nil)
}

func inicio(w http.ResponseWriter, r *http.Request) {

	// Conectando base de datos
	conexionEstablecida := conecionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM Ordenes")
	if err != nil {
		panic(err.Error())
	}

	orden := Orden{}
	conjuntoOrden := []Orden{}

	// Ejecutamos la sentencia SQL, recorriendo cada registro
	for registros.Next() {
		var lote, tamanio_lote int
		var producto, codigo string

		err = registros.Scan(&lote, &tamanio_lote, &producto, &codigo) // Obtenemos los datos
		if err != nil {
			panic(err.Error())
		}

		orden.Lote = lote
		orden.Tamanio_lote = tamanio_lote
		orden.Producto = producto
		orden.Codigo = codigo

		conjuntoOrden = append(conjuntoOrden, orden)
	}
	/*
		for _, empleado := range conjuntoEmpleado {
			fmt.Println(empleado)
		} //fmt.Println(conjuntoEmpleado)
	*/
	//fmt.Fprintf(w, "Hola mundo")

	plantillas.ExecuteTemplate(w, "inicio", conjuntoOrden)
}

func crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hola mundo")
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func insertar(w http.ResponseWriter, r *http.Request) {
	// Si hay datos "POST", los redireccionamos
	if r.Method == "POST" {
		lote := r.FormValue("lote")
		tamanio_lote := r.FormValue("tamanio_lote")
		producto := r.FormValue("producto")
		codigo := r.FormValue("codigo")

		// Conectando base de datos
		conexionEstablecida := conecionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO Ordenes(Lote, Tamanio_lote, Producto, Codigo) VALUES(?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(lote, tamanio_lote, producto, codigo)

		// redireccionamos a una pagina
		http.Redirect(w, r, "/", 301)
	}
}

// Recepción de datos para eliminarlos
func borrar(w http.ResponseWriter, r *http.Request) {

	// Recibimos el dato ID
	loteOrden := r.URL.Query().Get("lote")
	fmt.Println(loteOrden)

	// Conectando base de datos
	conexionEstablecida := conecionBD()
	eliminarRegistro, err := conexionEstablecida.Prepare("DELETE FROM Ordenes WHERE Lote=?")
	if err != nil {
		panic(err.Error())
	}

	eliminarRegistro.Exec(loteOrden)

	// redireccionamos a una pagina
	http.Redirect(w, r, "/", 301)
}

// Editar datos del usuario
func editar(w http.ResponseWriter, r *http.Request) {

	// Recibimos el dato ID
	loteOrden := r.URL.Query().Get("lote")

	// Conectando base de datos
	conexionEstablecida := conecionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM Ordenes WHERE Lote=?", loteOrden)
	if err != nil {
		panic(err.Error())
	}

	// Ejecutamos la sentencia SQL, recorriendo cada registro
	orden := Orden{}
	for registro.Next() {
		var lote, tamanio_lote int
		var producto, codigo string

		err = registro.Scan(&lote, &tamanio_lote, &producto, &codigo) // Obtenemos los datos
		if err != nil {
			panic(err.Error())
		}

		orden.Lote = lote
		orden.Tamanio_lote = tamanio_lote
		orden.Producto = producto
		orden.Codigo = codigo
	}

	// Ejecutamos la plantilla
	plantillas.ExecuteTemplate(w, "editar", orden)
}

// Editar datos del usuario
func retirar(w http.ResponseWriter, r *http.Request) {

	// Recibimos el dato ID
	loteOrden := r.URL.Query().Get("lote")

	// Conectando base de datos
	conexionEstablecida := conecionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM Ordenes WHERE Lote=?", loteOrden)
	if err != nil {
		panic(err.Error())
	}

	// Ejecutamos la sentencia SQL, recorriendo cada registro
	orden := Orden{}
	for registro.Next() {
		var lote, tamanio_lote int
		var producto, codigo string

		err = registro.Scan(&lote, &tamanio_lote, &producto, &codigo) // Obtenemos los datos
		if err != nil {
			panic(err.Error())
		}

		orden.Lote = lote
		orden.Tamanio_lote = tamanio_lote
		orden.Producto = producto
		orden.Codigo = codigo
	}

	// Ejecutamos la plantilla
	plantillas.ExecuteTemplate(w, "retirar", orden)
}

func actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		lote := r.FormValue("lote")
		tamanio_lote := r.FormValue("tamanio_lote")
		producto := r.FormValue("producto")
		codigo := r.FormValue("codigo")

		conexionEstablecida := conecionBD()
		editarRegistro, err := conexionEstablecida.Prepare("UPDATE Ordenes SET Tamanio_lote=?, Producto=?, Codigo=? WHERE Lote=?")
		if err != nil {
			panic(err.Error())
		}

		editarRegistro.Exec(tamanio_lote, producto, codigo, lote)

		// redireccionamos a una pagina
		http.Redirect(w, r, "/", 301)
	}
}

func actualizarRetiro(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		lote := r.FormValue("lote")
		tamanio_lote := r.FormValue("tamanio_lote")
		nuevo_tamanio := r.FormValue("nuevo_tamanio")

		conexionEstablecida := conecionBD()
		editarRegistro, err := conexionEstablecida.Prepare("UPDATE Ordenes SET Tamanio_lote = ? - ? WHERE Lote=? AND (Tamanio_lote >= 0)")
		if err != nil {
			panic(err.Error())
		}

		editarRegistro.Exec(tamanio_lote, nuevo_tamanio, lote)

		// redireccionamos a una pagina
		http.Redirect(w, r, "/", 301)
	}
}
