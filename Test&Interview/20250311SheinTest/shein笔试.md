# shein笔试

## 问题总结

1.高并发场景下，对于数据一致性问题，如何解决（加锁的具体处理）、比如银行存取、网站评论和动态读取。

2.goroutine的代码读，自己不熟悉。

3.docker相关的操作自己也不熟悉，基本的启动代码什么的得记住。

4.还有k8s。

## 处理

### 高并发场景

> 在银行存取款场景中，确保高并发情况下的数据一致性至关重要。以下是一些关键策略和技术，可以帮助你在高并发环境中实现这一目标：
>
> ### 1. 使用数据库事务
>
> 数据库事务是保证数据一致性的基本手段。通过将存取款操作封装在事务中，可以确保这些操作的原子性，即要么全部成功，要么全部失败。
>
> ```go
> func transferMoney(db *sql.DB, fromAccount, toAccount string, amount float64) error {
>     tx, err := db.Begin()
>     if err != nil {
>         return err
>     }
>     defer tx.Rollback()
> 
>     // 锁定账户以防止并发问题
>     _, err = tx.Exec("SELECT balance FROM accounts WHERE account_id = ? FOR UPDATE", fromAccount)
>     if err != nil {
>         return err
>     }
>     _, err = tx.Exec("SELECT balance FROM accounts WHERE account_id = ? FOR UPDATE", toAccount)
>     if err != nil {
>         return err
>     }
> 
>     // 执行转账操作
>     _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE account_id = ?", amount, fromAccount)
>     if err != nil {
>         return err
>     }
>     _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE account_id = ?", amount, toAccount)
>     if err != nil {
>         return err
>     }
> 
>     // 提交事务
>     return tx.Commit()
> }
> ```
>
> ### 2. 使用乐观锁或悲观锁
>
> - **悲观锁**：如上例所示，使用 `FOR UPDATE` 语句对账户进行锁定，确保在事务提交之前，其他事务不能修改这些记录。
> - **乐观锁**：通过版本号或时间戳来实现。在更新数据时，检查版本号是否一致，如果不一致则重试。
>
> ```go
> func updateBalanceWithOptimisticLock(db *sql.DB, accountID string, amount float64, version int) error {
>     for {
>         // 查询当前余额和版本号
>         var balance float64
>         var currentVersion int
>         err := db.QueryRow("SELECT balance, version FROM accounts WHERE account_id = ?", accountID).Scan(&balance, &currentVersion)
>         if err != nil {
>             return err
>         }
> 
>         // 更新余额和版本号
>         result, err := db.Exec("UPDATE accounts SET balance = ?, version = version + 1 WHERE account_id = ? AND version = ?", balance+amount, accountID, currentVersion)
>         if err != nil {
>             return err
>         }
> 
>         // 检查受影响的行数
>         rowsAffected, err := result.RowsAffected()
>         if err != nil {
>             return err
>         }
> 
>         if rowsAffected > 0 {
>             return nil // 更新成功
>         }
> 
>         // 如果受影响的行数为0，说明版本号不匹配，重试
>     }
> }
> ```
>
> ### 3. 使用分布式锁
>
> 在高并发分布式系统中，可以使用分布式锁（如Redis的 `SETNX` 命令或Zookeeper）来确保同一时间只有一个节点可以修改某个账户的数据。
>
> ### 4. 消息队列
>
> 使用消息队列（如Kafka、RabbitMQ）来解耦和削峰。将存取款请求放入消息队列中，然后由消费者逐个处理这些请求，可以有效减轻数据库的压力，并确保操作的顺序性。
>
> ```go
> func enqueueTransfer(db *sql.DB, fromAccount, toAccount string, amount float64) error {
>     // 将转账请求放入消息队列
>     msg := TransferMessage{
>         FromAccount: fromAccount,
>         ToAccount:   toAccount,
>         Amount:      amount,
>     }
>     return messageQueue.Publish("transfer_queue", msg)
> }
> 
> func processTransfers() {
>     for msg := range messageQueue.Consume("transfer_queue") {
>         transferMoney(db, msg.FromAccount, msg.ToAccount, msg.Amount)
>     }
> }
> ```
>
> ### 5. 数据库分片和读写分离
>
> 通过数据库分片和读写分离，可以将负载分散到多个数据库实例上，提高系统的吞吐量和可用性。
>
> ### 6. 监控和日志
>
> 实施全面的监控和日志记录，及时发现和解决性能瓶颈和数据一致性问题。
>
> ### 7. 测试和压测
>
> 进行充分的单元测试、集成测试和压力测试，确保系统在高并发情况下的稳定性和数据一致性。
>
> 通过以上策略和技术，可以在高并发场景下尽量保证银行存取款操作的数据一致性。
>
> - **悲观锁**：适用于写操作频繁、冲突概率高的场景。可以确保数据一致性，但可能会影响性能。
> - **乐观锁**：适用于读操作频繁、冲突概率低的场景。可以提高并发性能，但需要处理重试逻辑。

