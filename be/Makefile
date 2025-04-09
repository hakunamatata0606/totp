target = .
tmp := $(shell mktemp -d)

CC = gcc
CFLAG = -fPIC -shared -Ilibc/

all: server libotp libotp_test
	mkdir -p $(target)
	cp $(tmp)/server $(target)/server
	cp $(tmp)/libotp.so $(target)/libotp.so
	cp $(tmp)/otp_test $(target)/otp_test
migrate_up:
	migrate -path db/migrations -database "mysql://root:123@tcp(127.0.0.1:3306)/totp" -verbose up
migrate_down:
	migrate -path db/migrations -database "mysql://root:123@tcp(127.0.0.1:3306)/totp" -verbose down
sqlc:
	sqlc generate

server: libotp main.go
	go build -o $(tmp)/server main.go

libotp: libc/otp.c
	$(CC) $(CFLAG) libc/otp.c -o $(tmp)/libotp.so -lssl -lcrypto -lm -lpthread

libotp_test: libotp libc/otp_test.c
	$(CC) -Ilibc/ libc/otp_test.c -o $(tmp)/otp_test -L$(tmp)/ -lotp
clean:
	rm -f $(target)/server
.PHONY: migrate_up migrate_down sqlc all server libotp libotp_test clean