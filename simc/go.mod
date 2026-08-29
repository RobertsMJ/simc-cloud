module github.com/RobertsMJ/simc-cloud/simc

go 1.26.7

replace github.com/RobertsMJ/simc-cloud/test-utils => ./../test-utils

require (
	github.com/RobertsMJ/simc-cloud/test-utils v0.0.0-00010101000000-000000000000
	github.com/samber/lo v1.53.0
	github.com/stretchr/testify v1.12.1
)

require (
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/text v0.22.0 // indirect
)
