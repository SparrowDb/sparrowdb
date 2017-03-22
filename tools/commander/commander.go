package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"

	"github.com/SparrowDb/sparrowdb/slog"
)

var (
	flagHost         = flag.String("h", "127.0.0.1", "Host")
	flagPort         = flag.Int("P", 8081, "Port")
	flagCommand      = flag.String("c", "127.0.0.1", "Command (SEND, DELETE)")
	flagDatabaseName = flag.String("db", "", "Database name")
	flagImageName    = flag.String("iname", "", "Image name")
	flagImageFolder  = flag.String("ifolder", "", "Image folder")
	flagImagePath    = flag.String("ipath", "", "Image path")
	flagImageScript  = flag.String("iscript", "", "Image script")
	flagImageUpsert  = flag.Bool("upsert", false, "Image path")
	address          string
)

const (
	version = "1.0.0"

	contentTypeJSON = "application/json"
	contentTypeForm = "multipart/form-data"
)

func httpRequest(method, urlParms, contentType string, body io.Reader) (string, string, error) {
	slog.Infof("URL: %s", fmt.Sprintf("%s/%s", address, urlParms))
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", address, urlParms), body)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	bResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	return string(bResp), resp.Status[:3], nil
}

func createFormWriter(dbname, imgName, imgPath string) (*multipart.Writer, bytes.Buffer, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, b, err
	}
	defer f.Close()

	fw, err := w.CreateFormFile("uploadfile", imgPath)
	if err != nil {
		return nil, b, err
	}

	if _, err = io.Copy(fw, f); err != nil {
		return nil, b, err
	}

	if err := w.WriteField("dbname", dbname); err != nil {
		return nil, b, err
	}

	if err := w.WriteField("key", imgName); err != nil {
		return nil, b, err
	}

	if *flagImageUpsert == true {
		if err := w.WriteField("upsert", "true"); err != nil {
			return nil, b, err
		}
	}

	if *flagImageScript != "" {
		if err := w.WriteField("script", *flagImageScript); err != nil {
			return nil, b, err
		}
	}

	w.Close()

	return w, b, nil
}

func cmdSend(dbname, imgName, imgPath string) {
	addr := fmt.Sprintf("%s/%s/%s", "api", dbname, imgName)

	if dbname == "" || imgName == "" {
		slog.Fatalf("Invalid SEND params [db:%s, image:/%s]", dbname, imgName)
	}
	fmt.Printf("%s - %s\n", dbname, imgName)

	w, b, err := createFormWriter(dbname, imgName, imgPath)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	resp, status, err := httpRequest("PUT", addr, w.FormDataContentType(), &b)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	printResponse(resp, status)
}

func cmdSendFolder(dbname, path string) {
	allowedImages := map[string]byte{
		".jpg":  0,
		".jpeg": 0,
		".gif":  0,
		".png":  0,
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if _, ok := allowedImages[ext]; ok {
			cmdSend(dbname, strings.Replace(f.Name(), ext, "", 1), filepath.Join(path, f.Name()))
		}
	}
}

func cmdDelete() {
	addr := fmt.Sprintf("%s/%s/%s", "api", *flagDatabaseName, *flagImageName)

	if *flagDatabaseName == "" || *flagImageName == "" {
		slog.Fatalf("Invalid send params [db:%s, image:/%s]", *flagDatabaseName, *flagImageName)
	}

	resp, status, err := httpRequest("DELETE", addr, contentTypeJSON, nil)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	printResponse(resp, status)
}

func cmdListDbs() {
	addr := fmt.Sprintf("%s/%s", "api", "_all")

	resp, status, err := httpRequest("GET", addr, contentTypeJSON, nil)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	printResponse(resp, status)
}

func cmdListImgs() {
	addr := fmt.Sprintf("%s/%s/%s", "api", *flagDatabaseName, "_keys")

	if *flagDatabaseName == "" {
		slog.Fatalf("Invalid database name")
	}

	resp, status, err := httpRequest("GET", addr, contentTypeJSON, nil)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	printResponse(resp, status)
}

func printResponse(resp, status string) {
	var d bytes.Buffer
	if err := json.Indent(&d, []byte(resp), " ", " "); err != nil {
		slog.Fatalf(err.Error())
	}

	msg := fmt.Sprintf("[%s] %s", status, d.String())

	if status == "200" {
		slog.Infof(msg)
	} else {
		slog.Errorf(msg)
	}
}

func main() {
	flag.Parse()
	address = fmt.Sprintf("%s://%s:%v", "http", *flagHost, *flagPort)

	slog.Infof("SparrowDb Commander %s - commander tool", version)

	switch strings.ToLower(*flagCommand) {
	case "send":
		cmdSend(*flagDatabaseName, *flagImageName, *flagImagePath)
	case "sendf":
		cmdSendFolder(*flagDatabaseName, *flagImagePath)
	case "delete":
		cmdDelete()
	case "dbs":
		cmdListDbs()
	case "imgs":
		cmdListImgs()
	default:
		flag.Usage()
	}
}
