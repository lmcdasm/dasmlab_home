app=dasmlab-home
ver=latest

docker stop ${app}-instance
docker rm ${app}-instance
docker run -d \
	--restart=always \
	--name=${app}-instance \
	-p 8111:80 ${app}:${ver}
