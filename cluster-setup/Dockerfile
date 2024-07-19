FROM redis:7.2.5-alpine

CMD sleep 10; \
    redis-cli \
        --cluster create 127.0.0.1:6079 127.0.0.1:6179 127.0.0.1:6279 \
        --cluster-replicas 0 \
        --cluster-yes; \
    redis-cli --cluster add-node 127.0.0.1:6080 127.0.0.1:6079 --cluster-slave; \
    redis-cli --cluster add-node 127.0.0.1:6180 127.0.0.1:6179 --cluster-slave; \
    redis-cli --cluster add-node 127.0.0.1:6280 127.0.0.1:6279 --cluster-slave