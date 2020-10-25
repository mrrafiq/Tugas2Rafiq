package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

var db *sql.DB
var err error

type yamlconfig struct {
	Connection struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		User     string `yaml:"user"`
		Database string `yaml:"database"`
	}
}

type mahasiswa struct {
	NoBp     int    `xml:"MahasiswaID"`
	Nama     string `xml:"Nama"`
	Fakultas string `xml:"Fakultas"`
	Jurusan  string `xml:"Jurusan"`
	Alamat   struct {
		Jalan     string `xml:"Jalan"`
		Kelurahan string `xml:"Kelurahan"`
		Kecamatan string `xml:"Kecamatan"`
		Kabupaten string `xml:"Kabupaten"`
		Provinsi  string `xml:"Provinsi"`
	} `xml:"Alamat"`
	Nilai []nilai `json:"Nilai"`
}

type nilai struct {
	NoBp       int     `xml:"MahasiswaID"`
	IDMatkul   int     `xml:"MataKuliahID"`
	NamaMatkul string  `xml:"mataKuliah"`
	Nilai      float64 `xml:"Nilai"`
	Semester   string  `xml:"Semester"`
}

// Get all orders

func getNilai(w http.ResponseWriter, r *http.Request) {
	var mhs mahasiswa
	var nilaix nilai

	params := mux.Vars(r)

	sql := `SELECT
			MahasiswaID,
				IFNULL(nama,'') nama,
				IFNULL(jalan,'') jalan,
				IFNULL(kelurahan,'') kelurahan,
				IFNULL(kecamatan,'') kecamatan,
				IFNULL(kabupaten,'') kabupaten,
				IFNULL(provinsi,'') provinsi,
				IFNULL(fakultas,'') fakultas,
				IFNULL(jurusan,'') jurusan				
			FROM mahasiswa WHERE MahasiswaID IN (?)`

	result, err := db.Query(sql, params["MahasiswaID"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		err := result.Scan(&mhs.NoBp, &mhs.Nama, &mhs.Alamat.Jalan, &mhs.Alamat.Kelurahan, &mhs.Alamat.Kecamatan, &mhs.Alamat.Kabupaten, &mhs.Alamat.Provinsi, &mhs.Fakultas, &mhs.Jurusan)

		if err != nil {
			panic(err.Error())
		}

		sqlNilai := `SELECT
					MahasiswaID 		
						, mata_kuliah.MataKuliahID
						, mata_kuliah.mataKuliah
						, nilai
						, semester
					FROM
						nilai INNER JOIN mata_kuliah 
							ON (nilai.MataKuliahID = mata_kuliah.MataKuliahID)
					WHERE MahasiswaID = ?`

		noBp := &mhs.NoBp
		fmt.Println(noBp)
		resultNilai, errNilai := db.Query(sqlNilai, noBp)

		defer resultNilai.Close()

		if errNilai != nil {
			panic(err.Error())
		}

		for resultNilai.Next() {
			err := resultNilai.Scan(&nilaix.NoBp, &nilaix.IDMatkul, &nilaix.NamaMatkul, &nilaix.Nilai, &nilaix.Semester)
			if err != nil {
				panic(err.Error())
			}
			mhs.Nilai = append(mhs.Nilai, nilaix)
		}
	}
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	xml.NewEncoder(w).Encode(mhs)
}

// Main function
func main() {
	yamlFile, err := ioutil.ReadFile("../Yaml/config.yml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}
	var yamlConfig yamlconfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	host := yamlConfig.Connection.Host
	port := yamlConfig.Connection.Port
	user := yamlConfig.Connection.User
	pass := yamlConfig.Connection.Password
	data := yamlConfig.Connection.Database

	var (
		mySQL = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, data)
	)

	db, err = sql.Open("mysql", mySQL)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/nilai/{MahasiswaID}", getNilai).Methods("GET")
	// r.HandleFunc("/nilai", getNilaiAll).Methods("GET")

	fmt.Println("Server on :8282")
	// Start server
	log.Fatal(http.ListenAndServe(":8282", r))
}
