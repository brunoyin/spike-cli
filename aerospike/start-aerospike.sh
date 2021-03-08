cd [spike-cli loacal repo]
docker run -d \
    -v $PWD/aerospike/data:/opt/aerospike/data \
    -v $PWD/aerospike//conf:/opt/aerospike/etc \
    --name aerospike -p 3000:3000 -p 3001:3001 -p 3002:3002 -p 3003:3003 \
    aerospike /usr/bin/asd --foreground --config-file /opt/aerospike/etc/aerospike.conf