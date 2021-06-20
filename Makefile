README.md: *.go
	goreadme -badge-codecov -badge-godoc -badge-goreportcard -constants -variabless -functions -types -factories -methods -import-path github.com/jspc/bottom > README.md
