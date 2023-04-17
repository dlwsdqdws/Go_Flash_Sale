# Go High Concurrency

## Project Structure
1. model
2. repositories
3. services
4. controllers
5. views

## Data Flow
<p align="center"><img src="static/img/data_flow.png" alt="data_flow" width="500" /></p> 

## Optimization
1. Bidirectional Encryption
2. Distributed System
<br>Load balancing ：Consistent Hashing
3.  RabbitMq

## Test
### Servers
| Server Name      |  Intranet IP  |      Public IP |
|:-----------------|:-------------:|---------------:|
| Validate1 + SLB  | 172.26.194.42 |   47.250.49.25 |
| Validate2 + SLB  | 172.26.194.41 | 47.250.147.115 |
| Stress Test      | 172.17.169.89 | 47.250.131.235 |
| Quantity Control | 172.26.194.43 | 47.250.144.234 |
| RabbitMQ Simple  | 172.17.169.88 |  47.254.251.25 |
### WRK


## Environment
Go 1.18.1 arm64
<br>IRIS 11.1.1
<br>Mysql 8.0.31
<br>RabbitMQ 3.11.11