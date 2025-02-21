minikube start --memory=4096

minikube mount /Users/maverick/go/src/go-microservices:/go-microservices

k9s -n kafka # Run k9s in kafka namespace

kubectl create namespace kafka

kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka

kubectl apply -f https://strimzi.io/examples/latest/kafka/kraft/kafka-single-node.yaml -n kafka # Apply the `Kafka` Cluster CR file

kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n kafka 

minikube addons enable ingress

kubectl get pods -n ingress-nginx

# Add 127.0.0.1 payment.com to /etc/hosts

minikube tunnel
