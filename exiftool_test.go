package main

import (
	// "log"
	"testing"
)

func Test_exif_tool_GET(t *testing.T) {
	exp := ExifProc{}
	exp.Init()
	exp.StartProc()
	datum := exp.GETMeta("./co_exif")
	if datum.Items[0]["FileType"] != "ELF executable" {
		t.Error(datum.Items[0]["FileType"], "!=", "ELF executable")
	}
	exp.KillProc()
}

func Test_exif_tool_UPDATE(t *testing.T) {
	exp := ExifProc{}
	exp.Init()
	exp.StartProc()
	// equivalent in SHELL
	// echo "./temp/a.png\n -json=temp/cf28c49a-9cb2-42e1-a344-147294757e02.json\n-execute\n" | ./exif_tool/exiftool -stay_open True -@ -

	updated_datum := make(StandartJSON, 10)
	updated_datum["Artist"] = "Navar"
	updated_datum["Author"] = "Navarov"
	updated_datum["Comment"] = "Lorem"
	updated_datum["Copyright"] = "Ipsum"
	err := exp.UPDATEMeta("./temp/b.png", updated_datum)
	if err != nil {
		t.Error(err)
	}
	datum_updated := exp.GETMeta("./temp/b.png").Items

	test_helper := func(field string) {
		if datum_updated[0][field] != updated_datum[field] {
			t.Error(datum_updated[0][field], "!=", updated_datum[field])
		}
	}
	test_helper("Artist")
	test_helper("Author")
	test_helper("Comment")
	test_helper("Copyright")
	exp.KillProc()
}
