.PHONY: all clean

all: run-go run-python

run-go:
	go run main.go

run-python:
	python preprocessing.py

clean:
	rm -f dataTHPT2023.csv
install:
	go get
	pip install -r requirements.txt
