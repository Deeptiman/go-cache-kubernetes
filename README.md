# go-cache-kubernetes
<p>     
	<img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/Deeptiman/go-cache-kubernetes">
	<img alt="GitHub language count" src="https://img.shields.io/github/languages/count/Deeptiman/go-cache-kubernetes"> 
	<img alt="GitHub top language" src="https://img.shields.io/github/languages/top/Deeptiman/go-cache-kubernetes"> 
</p>
The web application is a Data Caching service designed and implemented using microservices architecture. The cloud deployment environments are used Kubernetes, Docker, and written in Go programming language. The application also uses a MongoDB as NoSQL database with Redis in-memory database for the caching services.

## Features

 - **MongoDB** is implemented to perform several database operations. The installation can be done using the go dependency module.
		
		go get go.mongodb.org/mongo-driver/mongo
	  https://github.com/mongodb/mongo-go-driver

 - **Redis Cache** is implemented to integrate the data caching in the application. So, the <b>go-redis</b> will cache the second GET request while reading the user details.
 
		go get github.com/go-redis/redis/v8
	  https://github.com/go-redis/redis

 - **Kafka Message Broker**: The <b>confluent-kafka-go</b> is used as a Go client library for Kafka message broker. The library will provide **Producer** and **Consumer** architecture to stream messages to the user for a subscribed topic. So, there will be two REST APIs that the user can use for Producing the messages reading from MongoDB and Consuming or Reading messages from the message broker.

		go get github.com/confluentinc/confluent-kafka-go/kafka
	https://github.com/confluentinc/confluent-kafka-go

	Note: It's recommended to install **confluent-kafka-go v1.4.0**, as the **librdkafka** will come with the bundle and no need to install separately.

