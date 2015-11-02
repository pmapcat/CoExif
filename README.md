**CoExif** is a fast cacheable **JSON** **REST** File Metadata Server,
which supports nearly any file format.
It uses one of the best metadata management tools outhere [Exiftool](http://www.sno.phy.queensu.ca/~phil/exiftool/)
Internally, **CoExif** uses a pool of long running ExifTool processes.

#Getting Started

## Installation
Download latest version of [Exiftool](http://www.sno.phy.queensu.ca/~phil/exiftool/).
Download one of these builds:
* [coexif_86_linux_1](https://github.com/MichaelLeachim/CoExif/releases/download/v1.0/coexif_86_linux_1.1) 
* [coexif_64_linux_1](https://github.com/MichaelLeachim/CoExif/releases/download/v1.0/coexif_64_linux_1.1) 
* [coexif_86_osx_1  ](https://github.com/MichaelLeachim/CoExif/releases/download/v1.0/coexif_86_osx_1.1) 
* [coexif_64_osx_1  ](https://github.com/MichaelLeachim/CoExif/releases/download/v1.0/coexif_64_osx_1.1)

## Start Server
```bash
co_exif -root "/home/mik/" -port 9200 -auth-name="admin" -auth-pass="admin"
```

Now, you can query server however you like.

## GET
```bash
curl -u admin:admin 127.0.0.1:9200/2.jpg
{
  "Items": [
    {
      "APP14Flags0": "[14]",
      "APP14Flags1": "(none)",
      "AlreadyApplied": true,
      "ApplicationRecordVersion": 0,
      "AutoLateralCA": 0,
      "BitsPerSample": 8,
      "Blacks2012": "+14",
      "BlueHue": 0,
      "BlueMatrixColumn": "0.14307 0.06061 0.7141",
      "BlueSaturation": 0,
      "BlueTRC": "(Binary data 2060 bytes, use -b option to extract)",
      "Brightness": "+50",
      "CMMFlags": "Not Embedded, Independent",
      "CameraProfile": "Adobe Standard",
      ......
      }
   ]
}
```
## GET many

```bash
curl -u admin:admin http://127.0.0.1:9200/
{
  "Items": [
    {
      "BitDepth": 24,
      "Compression": "FMP4",
      "Directory": "/home/mik",
      "Duration": "0:01:21",
      "ExifToolVersion": 10.04,
      "FileAccessDate": "2015:08:02 20:10:40+03:00",
      "FileInodeChangeDate": "2015:08:03 13:35:03+03:00",
      "FileModifyDate": "2015:08:03 13:35:03+03:00",
      "FileName": "output.avi",
      "FilePermissions": "rw-r--r--",
      ......
      e.t.c.
      }
  ]
}
```
## FILTERED GET

```bash
curl -u admin:admin "127.0.0.1:9200/2.jpg?tags=CameraProfile&tags=Brightness"
{
  "Items": [
    {
      "Brightness": "+50",
      "CameraProfile": "Adobe Standard",
      "SourceFile": "/home/mik/2.jpg"
    }
  ]
}
```
## Filtered many
It is good to have been able to select  specific tags:
```bash
curl -u admin:admin "http://127.0.0.1:9999/home/mik/?tags=FileAccessDate&tags=FileName"
{
  "Items": [
    {
      "FileAccessDate": "2015:08:02 20:10:40+03:00",
      "FileName": "output.avi",
      "SourceFile": "/home/mik/output.avi"
    },
    {
      "FileAccessDate": "2015:11:02 15:24:44+03:00",
      "FileName": "DSC_0158.NEF",
      "SourceFile": "/home/mik/DSC_0158.NEF"
    },
  ]
}    
```

## POST metadata
Post will replace metadata in a file with specified fields.
Not specified fields will remain the same.
Example:
```bash
// POST DATA
curl -XPOST http://127.0.0.1:9999/a.png -u admin:admin -H "Content-Type: application/json" -d '{"Artist":"Mik-s picture","Author":"Blablablab"}'
// RESULTS
curl -u admin:admin "127.0.0.1:9999/a.png?tags=Artist&tags=Author"
{
  "Items": [
    {
      "Artist": "Mik-s picture",
      "Author": "Blablablab",
      "SourceFile": "/home/mik/a.png"
    }
  ]
}
```
# SERVER PARAMS
```bash
  -auth-name="admin": Enter auth name
  -auth-pass="admin": Enter auth pass
  -auto-spawn=false: Should I autospawn processes
  -exif-path="./exif_tool/exiftool": Enter path to exiftool
  -max-prox=10: Enter number of ExifTool processes
  -port="9999": Enter a server port number
  -root="/": Enter default root path
```
  
  
<!-- # Build from source -->
<!-- ``` -->
<!-- git clone michaelleachim/coexif; -->
<!-- cd coexif; -->
<!-- ``` -->

<!-- # Bindings -->
<!-- Python -->
<!-- GoLang -->
<!-- Node.JS -->
<!-- Clojure -->
