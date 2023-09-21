sudo docker build -t proxy .

sudo docker run -p 8080:8080 -p 8000:8000 -t proxy