## Kubernetes tools
Kubernetes provides several tools that can be useful to setup Kubernetes in the local environment.

 - **minikube** tool will run a single-node Kubernetes cluster running inside a Virtual Machine. Virtualization has to be supported in the computer and Hypervisor needed to be enabled.
	    
	 **Installation**
	    follows with the Hypervisor installation and [Hyperkit](https://minikube.sigs.k8s.io/docs/drivers/hyperkit/) is the recommended virtualization toolkit.   

		$ sudo install minikube
		
		$ minikube start
	https://kubernetes.io/docs/setup/learning-environment/minikube/
	 
	
		

 - **kubectl** command-line tool will work to manage a Kubernetes cluster. The tool will be used to deploy, create, analyze, inspect pods that are running under a Kubernetes cluster.


	**Installation**
	
	```curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" ```

	https://kubernetes.io/docs/tasks/tools/install-kubectl/


## Build the Docker images

The application uses Docker for container-based development. The docker image gets published as a public repository at Docker Hub.

 - **Build the image**

		$ docker build -t go-cache-kubernetes-v1 .

 - **Tag the image**

		$ docker tag go-cache-poc deeptiman1991/go-cache-kubernetes-v1:1.0.0

 - **Login to docker hub**

		$ docker login

		Type Username and Password to complete the authentication

 - **Push the image to docker hub**

		$ docker push deeptiman1991/go-cache-kubernetes-v1:1.0.0

## Kubernetes Deployment
There will be several deployments, services that need to be running in the cluster as a Pod. The creation of a Pod requires a YAML file that will specify the kind, spec, containerPort, metadata, volume, and more. So, these parameters will be used to provide resources to the Kubernetes cluster.

**Start minikube** to begin the deployment process start the minikube 
	
	$ minikube start

<h3>Kubernetes Secret Management</h3> 
The application will be using few MongoDB credentials for database connection. So the username and password will be secure using Secret Management via Environment Variables.
<br><br>
 <p> <b> Create Secret literals using kubectl </b> </p>
    <p><code>$ kubectl create secret generic mongosecret --from-literal='username=admin' --from-literal='password=admin123'</code></p>	
    
 <p> <b> Implement the Secret literals in the pod deployment </b> </p>
		
	    spec:
	      containers:
	      - name: go-cache-kubernetes-container-poc
		image: deeptiman1991/go-cache-kubernetes-v1:1.0.0
		env:
		- name: SECRET_USERNAME
		  valueFrom:
		    secretKeyRef:
		      name: mongosecret
		      key: username
		- name: SECRET_PASSWORD
		  valueFrom:
		    secretKeyRef:
		      name: mongosecret
		      key: password		
  <table class="table table-striped table-bordered">
	<tbody>
	<tr>
		<td><b>YAML</b></td>
		<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/go_cache_app/go-cache-poc-app.yaml#L19" target="_blank">go-cache-poc-app.yaml</a></td>
	</tr>
	</tbody>
	</table>
</p>	
<p>So, now SECRET_USERNAME & SECRET_PASSWORD environment variables can be used to connect to the MongoDB database from the application. </p>
    
<h3>Deploy PersistentVolumeClaim</h3> 

This will allocate a volume of 1GB storage in the cluster
<br>
	<table class="table table-striped table-bordered">
	<tbody>
	<tr>
		<td><b>Name</b></td>
		<td>go-cache-poc-pvc</td>
	</tr>
	<tr>
		<td><b>Kind</b></td>
		<td>PersistentVolumeClaim</td>
	</tr>
	<tr>
		<td><b>YAML</b></td>
		<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/go_cache_app/go-cache-poc-pvc.yaml" target="_blank">go-cache-poc-pvc.yaml</a></td>
	</tr>
	</tbody>
	</table>

    $ kubectl apply -f go-cache-poc-pvc.yaml

<h3>Deploy Go Web App</h3> 

This will load the web app Docker image in the cluster.	
<br>
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
		<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/go_cache_app/go-cache-poc-app.yaml" target="_blank">go-cache-poc-app.yaml</a></td>
	</tr>
	</tbody>
	</table>

    $ kubectl apply -f go-cache-poc-app.yaml	

**Verify**

	$ kubectl get deployments
	NAME           		      READY   UP-TO-DATE   AVAILABLE   AGE
	go-cache-kubernetes-app-poc   1/1     1            1           14s	
	  
	There is only one pod is running under this deployment.

<h3>Deploy Go Web App Service</h3>

This service will create an external endpoint using a LoadBalancer.
<br>
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
		<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/go_cache_app/go-cache-poc-svc.yaml" target="_blank">go-cache-poc-svc.yaml</a></td>
	</tr>
	</tbody>
	</table>


	$ kubectl apply -f go-cache-poc-svc.yaml
	
**Verify**

	$ kubectl get services

 <h3>Deploying MongoDB ReplicaSet as a Kubernetes StatefulSet</h3>


<p>Kubernetes provides a feature that will allow us to create a stateful application in the cluster. There will be a storage class and services running under the cluster that will allow the databases to connect with services and store records in their persistent database.</p>

 - **MongoDB service** will create the Mongo services in the cluster.
	<table class="table table-striped table-bordered">
	<tbody>
	<tr>
		<td><b>Name</b></td>
		<td>mongodb-service</td>
	</tr>
	<tr>
		<td><b>Kind</b></td>
		<td>Service</td>
	</tr>
	<tr>
		<td><b>YAML</b></td>
		<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/mongodb/mongodb-service.yaml" target="_blank">mongodb-service.yaml</a></td>
	</tr>
	</tbody>
	</table>

	   $ kubectl apply -f mongodb-service.yaml
	 
- **MongoDB StatefulSet** will create the StatefulSet app in the cluster.
	<table class="table table-striped table-bordered">
	<tbody>
	<tr>
		<td><b>Name</b></td>
		<td>mongod</td>
	</tr>
	<tr>
		<td><b>Kind</b></td>
		<td>StatefulSet</td>
	</tr>
	<tr>
		<td><b>YAML</b></td>
		<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/mongodb/mongodb-stateful.yaml" target="_blank">mongodb-stateful.yaml</a></td>
	</tr>
	</tbody>
	</table>

	   $ kubectl apply -f mongodb-stateful.yaml

#### Define the Administrator
There will be three mongo containers in the cluster. We need to connect to anyone of them to define the administrator.
Command to exec

	 $ kubectl exec -it mongod-app-0 -c mongod-container-app bash
	-it: mongo app name
	-c:  mongo container name
Bash

  	$ hostname -f
	     
	 mongod-app-0.mongodb-service.default.svc.cluster.local

Mongo Shell

	 $ mongo

Type to the following query to generate the replica set

	> rs.initiate({_id: "MainRepSet", version: 1, members: [
		{ _id: 0, host : "mongod-app-0.mongodb-service.default.svc.cluster.local:27017" }
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

<h3>Deploy Redis in Kubernetes</h3>
There will be several steps to follow for deploying Redis into the Kubernete cluster. 
<br><br>
	<p><b>Download Docker images for Redis</b></p>
	<p> 
   	   <code>$ docker run -p 6379:6379 redislabs/redismod</code>	
	</p>	
	 <p><b>Redis Deployment</b></p>
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
				<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/redis/redis-deployment.yaml" target="_blank">redis-deployment.yaml</a></td>
			</tr>
		  </tbody>
		</table>
	 	<p><code>$ kubectl apply -f redis-deployment.yaml</code></p>
	
   <p><b>Redis Service</b></p>
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
				<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/redis/redis-service.yaml" target="_blank">redis-service.yaml</a></td>
			</tr>
		  </tbody>
		</table>
	<p><code>$ kubectl apply -f redis-service.yaml</code></p>	
	<p><b>Deploy the redismod image</b></p>
	<p>
	  <code>$ kubectl run redismod --image=redislabs/redismod --port=6379</code>
	</p>
	<p><b>Expose the deployment</b></p>
	<p>
	  <code>$ kubectl expose deployment redismod --type=NodePort</code>
	</p>
	<p><b>Now, check for the Redis Connection</b></p>
	<p>
	  <code>$ redis-cli -u $(minikube service --format "redis://{{.IP}}:{{.Port}}" --url redismod)</code>
	</p>
	 <code>You must be getting an ip address with a port that can be used as a connection string for Redis</code>

<h3>Deploy Kafka in Kubernetes</h3>
There will be a deployment of ZooKeeper, Kafka Service, and running kafka/zookeeper server script. Please install <a href="https://kafka.apache.org/downloads">Apache Kafka</a> in your local machine and gcloud.
<br>
<h4>Zookeeper</h4>
There will be deployment and service similar to the other Pods running in the cluster.
	 <h5>zookeeper-deployment</h5>
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
				<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/kafka/zookeeper-deployment.yaml" target="_blank">zookeeper-deployment.yaml</a></td>
			</tr>
		  </tbody>
		</table>
	 	
	$ kubectl apply -f zookeeper-deployment.yaml
		    
 <h5>zookeeper-service</h5>
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
				<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/kafka/zookeeper-service.yaml" target="_blank">zookeeper-service.yaml</a></td>
			</tr>
		  </tbody>
		</table>
	 	
    $ kubectl apply -f zookeeper-service.yaml
		    
<h4>Kafka</h4>
<h5>kafka-service</h5>
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
				<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/kafka/kafka-service.yaml" target="_blank">kafka-service.yaml</a></td>
			</tr>
		  </tbody>
		</table>
		   
	         
	 	
 	$ kubectl apply -f kafka-service.yaml

<h5>kafka-replication-controller</h5>
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
				<td><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/deploy_kubernetes/kafka/kafka-repcon.yaml" target="_blank">kafka-repcon.yaml</a></td>
			</tr>
		  </tbody>
		</table>

	 	
	$ kubectl apply -f kafka-repcon.yaml
		

<h3>Start zookeeper/kafka server</h3>
   <h4>zookeeper server</h4>
		  
	$ cd kafka/
	$~/kafka/ bin/zookeeper-server-start.sh config/zookeeper.properties	

   <h4>kafka server</h4>
		 
	$ cd kafka/
	$~/kafka/ bin/kafka-server-start.sh config/server.properties
				

## Troubleshoot with kubectl

The kubectl is a very handy tool while troubleshooting application into the Kubernetes.

**Few useful commands**
<ol>
	<li> kubectl get pods </li>
	<li> kubectl describe pods <pod-name>  </li>
	<li> kubectl logs <pod-name>  </li>
	<li> kubectl exec -ti <pod-name> --bash </li>
</ol>

## Swagger API documentation

The <b>go-swagger</b> toolkit is integrated for the REST APIs documentation. The API doc can be accessible via http://localhost:5000/docs

<p><a href="https://github.com/go-swagger/go-swagger" target="_blank">https://github.com/go-swagger/go-swagger</a></p>
<p><a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/swagger.yaml" target="_blank">swagger.yaml</a></p>

## More Info
<ul>
	<li> <a href="https://kubernetes.io/docs/setup/">Getting started Kubernetes</a> </li>
	<li> <a href="https://kubernetes.io/docs/tutorials/hello-minikube/">Hello Minikube</a></li>
	<li> <a href="https://docs.docker.com/get-started/">Docker Documentation</a></li>
	<li> <a href="https://redocly.github.io/redoc/">Swagger ReDoc</a></li>
</ul>

<h2>License</h2>
<p>This project is licensed under the <a href="https://github.com/Deeptiman/go-cache-kubernetes/blob/master/LICENSE">MIT License</a></p>
