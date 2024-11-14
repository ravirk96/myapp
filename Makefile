.PHONY: all build deploy clean

all: build deploy

build:
	./build.sh
	
deploy:
	cd terraform && terraform init && terraform apply -auto-approve

clean:
	rm -f function.zip
	cd terraform && terraform destroy -auto-approve