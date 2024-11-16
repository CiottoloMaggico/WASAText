swagger-cli bundle doc/api/template.yaml --outfile doc/api.yaml --type yaml
lint-openapi doc/api.yaml -r linters/openApi/spectral.js
