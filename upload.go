package main

import(
 	"log"
    "net/http"
    "bufio"
    "path/filepath"
    "io"
    "os"
    "encoding/json"
)


func uploadSingleFile(w http.ResponseWriter, r *http.Request){
	multipartReader, err := r.MultipartReader()	
	if err != nil {
		log.Println(err)
		return 
	}

	firstReader, err := multipartReader.NextPart()
	if err != nil {
		log.Println(err)
		return 
	}
	defer firstReader.Close()

	file := filepath.Join("./public", firstReader.FileName())
	output, err := os.Create(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer output.Close()
	
	bWriter := bufio.NewWriter(output)
	io.Copy(bWriter, firstReader)

	m := make(map[string]string)
	m["filename"] = "/static/" + firstReader.FileName()
	jData, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
