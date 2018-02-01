
all: certs.go

server.key:
	openssl genrsa -out server.key 2048

server.pem: server.key
	openssl req -new -x509 -key server.key -out server.pem -days 3650 -subj "/CN=localhost:8888/emailAddress=admin@localhost/C=US/ST=Utah/L=Provo/O=SlickQA/OU=Slick"

jwt.key:
    openssl genrsa -out jwt.key 1024

jwt.pub: jwt.key
    openssl rsa -in jwt.key -pubout > jwt.pub

certs.go: server.pem server.key jwt.key jwt.pub
	echo "package certs\n\nconst (\n\tKey = \`" >certs.go
	cat server.key >> certs.go
	echo "\`\n\n\tCert = \`" >> certs.go
	cat server.pem >> certs.go
	echo "\`\n)\n" >> certs.go

