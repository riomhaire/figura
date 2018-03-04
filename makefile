.DEFAULT_GOAL := everything

dependencies:
	@echo Downloading Dependencies
	@go get ./...

build: dependencies
	@echo Compiling Apps
	@echo   --- figura server
	@go build github.com/riomhaire/figura/frameworks/application/figura
	@go install github.com/riomhaire/figura/frameworks/application/figura
	@echo Done Compiling Apps
test:
	@echo Running Unit Tests
	@go test ./...

profile:
	@echo Profiling Code
	@go get -u github.com/haya14busa/goverage 
	@goverage -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out	

# @go test -coverprofile coverage.out  github.com/riomhaire/figura/interfaces 
# @go tool cover -html=coverage.out -o coverage-interfaces.html
# @rm coverage.out	
# @go test -coverprofile coverage.out  github.com/riomhaire/figura/usecases
# @go tool cover -html=coverage.out -o coverage-usecases.html
# @rm coverage.out	
# @go test -coverprofile coverage.out  github.com/riomhaire/figura/entities
# @go tool cover -html=coverage.out -o coverage-entities.html
# @rm coverage.out

clean:
	@echo Cleaning
	@go clean
	@rm -f figura 
	@rm -f coverage-*.html
	@find . -name "debug.test" -exec rm -f {} \;

everything: clean build test profile  
	@echo Done
