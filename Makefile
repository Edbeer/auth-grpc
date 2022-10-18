lint:
	cd proto && buf lint

gen: lint
	cd proto && buf generate