网站场景下：

> 在网站评论和动态读取的场景下，既要满足高并发，又要保证数据一致性，可以采用以下策略和技术：
>
> ### 1. 数据库设计与优化
>
> #### a. 合理的数据库表结构
> • **分表分库**：根据业务需求，将评论和动态数据按照一定规则（如用户ID、时间等）进行分表或分库，分散数据库压力。
> • **索引优化**：为常用的查询字段（如动态ID、用户ID、时间戳等）建立合适的索引，提高查询效率。
>
> #### b. 使用缓存
> • **Redis**：将热点数据（如最新评论、热门动态）缓存到Redis中，减少数据库读取压力，提高读取速度。
> • **缓存更新策略**：采用合适的缓存更新策略（如Cache-Aside、Read-Through、Write-Through等），确保缓存与数据库数据的一致性。
>
> ### 2. 并发控制
>
> #### a. 乐观锁
> • **版本号机制**：在评论和动态表中增加一个版本号字段，每次更新数据时，检查版本号是否一致，若不一致则重试。
>
> #### b. 悲观锁
> • **行级锁**：在更新或删除评论时，使用 `SELECT ... FOR UPDATE` 对相关行加锁，防止并发修改导致的数据不一致。
> • **适用场景**：适用于写操作频繁且对数据一致性要求极高的场景。
>
> ### 3. 消息队列
>
> • **异步处理**：将评论和动态的写入操作放入消息队列（如Kafka、RabbitMQ），由消费者异步处理，减轻数据库压力，提高系统吞吐量。
> • **顺序性保证**：确保同一动态下的评论按顺序处理，避免评论顺序错乱。
>
> ### 4. 分布式锁
>
> • **Redis分布式锁**：在高并发写入场景下，使用Redis的 `SETNX` 命令或Redlock算法实现分布式锁，确保同一时间只有一个节点能修改特定数据。
> • **适用场景**：适用于需要在分布式环境下保证数据一致性的场景。
>
> ### 5. 数据库事务
>
> • **事务隔离级别**：根据业务需求选择合适的事务隔离级别（如Read Committed、Repeatable Read等），平衡数据一致性与并发性能。
>
> ### 6. 读写分离
>
> • **主从复制**：配置数据库主从复制，主库负责写操作，从库负责读操作，分散读写压力，提高系统性能。
> • **负载均衡**：在应用层通过负载均衡策略（如轮询、最小连接数等）将读请求分发到不同的从库。
>
> ### 7. 限流与降级
>
> • **限流**：通过令牌桶、漏桶等算法限制请求速率，防止系统过载。
> • **降级**：在高并发情况下，暂时关闭部分非核心功能，确保核心功能的稳定运行。
>
> ### 8. 监控与日志
>
> • **实时监控**：使用Prometheus、Grafana等工具实时监控系统状态，及时发现和处理性能瓶颈。
> • **日志记录**：详细记录操作日志，便于排查问题和审计。
>
> ### 9. 数据一致性保障
>
> #### a. 最终一致性
> • **异步补偿**：对于一些非关键操作，允许短暂的数据不一致，通过后台任务进行补偿，最终达到一致状态。
>
> #### b. 分布式事务
> • **两阶段提交（2PC）**：在涉及多个数据库或服务的操作中，使用两阶段提交协议确保所有节点的数据一致性。
> • **Saga模式**：通过一系列本地事务和补偿操作，确保分布式系统中的数据最终一致。
>
> ### 总结
>
> 在网站评论和动态读取的高并发场景下，保证数据一致性需要综合运用数据库优化、缓存策略、并发控制、消息队列、分布式锁等多种技术手段。通过合理的设计与实现，可以在确保数据一致性的同时，提升系统的性能和用户体验。
>

