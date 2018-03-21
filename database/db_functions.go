package database


import (
	"net/http"
	"fmt"
	_ "database/sql"
	_"github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"html/template"
	"github.com/gorilla/mux"
	"net"
	"encoding/base64"

	"io/ioutil"
)



var db, _ = sql.Open("sqlite3", "test_db.db")
type Payload struct {
	id int
	name string
	content_type string
	host_blacklist	string
	host_whitelist	string
	data_file	string
	data_b64	string
	type_id 		int
	//one_liner	string

}

type Host struct {
	name string
	htype string
	data string
}

func CreateTable()  {

	// Insert payload query
	//insert into payloads values (NULL,'test123','wtf',NULL,NULL,'/tmp/shit',NULL,1);
	// Payload types query
	//insert into payload_types values (NULL,'javascript',NULL)



	// This function will create the requrired shits
	createTableSql := `CREATE TABLE payloads (
							id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
							name	TEXT NOT NULL UNIQUE,
							content_type	TEXT,
							host_blacklist	TEXT,
							host_whitelist	TEXT,
							data_file	TEXT,
							data_b64	TEXT,
							type_id	INTEGER NOT NULL
						);
						`

	createHostSql := `CREATE TABLE hosts (
							id	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
							type	TEXT NOT NULL,
							data	TEXT NOT NULL
						);`

	createTypesSql := `CREATE TABLE types (
							id	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
							name	TEXT NOT NULL UNIQUE,
							type_template	TEXT
						);`



	fmt.Println(createTableSql)
	fmt.Println(createHostSql)
	fmt.Println(createTypesSql)
}


func EditPayload(w http.ResponseWriter,r *http.Request)  {
	vars := mux.Vars(r)
	pid := vars["pid"]
	ReadPayload := "SELECT * FROM payloads WHERE p_id=?"
	rows, err := db.Query(ReadPayload, pid)
	if err != nil {
		panic(err)
	}
	payload := Payload{}
	rows.Next()
	err_sql := rows.Scan(&payload.id,&payload.name,&payload.content_type,&payload.data_b64)

	if err_sql != nil{
		panic(err_sql)
	}
	t,err_tmpl := template.ParseFiles("templates/edit.html")
	if err_tmpl != nil{
		panic(err_tmpl)
	}
	t.Execute(w,payload)
}
func DeletePayload(w http.ResponseWriter,r *http.Request)  {

}
func ShowShit(w http.ResponseWriter,r *http.Request)  {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w,"What you received is shit")
}
func GetPayloads(w http.ResponseWriter,r *http.Request)  {


}
func GetPayload(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	puid := vars["puid"]
	GetPayloadQuery := `SELECT id,
						name,
						content_type,
						COALESCE(host_blacklist, '') as host_blacklist, 
						COALESCE(host_whitelist, '') as host_whitelist,
						COALESCE(data_file, '') as data_file, 
						COALESCE(data_b64, '') as data_b64 ,
						type_id 
						from payloads 
						WHERE id=?`


	rows, err := db.Query(GetPayloadQuery, puid)
	if err != nil {
		panic(err)
	}
	payload := Payload{}
	rows.Next()
	err_sql := rows.Scan(&payload.id,&payload.name,&payload.content_type,&payload.host_blacklist,&payload.host_whitelist,&payload.data_file,&payload.data_b64,&payload.type_id)

	if err_sql != nil{
		panic(err_sql)
	}

	ip , _ , _ := net.SplitHostPort(r.RemoteAddr)

	log.Println(fmt.Sprintf("Delivering payload %s to IP : %s",payload.name,ip))

	w.Header().Set("Content-Type",payload.content_type)
	w.WriteHeader(http.StatusOK)

	if payload.data_file == ""{
		if payload.data_b64 != ""{
			data, err := base64.StdEncoding.DecodeString(payload.data_b64)
			if err != nil{
				log.Println("ERROR : Decoding b64 payload failed.")
			}
			w.Write([]byte(data))
		}else{
			log.Println("ERROR : Payload delivery failed. No content or file specified.")
		}
	}else{
		// Write data from file
		data, err := ioutil.ReadFile(payload.data_file)
		if err != nil{
			log.Println(fmt.Sprintf("ERROR: Payload file %s not found.", payload.data_file))
			return
		}
		w.Write(data)
	}


}


func CreatePayloadGet(w http.ResponseWriter,r *http.Request)  {
	w.WriteHeader(http.StatusOK)
	t, err_tmpl := template.ParseFiles("templates/create.html")
	if err_tmpl != nil{
		panic(err_tmpl)
	}
	t.Execute(w,nil)
}

func CreatePayload(w http.ResponseWriter,r *http.Request)  {

	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare("INSERT INTO payloads VALUES (NULL,?,?,?);")
	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	_, err := stmt.Exec(r.FormValue("p_name"),r.FormValue("p_ct"),r.FormValue("p_content"))
	tx.Commit()
	if err != nil {
		log.Println("Payload Insertion Failed")
		w.WriteHeader(http.StatusBadRequest)
	}else{
		http.Redirect(w,r,"/create.html",http.StatusSeeOther)
		log.Println("Payload created.")


	}

}