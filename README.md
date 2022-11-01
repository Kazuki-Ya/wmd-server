# wmd-server
wmd-server is my mini microservice using gRPC. This microservice is composed of web-service, inference-service (machine learning service) and log-service.

## Structure
![alt text](/img.png)

## Web-server
Web-server is implemented in Go. This server recieves a client's request and sends RPC to inference or log service.

## Log-server
Log-server is implemented in Go, Using Raft algorithm and Serf Discovery service. Because of Raft , this log-server is fault-torelance and managing a replicated log.

## Inference-server
Inference-server is implemented in Python. When this server recieves RPC (InferenceCall  api) , it output label from BERT model and sends back to web-server.

## Port
In order to deploy locally or to Cloud, port number is 

|         |  port number |
| :---:   |  :---:       |
| Log RPC |  8400        |
| Serf port | 8401       |
| Web port  | 8080       |
| Inference port | 8403  |