## docker

---

### 2. **启动 Docker 服务**
在 Linux 系统上，启动 Docker 服务：
```bash
sudo systemctl start docker
```
设置开机自启：
```bash
sudo systemctl enable docker
```

---

### 3. **检查 Docker 版本**
查看 Docker 的版本信息：
```bash
docker --version
```

---

### 4. **镜像操作**
#### 4.1 拉取镜像
从 Docker Hub 拉取镜像：
```bash
docker pull <镜像名称>:<标签>
```
例如：
```bash
docker pull ubuntu:20.04
```

#### 4.2 查看本地镜像
列出本地所有的镜像：
```bash
docker images
```

#### 4.3 删除镜像
删除本地镜像：
```bash
docker rmi <镜像ID或镜像名称>
```
例如：
```bash
docker rmi ubuntu:20.04
```

#### 4.4 搜索镜像
在 Docker Hub 上搜索镜像：
```bash
docker search <镜像名称>
```
例如：
```bash
docker search nginx
```

---

### 5. **容器操作**
#### 5.1 创建并运行容器
创建并启动一个容器：
```bash
docker run [选项] <镜像名称>:<标签> [命令]
```
常用选项：
• `-d`：后台运行容器
• `-p`：端口映射（主机端口:容器端口）
• `-v`：挂载主机目录到容器
• `--name`：指定容器名称

示例：
```bash
docker run -d --name my-nginx -p 8080:80 nginx
```

#### 5.2 查看运行中的容器
列出所有运行中的容器：
```bash
docker ps
```
列出所有容器（包括停止的）：
```bash
docker ps -a
```

#### 5.3 停止容器
停止一个容器：
```bash
docker stop <容器ID或容器名称>
```
例如：
```bash
docker stop my-nginx
```

#### 5.4 启动容器
启动一个已停止的容器：
```bash
docker start <容器ID或容器名称>
```

#### 5.5 删除容器
删除一个容器：
```bash
docker rm <容器ID或容器名称>
```
删除所有已停止的容器：
```bash
docker container prune
```

#### 5.6 进入容器
进入容器的终端：
```bash
docker exec -it <容器ID或容器名称> /bin/bash
```
例如：
```bash
docker exec -it my-nginx /bin/bash
```

#### 5.7 查看容器日志
查看容器的日志：
```bash
docker logs <容器ID或容器名称>
```
例如：
```bash
docker logs my-nginx
```

---

### 6. **Docker 网络**
#### 6.1 查看网络
列出所有 Docker 网络：
```bash
docker network ls
```

#### 6.2 创建网络
创建一个新的网络：
```bash
docker network create <网络名称>
```
例如：
```bash
docker network create my-network
```

#### 6.3 删除网络
删除一个网络：
```bash
docker network rm <网络名称>
```

---

### 7. **Docker 数据卷**
#### 7.1 创建数据卷
创建一个数据卷：
```bash
docker volume create <卷名称>
```
例如：
```bash
docker volume create my-volume
```

#### 7.2 查看数据卷
列出所有数据卷：
```bash
docker volume ls
```

#### 7.3 删除数据卷
删除一个数据卷：
```bash
docker volume rm <卷名称>
```

#### 8.2 使用 Docker Compose
创建一个 `docker-compose.yml` 文件，然后运行：
```bash
docker-compose up -d
```
停止服务：
```bash
docker-compose down
```

---

### 9. **清理未使用的资源**
清理未使用的镜像、容器、网络和数据卷：
```bash
docker system prune
```
清理所有未使用的资源（包括数据卷）：
```bash
docker system prune -a --volumes
```

---

### 10. **其他常用命令**
• 查看 Docker 系统信息：
  ```bash
  docker info
  ```
• 查看 Docker 详细信息：
  ```bash
  docker inspect <容器ID或镜像名称>
  ```

---

## k8s

Kubernetes（简称K8s）是一个开源的容器编排平台，用于自动化部署、扩展和管理容器化应用程序。以下是一些Kubernetes的基础操作，包括集群创建、资源管理、应用部署和监控等。

---

