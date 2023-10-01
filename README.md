# DevOps
让你轻松快速拥有DevOps环境


docker build -t devops-go:v1 .

docker run -d --name nginx-go \
    --restart=always \
    --privileged=true \
    -p 8080:80 \
    -p 22333:22333 \
    devops-go:v1
