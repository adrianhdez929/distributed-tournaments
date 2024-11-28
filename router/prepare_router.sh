# Check if router image already exists
if ! docker images router -q | grep -q .; then
    docker build -t router -f router/router.Dockerfile .
fi

docker run -itd --rm --name router router
docker network connect --ip 10.0.10.254 clients router
docker network connect --ip 10.0.11.254 servers router