### 1. **Kubernetes集群基础**
#### 1.1 安装Kubernetes
- **Minikube（单节点集群）**：适合本地开发和测试。
  ```bash
  # 安装Minikube
  curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
  sudo install minikube-linux-amd64 /usr/local/bin/minikube
  
  # 启动Minikube
  minikube start
  ```

- **Kubeadm（生产级集群）**：
  1. 安装Docker和Kubernetes组件。
     ```bash
     sudo apt-get update
     sudo apt-get install -y apt-transport-https curl
     curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
     echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
     sudo apt-get update
     sudo apt-get install -y kubelet kubeadm kubectl
     sudo apt-mark hold kubelet kubeadm kubectl
     ```
  2. 初始化Kubernetes主节点。
     ```bash
     sudo kubeadm init --pod-network-cidr=10.244.0.0/16
     ```
  3. 设置kubectl配置。
     ```bash
     mkdir -p $HOME/.kube
     sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
     sudo chown $(id -u):$(id -g) $HOME/.kube/config
     ```
  4. 部署Pod网络插件（如Flannel）。
     ```bash
     kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
     ```

#### 1.2 集群状态检查
- 检查集群状态：
  ```bash
  kubectl cluster-info
  ```
- 查看节点状态：
  ```bash
  kubectl get nodes
  ```

---

### 2. **Kubernetes资源管理**
#### 2.1 常用资源对象
- **Pod**：Kubernetes的最小部署单元。
- **Deployment**：用于管理无状态应用，支持自动扩缩容和滚动更新。
- **Service**：用于定义Pod的访问方式，如ClusterIP、NodePort、LoadBalancer。
- **Namespace**：用于资源隔离，一个集群可以有多个Namespace。

#### 2.2 创建和管理资源
- **创建Deployment**：
  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: my-app
  spec:
    replicas: 3
    selector:
      matchLabels:
        app: my-app
    template:
      metadata:
        labels:
          app: my-app
      spec:
        containers:
        - name: my-container
          image: nginx:latest
  ```
  应用配置文件：
  ```bash
  kubectl apply -f my-deployment.yaml
  ```

- **查看Deployment状态**：
  ```bash
  kubectl get deployments
  ```

- **更新Deployment（滚动更新）**：
  ```bash
  kubectl set image deployment/my-app my-container=nginx:1.21
  ```

- **创建Service**：
  ```yaml
  apiVersion: v1
  kind: Service
  metadata:
    name: my-service
  spec:
    selector:
      app: my-app
    ports:
      - protocol: TCP
        port: 80
        targetPort: 80
    type: ClusterIP
  ```
  应用配置文件：
  ```bash
  kubectl apply -f my-service.yaml
  ```

- **查看Service状态**：
  ```bash
  kubectl get services
  ```

- **删除资源**：
  ```bash
  kubectl delete deployment my-app
  kubectl delete service my-service
  ```

---

### 3. **应用部署**
#### 3.1 部署应用
- 使用`kubectl create`命令快速部署：
  ```bash
  kubectl create deployment my-app --image=nginx
  ```

- 使用`kubectl expose`创建Service：
  ```bash
  kubectl expose deployment my-app --port=80 --type=NodePort
  ```

#### 3.2 访问应用
- 查看Service的NodePort：
  ```bash
  kubectl get svc my-app
  ```
  使用`<NodeIP>:<NodePort>`访问应用。

---

### 4. **监控和日志**
#### 4.1 查看Pod日志
- 查看单个Pod日志：
  ```bash
  kubectl logs <pod-name>
  ```

- 查看多容器Pod日志：
  ```bash
  kubectl logs <pod-name> -c <container-name>
  ```

#### 4.2 监控集群状态
- 查看Pod状态：
  ```bash
  kubectl get pods -o wide
  ```

- 查看资源使用情况：
  ```bash
  kubectl top pod
  ```

---

### 5. **其他常用命令**
- **进入Pod内部**：
  
  ```bash
  kubectl exec -it <pod-name> -- /bin/sh
  ```
  
- **查看资源详细信息**：
  
  ```bash
  kubectl describe pod <pod-name>
  ```
  
- **查看API Server支持的资源类型**：
  
  ```bash
  kubectl api-resources
  ```

---

### 6. **清理集群**
- **清理Minikube**：
  ```bash
  minikube stop
  minikube delete
  ```

- **清理Kubeadm集群**：
  ```bash
  kubeadm reset
  ```

---

