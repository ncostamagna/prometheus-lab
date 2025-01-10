curl -LO https://github.com/prometheus/prometheus/releases/download/v3.1.0/prometheus-3.1.0.linux-amd64.tar.gz
tar -xzf prometheus-3.1.0.linux-amd64.tar.gz
mv prometheus-3.1.0.linux-amd64 prometheus
cd prometheus/
./prometheus &