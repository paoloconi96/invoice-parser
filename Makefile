react_dev:
	echo "== React dev =="
	cd web && npm run dev
go_dev:
	echo "== Go dev =="
	go run github.com/paoloconi96/invoice-parser/cmd/api-server
