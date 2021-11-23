NODE_DIR=~/Projects/ln-test-net

go build
cp go-hosted-channels $NODE_DIR/l1-regtest/plugins
cp go-hosted-channels $NODE_DIR/l2-regtest/plugins