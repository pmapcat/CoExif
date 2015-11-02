package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type ExifProc struct {
	proc   *exec.Cmd
	stdin  chan string
	stdout chan string
	stderr chan string
	Busy   bool
}

type StandartJSON map[string]interface{}

type IdioticJSON struct {
	Items []map[string]interface{}
}

func (ex *ExifProc) Init() {
	ex.Busy = false
	ex.stdin = make(chan string)
	ex.stdout = make(chan string)
	ex.stderr = make(chan string)
}

func (ex *ExifProc) StartProc() {
	cmd := exec.Command(EXIF_DIR, "-stay_open", "True", "-@", "-")
	ex.proc = cmd
	log.Println("Spawning Exiftool proc")
	go long_running_proc(ex.stdin, ex.stdout, ex.stderr, cmd)
}

func (ex *ExifProc) KillProc() {
	log.Println("Killing Exiftool Proc")
	if err := ex.proc.Process.Kill(); err != nil {
		log.Println("failed to kill: ", err)
	}
}

func (ex *ExifProc) GETMeta(path string, probably_filter ...string) IdioticJSON {
	ex.Busy = true
	filtering := ""
	if len(probably_filter) != 0 {
		filtering = "\n-" + strings.Join(probably_filter, "\n -")
	}
	fmt.Printf("%s%s\n-JSON\n-fast\n-execute\n", path, filtering)

	ex.stdin <- fmt.Sprintf("%s%s\n-JSON\n-fast\n-execute\n", path, filtering)
	returnJSON := IdioticJSON{}
	message := make([]string, 0)
	for {
		msg := <-ex.stdout
		if strings.Contains(msg, "{ready}") {
			log.Println("METADATA COLLECTED")
			break
		}
		message = append(message, msg)
	}

	// we must hack some json surroundings to allow go parse this json
	message = append([]string{`{"Items":`}, message...)
	message = append(message, "}")
	// ========
	full_message := strings.Join(message, "")
	json.Unmarshal([]byte(full_message), &returnJSON)
	ex.Busy = false
	return returnJSON
}
func (ex *ExifProc) UPDATEMeta(file_path string, datum StandartJSON) error {
	// exiftool -json=a.json a.png -overwrite_original
	ex.Busy = true
	// ====== MERGE MAPS ======
	old_meta := ex.GETMeta(file_path).Items
	new_corpus := datum
	var old_corpus StandartJSON
	if len(old_meta) > 0 {
		old_corpus = old_meta[0]
	} else {
		ex.Busy = false
		return errors.New("No suitable corpus for file")
	}
	mergo.Map(&new_corpus, old_corpus)
	// ====== WRITE JSON ======
	new_corpus_byte, err := json.Marshal(new_corpus)
	if err != nil {
		ex.Busy = false
		return err
	}
	// in case there is no temp dir
	err = os.MkdirAll(TEMP_DIR, 0777)
	if err != nil {
		ex.Busy = false
		return err
	}

	// ====== CREATE TEMP JSON FILE =====
	json_path := path.Join(TEMP_DIR, uuid.NewV4().String()+".json")

	err = ioutil.WriteFile(json_path, new_corpus_byte, 0644)
	if err != nil {
		ex.Busy = false
		return err
	}
	// ====== POST DATUM FOR WORKAGE =====
	// exiftool -json=a.json a.png -overwrite_original

	ex.stdin <- fmt.Sprintf("%s\n-json=%s\n-overwrite_original\n-execute\n", file_path, json_path)
	msg := <-ex.stdout
	log.Println(msg)
	// ====== DELETE STALE FILE ====
	err = os.Remove(json_path)
	if err != nil {
		ex.Busy = false
		return err
	}
	ex.Busy = false
	return nil
}
