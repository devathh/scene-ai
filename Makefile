all: keys

.PHONY: all keys

keys:
	openssl genrsa -out jwtRS256.key 2048
	openssl rsa -in jwtRS256.key -pubout -out jwtRS256.key.pub

tlslocal:
	mkcert -install

	mkdir -p ./certs
	mkcert localhost 127.0.0.1
	mv localhost+1.pem server.crt
	mv server.crt ./certs/
	mv localhost+1-key.pem server.key
	mv server.key ./certs/