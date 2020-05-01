outObj := fjs_go_oss

all : 
	# go build ./
	go build -o ./bin/$(outObj) ./main/main.go
	./bin/$(outObj)

clear :
	rm -rf ./bin/$(outObj)