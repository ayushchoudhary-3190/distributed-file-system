🚀 Distributed File System (DFS) — Cloud-Native Storage Platform

A scalable, fault-tolerant distributed file system built using Golang, gRPC, Docker, and Kubernetes. This system enables chunk-based file storage, metadata management, replication, and high availability while being cost-efficient by simulating a real distributed storage cluster using multiple Kubernetes DataNode replicas on a single physical machine.

📌 Key Highlights

Chunk-based distributed storage architecture

Centralized metadata management (MetaServer)

Parallel chunk upload & download pipeline

Automatic replication for fault tolerance

Multi-replica DataNode simulation using Kubernetes (no multiple EC2 instances required)

Containerized deployment using Docker

Horizontal scalability using Kubernetes orchestration

gRPC-based high-performance service communication
🧠 High-Level Architecture (HLD)
Architecture Overview

The system follows a control plane + data plane separation model:

MetaServer → Control Plane (metadata, chunk mapping, node coordination)

DataNodes → Data Plane (actual file chunk storage)

Client → Upload/download orchestration

Kubernetes Layer → Node orchestration, scaling, failover
🧩 Cost-Efficient Distributed Cluster Simulation (Kubernetes Design)

Instead of provisioning multiple cloud VMs, this system uses Kubernetes replica-based orchestration to simulate a real distributed storage cluster.

How It Works

A single physical machine (local system or VM) acts as a Kubernetes Node

Multiple DataNode Pods are deployed using replicas

Each Pod:

Runs an independent DataNode service instance

Has its own isolated process space
Uses dedicated persistent storage volumes

Has a unique internal IP address

This design allows the system to behave like:

3 Separate Storage Servers
      ↓
Single Physical Machine
      ↓
3 Kubernetes DataNode Pods

Benefits

Eliminates cloud infrastructure cost

Preserves real-world distributed behavior

Enables fault simulation and recovery testing

Supports horizontal scaling using replica updates

🔄 Complete Data Flow
📤 File Upload Workflow
Client → MetaServer → Kubernetes Service → DataNodes

Steps:

Client sends upload request to MetaServer

MetaServer assigns:

Chunk IDs

Target DataNode replicas

Client splits file into fixed-size chunks

Chunks are uploaded in parallel to different DataNode Pods

Kubernetes load-balances traffic across DataNode replicas

Replication ensures redundancy across pods

MetaServer updates centralized metadata registry

📥 File Download Workflow
Client → MetaServer → DataNodes → Client

Steps:

Client requests metadata

MetaServer returns chunk locations

Client fetches chunks from nearest available DataNode Pod

Client merges chunks locally

Final file is reconstructed

🐳 Containerization Architecture

All system components are containerized using Docker.

Services Running as Containers

MetaServer

DataNode Service

Client Service

Metadata Database

Advantages

Environment consistency

Faster deployments

Easy versioning

Portable builds

☸ Kubernetes Orchestration Layer

Managed using Kubernetes.

Kubernetes Components Used
📌 Deployments

MetaServer Deployment

DataNode Deployment (replica-based scaling)

Enables:

Zero downtime updates

Automated restarts

Horizontal scaling

📌 Services
Service Type	Usage
ClusterIP	Internal service discovery
LoadBalancer / NodePort	External client access
📌 Persistent Volumes

Each DataNode Pod is attached to its own:

Persistent Volume Claim (PVC)

Dedicated storage directory

This guarantees:

Physical separation of chunk replicas

Crash-safe storage

Stateful workload behavior

📌 Auto Healing

If a DataNode Pod crashes:

Kubernetes automatically recreates it

MetaServer reassigns chunk replicas

System remains operational

⚙ Core System Design Principles
✅ Stateless Control Services

MetaServer remains stateless

All metadata persisted in database

Enables seamless scaling and failover

✅ Fault Tolerance

Multi-replica chunk storage

Pod failure recovery

Kubernetes self-healing

✅ Horizontal Scalability

Scaling DataNodes:

kubectl scale deployment datanode --replicas=5


No downtime required.

✅ High Throughput Architecture

Parallel chunk uploads

gRPC streaming

Distributed I/O pipeline

📈 Why This Architecture Matters

This project demonstrates real-world engineering skills:

Distributed systems design

Kubernetes orchestration

Cloud-native architecture

Storage system internals

Cost-optimized infrastructure simulation

DevOps + backend integration
