FROM redis:7.2.5-alpine

CMD sleep 10; \
    redis-cli -c -h master0 -p 6079 FLUSHALL; \
    redis-cli -c -h master0 -p 6079 CLUSTER RESET; \
    redis-cli -c -h master1 -p 6179 FLUSHALL; \
    redis-cli -c -h master1 -p 6179 CLUSTER RESET; \
    redis-cli -c -h master2 -p 6279 FLUSHALL; \
    redis-cli -c -h master2 -p 6279 CLUSTER RESET; \
    redis-cli -c -h replica0 -p 6080 FLUSHALL; \
    redis-cli -c -h replica1 -p 6180 FLUSHALL; \
    redis-cli -c -h replica2 -p 6280 FLUSHALL; \
    redis-cli \
        --cluster create master0:6079 master1:6179 master2:6279 \
        --cluster-replicas 0 \
        --cluster-yes; \
    redis-cli --cluster add-node replica0:6080 master0:6079 --cluster-slave; \
    redis-cli --cluster add-node replica1:6180 master1:6179 --cluster-slave; \
    redis-cli --cluster add-node replica2:6280 master2:6279 --cluster-slave