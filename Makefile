.PHONY: build_bin
build_bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/gimli_exporter_amd64 -buildvcs=false .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o build/gimli_exporter_arm64 -buildvcs=false .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/gimli_exporter.exe -buildvcs=false .

.PHONY: upload
upload: build_bin
	aws s3 cp build s3://res-download/apisec/probe/stress_test/gimli_exporter/ --recursive --acl public-read
	aws s3 cp gimli_exporter.service s3://res-download/apisec/probe/stress_test/gimli_exporter/ --acl public-read
	aws s3 cp gimli_exporter.ps1 s3://res-download/apisec/probe/stress_test/gimli_exporter/ --acl public-read

