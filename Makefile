build:
	go build github.com/paoloconi96/invoice-parser/cmd/api-server

run:
	go run github.com/paoloconi96/invoice-parser/cmd/api-server

next_dev:
	echo "== Next dev =="
	cd web && npm run dev
