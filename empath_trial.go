package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const URL = `https://api.webempath.net/v2/analyzeWav`

type Payload struct {
	ApiKey string `json:"apikey"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	apikey := os.Getenv("EMPATH_API_KEY")
	wavPath := "./sample.wav"

	file, err := os.Open(wavPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("wav", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(part, file); err != nil {
		log.Fatal(err)
	}

	if err = writer.WriteField("apikey", apikey); err != nil {
		log.Fatal(err)
	}

	if err = writer.Close(); err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(URL, writer.FormDataContentType(), body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes))
}
