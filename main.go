package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/catinello/base62"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", handlerFunc)

	errer := http.ListenAndServe(":8080", nil)
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	if errer != nil {
		log.Fatalln(errer)
	}
	server.ListenAndServe()
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	query := r.FormValue("key")
	number, _ := base62.Decode(r.URL.Path[1:])
	fmt.Println(number)
	if number > 0 {
		link := dbSelect(&number)
		fmt.Println("ok")
		if link != "" {
			fmt.Println(link)
			http.Redirect(w, r, link, 302)
		}
	}
	//유효성 검사
	var str string
	var str2 string
	if len(query) > 8 {
		str = query[0:8]
		str2 = query[0:7]
	}

	if (str == "https://") || (str2 == "http://") {
		uid := dbfunc(&query)

		var url string
		url = "http://127.0.0.1:8080/" + base62.Encode(uid)
		t.Execute(w, url)
	} else {
		t.Execute(w, "URL을 다시입력해 주세요")
	}

}

func dbfunc(query *string) int {
	//db 연결
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3307)/study_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 쿼리
	result, err := db.Exec("INSERT INTO url_tb(url) VALUE(?)", query)
	if err != nil {
		log.Fatal(err)
	}

	n, err := result.RowsAffected()
	if n == 1 {
		fmt.Println("1 row inserted.")
	}
	var name int
	errr := db.QueryRow("SELECT max(id) from url_tb").Scan(&name)
	if errr != nil {
		panic(err.Error())
	}

	return name
}

func dbSelect(id *int) string {
	//db 연결
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3307)/study_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var name string
	errr := db.QueryRow("SELECT url FROM url_tb WHERE id=?", id).Scan(&name)
	if errr != nil {
		name = ""
	}

	return name
}
