# go-cache-kubernete

The application is a Data Caching service designed and implemented using Microservices architecture. The cloud deployment environments are used **Kubernetes**, **Docker**, and written in **Go** programming language. The application also uses a **MongoDB** as NoSQL database with **Redis** in-memory database for the caching services.

## Features

 - **MongoDB**: The MongoDB Go driver is implemented to perform several database operations. The installation can be done using the go dependency module.
		
		go get go.mongodb.org/mongo-driver/mongo
	link: https://github.com/mongodb/mongo-go-driver

 - **Redis Cache**: The **go-redis** library is implemented to integrate the Redis data caching in the application. So, the redis will cache the second GET request while reading the user details.
 
		go get github.com/go-redis/redis/v8
	link : https://github.com/go-redis/redis

 - **Kafka Message Broker**: The confluent-kafka-go is used as a Go client library for Kafka message broker. The library will provide **Producer** and **Consumer** architecture to stream messages to the user for a subscribed topic. So, there will be two REST APIs that the user can use for Producing the messages reading from MongoDB and Consuming or Reading messages from the message broker.

		go get github.com/confluentinc/confluent-kafka-go/kafka
	link: https://github.com/confluentinc/confluent-kafka-go

	Note: It's recommended to install **confluent-kafka-go v1.4.0**, as the **librdkafka** will come with the bundle and no need to install separately.

