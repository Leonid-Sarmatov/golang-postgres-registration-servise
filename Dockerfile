FROM golang:1.20

WORKDIR /html_registration_web_site

COPY go.mod .
COPY . .

CMD ["go", "run", "reg_resiver_server.go"]
