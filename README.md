## 📂 **Project Structure from Fantastic Coffee**

###
docker build -t my-go-app -f Dockerfile.backend .


docker run -p 3000:3000 --name my-running-app my-go-app


docker build -t my-frontend -f Dockerfile.frontend .


docker run -p 4000:80 --name frontend-dev my-frontend

