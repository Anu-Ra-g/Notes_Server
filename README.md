# Notes_Server

This is a REST API made with Golang and Gin Framework. This API uses JWT for user authentication and PostgreSQL database for data persistence. The database schema includes users and notes tables. 

## API endpoints

- `POST /signup`  
- `POST /login`  
- `GET /notes` 🔒
- `POST /notes` 🔒 
- `DELETE /notes` 🔒

## To run the code
1. Pull the docker image
   `docker pull anurag101/notes_api:latest`
2. Start the container at port 3000 and listen on port 3000
   `docker run -p 3000:3000 anurag101/notes_api:latest`
3. Now you can use `curl` or an API client like Postman to make requests to this API.
