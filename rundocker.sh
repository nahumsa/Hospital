# Build docker container
sudo docker build -t hospital . 

# Execute container
sudo docker run -d -p 8080:8080 --name hospital hospital 