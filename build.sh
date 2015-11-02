function compile_go(){
    GOOS=$3 GOARCH=$1 go build -o $2
    mv $2 ./build/
}
# compile_go 386   metator_86_1.1.exe     windows
# compile_go amd64 metator_64_1.1.exe     windows
compile_go 386   coexif_86_linux_1.1   linux
compile_go amd64 coexif_64_linux_1.1   linux
compile_go 386   coexif_86_osx_1.1     darwin
compile_go amd64 coexif_64_osx_1.1     darwin
