# dorisctl

A tool to deploy apache [doris](https://github.com/apache/doris/) cluster.

## Installation

```sh
go install github.com/lobshunter/dorisctl/cmd/dorisctl
```

## Usage

```sh
# NOTE: `ssh root@<host> -i <ssh_private_key_path>` must work for every hosts in deployment yaml file

# deploy doris
$ dorisctl deploy examples/minimal-amd64.yaml

# start cluster
$ dorisctl start

# list clusters
$ dorisctl list
+--------------+-----------+---------------+
| CLUSTER NAME | FEMASTER  | FE QUERT PORT |
+--------------+-----------+---------------+
|   default    | 127.0.0.1 |     9030      |
+--------------+-----------+---------------+

# show cluster status
$ dorisctl status
Check Cluster Status  Done

Frontends:
+-----------+-----------+----------+-------+---------------------------+
|   HOST    | QUERYPORT | ISMASTER | ALIVE |          VERSION          |
+-----------+-----------+----------+-------+---------------------------+
| 127.0.0.1 |   9030    |   true   | true  | doris-2.0-beta-afe6bb9638 |
+-----------+-----------+----------+-------+---------------------------+

Backends:
+-----------+-------+---------------+---------------+---------+---------------------------+
|   HOST    | ALIVE | AVAILCAPACITY | TOTALCAPACITY | USEDPCT |          VERSION          |
+-----------+-------+---------------+---------------+---------+---------------------------+
| 127.0.0.1 | true  | 77.527 GB     | 96.727 GB     | 19.85 % | doris-2.0-beta-afe6bb9638 |
+-----------+-------+---------------+---------------+---------+---------------------------+

# stop cluster
$ dorisctl stop

# destroy cluster
$ dorisctl destory

# take over a manually deployed cluster so it can be managed by dorisctl
$ dorisctl takeover --cluster-name yelo --fe-hosts 172.30.0.9,172.30.0.10 --fe-master 172.30.0.9 --be-hosts 172.30.0.5,172.30.0.6,172.30.0.12 --fe-deploy-dir /doris/fe  --be-deploy-dir /doris/be

# handover managed cluster (remove manifest from dorisctl, without deleting the cluster)
$ dorisctl handover --cluster-name yelo
```
