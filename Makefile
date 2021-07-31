swagger:
	GO111MODULE=off swagger generate spec -o ./api/swagger.yaml --scan-models 
run_mongo:
	docker run --name mongo -p 27017:27017 -d -v /data/db mongo:latest