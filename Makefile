clean:
	- rm ./receita

build/win:
	- GOOS=windows GOARCH=386 go build -o receita.exe main.go

build:
	- go build -o receita main.go

run: clean build
	- ./receita
