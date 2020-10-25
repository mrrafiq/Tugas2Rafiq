package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type mahasiswa struct {
	NoBp     int    `json:"MahasiswaID"`
	Nama     string `json:"Nama"`
	Fakultas string `json:"Fakultas"`
	Jurusan  string `json:"Jurusan"`
	Alamat   struct {
		Jalan     string `json:"Jalan"`
		Kelurahan string `json:"Kelurahan"`
		Kecamatan string `json:"Kecamatan"`
		Kabupaten string `json:"Kabupaten"`
		Provinsi  string `json:"Provinsi"`
	} `json:"Alamat"`
	Nilai []nilai `json:"Nilai"`
}

type nilai struct {
	NoBp       int     `json:"MahasiswaID"`
	IDMatkul   int     `json:"MataKuliahID"`
	NamaMatkul string  `json:"mataKuliah"`
	Nilai      float64 `json:"Nilai"`
	Semester   string  `json:"Semester"`
}

func main() {

	url := "http://localhost:8282/nilai"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	mhs := mahasiswa{}
	jsonErr := json.Unmarshal(body, &mhs)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(mhs.NoBp)
	fmt.Println(mhs.Nama)
	fmt.Println(mhs.Alamat.Jalan)
	fmt.Println(mhs.Alamat.Kelurahan)
	fmt.Println(mhs.Alamat.Kecamatan)
	fmt.Println(mhs.Alamat.Kabupaten)
	fmt.Println(mhs.Alamat.Provinsi)

	for _, nilai := range mhs.Nilai {
		fmt.Println("No BP", nilai.NoBp)
		fmt.Println("ID Mata Kuliah", nilai.IDMatkul)
		fmt.Println("Nama Mata Kuliah", nilai.NamaMatkul)
		fmt.Println("Nilai", nilai.Nilai)
		fmt.Println("Semester", nilai.Semester)
	}

}
