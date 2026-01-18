.PHONY: build dev clean

build: build-frontend build-backend

build-frontend:
	cd web && npm install && npm run build

build-backend:
	go build -o skill-router .

dev:
	@echo "Run 'go run .' in one terminal"
	@echo "Run 'cd web && npm run dev' in another terminal"

clean:
	rm -rf skill-router web/dist web/node_modules
