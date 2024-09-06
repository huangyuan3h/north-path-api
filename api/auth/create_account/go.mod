module north-path.it-t.xyz/controller/auth/create_account

go 1.22.0

require github.com/aws/aws-lambda-go v1.46.0

require (
	north-path.it-t.xyz/utils v1.0.0
	github.com/go-playground/validator/v10 v10.19.0
)

require (
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace north-path.it-t.xyz/utils => ../../../utils
