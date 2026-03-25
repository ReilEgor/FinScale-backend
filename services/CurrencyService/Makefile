APP_MAIN = cmd/main.go
SWAG_OUT = api/docs

.PHONY: swag-init
swag-init: # swag init -g cmd/main.go -o api/docs --parseDependency --parseInternal
	@echo "Generating Swagger documentation..."
	swag init -g $(APP_MAIN) -o $(SWAG_OUT) --parseDependency --parseInternal
	@echo "Done! Docs are in $(SWAG_OUT)"
