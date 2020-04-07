package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type Data struct {
	Title    string
	Subtitle string
	Body     string
}

func main() {
	templateReader, err := zip.OpenReader("sample-template.docx")
	if err != nil {
		log.Fatalf("cannot open template file: %s", err.Error())
	}
	defer templateReader.Close()

	templateStr := getTemplateString(templateReader, err)
	tpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		log.Fatalf("cannot parse template file: %s", err.Error())
	}

	generateDoc(templateReader, tpl)

	log.Println("Success!")
}

func generateDoc(templateReader *zip.ReadCloser, tpl *template.Template) {
	resultFile, err := os.Create("result.docx")
	if err != nil {
		log.Fatalf("cannot create result file: %s", err.Error())
	}
	defer resultFile.Close()
	resultWriter := zip.NewWriter(resultFile)
	defer resultWriter.Close()

	for _, f := range templateReader.File {
		w, err := resultWriter.Create(f.Name)
		if err != nil {
			log.Fatalf("cannot create file for result docx: %s", err.Error())
		}
		if f.Name == "word/document.xml" {
			data := Data{
				Title:    "The Title",
				Subtitle: "Some description in subtitle",
				Body:     "A paragraph between lorem ipsum texts",
			}
			err = tpl.Execute(w, data)
			if err != nil {
				log.Fatalf("cannot execute template: %s", err.Error())
			}
		} else {
			fileContent, err := f.Open()
			if err != nil {
				log.Fatalf("cannot open file: %s", err.Error())
			}
			_, err = io.Copy(w, fileContent)
			if err != nil {
				log.Fatalf("cannot copy component: %s", err.Error())
			}
		}
	}
}

func getTemplateString(templateReader *zip.ReadCloser, err error) string {
	var templateFile *zip.File
	for _, f := range templateReader.File {
		if f.Name == "word/document.xml" {
			templateFile = f
		}
	}
	if templateFile == nil {
		log.Fatalf("word file not found")
	}

	docReader, err := templateFile.Open()
	if err != nil {
		log.Fatalf("cannot open document.xml: %s", err.Error())
	}
	defer docReader.Close()

	b, err := ioutil.ReadAll(docReader)
	if err != nil {
		log.Fatalf("cannot read document.xml: %s", err.Error())
	}
	templateStr := string(b)
	return templateStr
}