## Kubernetes tools
Kubernetes provides several tools that can be useful to setup Kubernetes in the local environment.

 - **minikube**: The tool will run a single-node Kubernetes cluster running inside a Virtual Machine. Virtualization has to be supported in the computer and Hypervisor needed to be enabled.
	    
	 **Installation**
	    The minikube installation follows with the Hypervisor installation and [Hyperkit](https://minikube.sigs.k8s.io/docs/drivers/hyperkit/) is the recommended virtualization toolkit.   

		sudo install minikube
	link: https://kubernetes.io/docs/setup/learning-environment/minikube/
**Quick check**
				minikube start

 - **kubectl**: The kubectl command-line tool will work to manage a Kubernetes cluster. The tool will be used to deploy, create, analyze, inspect pods that are running under a Kubernetes cluster.


	**Installation**
	
	```curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" ```

	link: https://kubernetes.io/docs/tasks/tools/install-kubectl/


## Build the Docker images

The application uses Docker for container-based development. The docker image gets published as a public repository at Docker Hub.

 - **Build the image**
	
		$ docker build -t go-cache-poc .

 - **Tag the image**

		$ docker tag go-cache-poc deeptiman1991/go-cache-poc-app:1.0.0

 - **Login to docker hub**

		$ docker login

	Type Username and Password to complete the authentication

 - **Push the image to docker hub**

		$ docker push deeptiman1991/go-cache-poc:1.0.0

## Kubernetes Deployment
There will be several deployments, services that need to be running in the cluster as a Pod. The creation of a Pod requires a YAML file that will specify the kind, spec, containerPort, metadata, volume, and more. So, these parameters will be used to provide resources to the Kubernetes cluster.

**Start minikube**
To begin the deployment process start the minikube 
	
	$ minikube start

### Deployment Steps

 1. Go Web App
This will load the web app Docker image in the cluster.	

<table class="table table-striped table-bordered">
<tbody>
<tr>
	<td><b>Name</b></td>
	<td>go-cache-poc</td>
</tr>
<tr>
	<td><b>Kind</b></td>
	<td>Deployment</td>
</tr>
<tr>
	<td><b>YAML</b></td>
	<td>go-cache-poc-app.yaml</td>
</tr>
</tbody>
</table>

**Command to deploy**

	    $ kubectl apply -f go-cache-poc-app.yaml	

**Verify**

	$ kubectl get deployments
	NAME           READY   UP-TO-DATE   AVAILABLE   AGE
	go-cache-poc   3/3     3            3           14s	
	  
Three pods are running under this deployment.

 - Go Web App Service
 This service will create an external endpoint using a LoadBalancer.

<table class="table table-striped table-bordered">
<tbody>
<tr>
	<td><b>Name</b></td>
	<td>go-cache-poc-service</td>
</tr>
<tr>
	<td><b>Kind</b></td>
	<td>Service</td>
</tr>
<tr>
	<td><b>YAML</b></td>
	<td>go-cache-poc-svc.yaml</td>
</tr>
</tbody>
</table>

 
**Command to deploy**

	$ kubectl apply -f go-cache-poc-svc.yaml
	
**Verify**

	$ kubectl get services
	// Need to write

 - Deploying MongoDB ReplicaSet as a Kubernetes StatefulSet
Kubernetes provides a feature that will allow us to create a stateful application in the cluster. There will be a storage class and services running under the cluster that will allow the databases to connect with services and store records in their persistent database.

 - **MongoDB StorageClass**
 This will create the StorageClass that will be used for the storage
<table class="table table-striped table-bordered">
<tbody>
<tr>
	<td><b>Name</b></td>
	<td>mongodb-storage</td>
</tr>
<tr>
	<td><b>Kind</b></td>
	<td>StorageClass</td>
</tr>
<tr>
	<td><b>YAML</b></td>
	<td>mongodb-storage.yaml</td>
</tr>
</tbody>
</table>


Command to deploy

	$ kubectl apply -f mongodb-app-svc.yaml

 - **MongoDB service**
This will create the StatefulSet app and the Mongo services in the cluster.
<table class="table table-striped table-bordered">
<tbody>
<tr>
	<td><b>Name</b></td>
	<td>mongodb-svc</td>
</tr>
<tr>
	<td><b>Kind</b></td>
	<td>StatefulSet, Service</td>
</tr>
<tr>
	<td><b>YAML</b></td>
	<td>mongodb-app-svc.yaml</td>
</tr>
</tbody>
</table>

Command to deploy

	   $ kubectl apply -f mongodb-app-svc.yaml

#### Define the Administrator
There will be three mongo containers in the cluster. We need to connect to anyone of them to define the administrator.
Command to exec

	 $ kubectl exec -it mongod-app-0 -c mongod-container bash
	-it: mongo contrainer name
	-c:  mongo container name
Bash
$ hostname -f
	     
	 mongod-app-0.mongodb-svc.default.svc.cluster.local
Mongo Shell

	 $ mongo

Type to the following query to generate the replica set

	> rs.initiate({_id: "MainRepSet", version: 1, members: [
		{ _id: 0, host : "mongod-app-0.mongodb-service.default.svc.cluster.local:27017" },
		{ _id: 1, host : "mongod-app-1.mongodb-service.default.svc.cluster.local:27017" },
		{ _id: 2, host : "mongod-app-2.mongodb-service.default.svc.cluster.local:27017" }
	]}); 	  	
	 	 
then verify	
			    	> rs.status();

Now create the Admin user

		> db.getSiblingDB("admin").createUser({
		      user : "admin",
		      pwd  : "admin123",
		      roles: [ { role: "root", db: "admin" } ]
		 });
So, now the MongoDB is complete setup with ReplicaSet and with an Administrator for the database.

 - Deploy Redis in Kubernetes
There will be deployment and service running in the Kubernetes cluster. The connection string will change the redis client for both local and server environments.
	 - **Connection URI**
	<table class="table table-striped table-bordered">
	<tbody>
	<tr>
		<td><b>Local</b></td>
		<td>localhost:6379</td>
	</tr>
	<tr>
		<td><b>Server</b></td>
		<td>redis.default.svc.cluster.local:6379</td>
	</tr>	 
	</tbody>
	</table>
	 - **Redis Deployment**
		<table class="table table-striped table-bordered">
		  <tbody>
			<tr>
				<td><b>Name</b></td>
				<td>redis-app</td>
			</tr>
			<tr>
				<td><b>Kind</b></td>
				<td>Deployment</td>
			</tr>
			<tr>
				<td><b>YAML</b></td>
				<td>redis-deployment.yaml</td>
			</tr>
		  </tbody>
		</table>
	Command to deploy
	 	
	$ kubectl apply -f redis-deployment.yaml
		
   - **Redis Service**
   		<table class="table table-striped table-bordered">
		  <tbody>
			<tr>
				<td><b>Name</b></td>
				<td>redis-service</td>
			</tr>
			<tr>
				<td><b>Kind</b></td>
				<td>Service</td>
			</tr>
			<tr>
				<td><b>YAML</b></td>
				<td>redis-service.yaml</td>
			</tr>
		  </tbody>
		</table>
	   
	     
Command to deploy
	 	
	 $ kubectl apply -f redis-service.yaml

 **7. Deploy Kafka in Kubernetes**
There will be a deployment of ZooKeeper, Kafka Service, and running kafka/zookeeper server script.

 - **Deploy Zookeeper**
There will be deployment and service similar to the other Pods running in the cluster.
	 - **zookeeper-deployment**
	 <table class="table table-striped table-bordered">
		  <tbody>
			<tr>
				<td><b>Name</b></td>
				<td>zookeeper-app</td>
			</tr>
			<tr>
				<td><b>Kind</b></td>
				<td>Deployment</td>
			</tr>
			<tr>
				<td><b>YAML</b></td>
				<td>zookeeper-deployment.yaml</td>
			</tr>
		  </tbody>
		</table>
		 
Command to deploy
	 	
		    $ kubectl apply -f zookeeper-deployment.yaml
		    
 - **zookeeper-service**
 <table class="table table-striped table-bordered">
		  <tbody>
			<tr>
				<td><b>Name</b></td>
				<td>zookeeper-service</td>
			</tr>
			<tr>
				<td><b>Kind</b></td>
				<td>Service</td>
			</tr>
			<tr>
				<td><b>YAML</b></td>
				<td>zookeeper-service.yaml</td>
			</tr>
		  </tbody>
		</table>
		 
Command to deploy
	 	
		    $ kubectl apply -f zookeeper-service.yaml
		    
- **Deploy Kafka**
	 

	- **kafka-service**
	<table class="table table-striped table-bordered">
		  <tbody>
			<tr>
				<td><b>Name</b></td>
				<td>kafka-service</td>
			</tr>
			<tr>
				<td><b>Kind</b></td>
				<td>Service</td>
			</tr>
			<tr>
				<td><b>YAML</b></td>
				<td>kafka-service.yaml</td>
			</tr>
		  </tbody>
		</table>
		   
	        Command to deploy
	 	
				$ kubectl apply -f kafka-service.yaml

	- **kafka-replication-controller**
	<table class="table table-striped table-bordered">
		  <tbody>
			<tr>
				<td><b>Name</b></td>
				<td>kafka-repcon</td>
			</tr>
			<tr>
				<td><b>Kind</b></td>
				<td>Deployment</td>
			</tr>
			<tr>
				<td><b>YAML</b></td>
				<td>kafka-repcon.yaml</td>
			</tr>
		  </tbody>
		</table>
			  
	Command to deploy
	 	
			    $ kubectl apply -f kafka-repcon.yaml
		

    - **Start Zookeeper/Kafka server**
				
	 - **zookeeper server**
		  
		  $ cd kafka/
		  $~/kafka/ bin/zookeeper-server-start.sh config/zookeeper.properties	

	- **Kafka server**
		 
		  $ cd kafka/
		  $~/kafka/ bin/kafka-server-start.sh config/server.properties
				

## Troubleshoot with kubectl

The kubectl is a very handy tool while troubleshooting application into the Kubernetes.

**Few useful commands**
<ol>
	<li> kubectl get pods //List all pods</li>
	<li> kubectl describe pods <pod-name> //Inspect a pod </li>
	<li> kubectl logs <pod-name> //Check the logs for a pod </li>
	<li> kubectl exec -ti <pod-name> --bash //Get inside into the pod shell</li>
</ol>

## Swagger API documentation

The go-swagger toolkit is being integrated for the REST APIs documentation. The API doc can be accessible via http://localhost:5000/docs

https://github.com/go-swagger/go-swagger
