FROM redis:7.2.5-alpine

CMD sleep 10; \
    redis-cli \
        --cluster create master0:6079 master1:6179 master2:6279 \
        --cluster-replicas 0 \
        --cluster-yes; \
    redis-cli --cluster add-node replica0:6080 master0:6079 --cluster-slave; \
    redis-cli --cluster add-node replica1:6180 master1:6179 --cluster-slave; \
    redis-cli --cluster add-node replica2:6280 master2:6279 --cluster-slave