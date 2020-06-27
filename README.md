# AWS Budgets Exporter

## How to run

1. Docker buil

	```
	$ docker build -t aws-budgets-exporter .
	```

1. Configure config.yaml
1. Docker run

	```
	$ docker run -it --rm \
		--name aws-budgets-exporter \
		-p 8080:8080 \
		-e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" \
		-e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" \
		-e "AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}" \
		-e "ROLE_SESSION_NAME=${ROLE_SESSION_NAME}" \
		-e "CONFIG_PATH=/opt/config.yaml" \
		-v $PWD/config.yaml:/opt/config.yaml \
		aws-budgets-exporter
	```

## Environment values

* LISTEN_ADDR
	* default: ":8080"
* METRICS_ENDPOINT
	* default: "/metrics"
* TIMEOUT
	* default: "5s"
* CONFIG_PATH
	* Path to config file
	* defualt: "/opt/config.yaml"
* AWS_ACCESS_KEY_ID
	* AWS access key
* AWS_SECRET_ACCESS_KEY
	* AWS secret key
* AWS_DEFAULT_REGION
	* AWS region
* ROLE_SESSION_NAME
	* An identifier for the assumed role session
