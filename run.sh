docker stop livy
docker rm livy
docker build -t livy:latest .
docker run --name livy -d -p 5000:5000 livy
