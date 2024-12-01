run: # binary in heroku dyno
	bin/crypto_dca_go

dev:
	source conf/dev.env && go run main.go

test:
	go test -count=1 ./...

prod_logs:
	heroku logs --remote production

ssh_staging:
	heroku run bash -a gemini-dca-staging

ssh_prod:
	heroku run bash -a gemini-dca-production