
# TO-DO App REST API  

## Overview  

This project is a REST API service for a TO-DO application developed using **Golang** and the **Gin** framework. Users can create TO-DO lists, add steps to these lists, and manage them. There are two user roles in the system: **regular user** and **admin**.  

## Features  

* Authentication with JWT  
* Role-based authorization (admin/user)  
* Creating, updating, and soft-deleting TO-DO lists  
* Creating, updating, and deleting steps for each list  
* Soft delete functionality to prevent data loss  
* Admin can access all users' data  

## Technologies Used  

* **Programming Language**: Golang  
* **Framework**: Gin  
* **Authentication**: JWT  
* **Database**: In-memory (Mock repository)  
* **API Testing**: Postman  
* **Frontend**: HTML, CSS, JavaScript  
* **Hosting**: Render (for frontend)  

## User Credentials  

**Regular User**  

    * Username: `user1`  
    * Password: `user123`  

**Admin**  

    * Username: `admin`  
    * Password: `admin123`  

## Installation  

```bash
git clone https://github.com/kullaniciadi/todo-api.git  
cd todo-api  
go mod download  
go run main.go  
```  

## API Endpoints (Testable via Postman)  

### Authentication  

| Method | Endpoint  | Description       |  
|--------|----------|------------------|  
| POST   | `/login` | Retrieve JWT token |  

### TO-DO Lists  

| Method | Endpoint     | Description                              |  
|--------|------------|-----------------------------------------|  
| GET    | `/todos`   | Retrieves the user's TO-DO lists       |  
| POST   | `/todos`   | Creates a new TO-DO list               |  
| PUT    | `/todos/:id` | Updates a TO-DO list                 |  
| DELETE | `/todos/:id` | Soft deletes a TO-DO list            |  

### TO-DO Steps  

| Method | Endpoint     | Description                                              |  
|--------|------------|----------------------------------------------------------|  
| GET    | `/steps`   | Retrieves user's steps (admin can view all)             |  
| POST   | `/steps`   | Creates a new step                                      |  
| PUT    | `/steps/:id` | Updates a step                                        |  
| DELETE | `/steps/:id` | Deletes a step (soft delete)                          |  

> The admin user can view all steps via the `/admin/steps` endpoint.  

## Architecture  

This project follows **Clean Architecture** principles:  

* **Handler (Controller)**: Handles HTTP requests.  
* **Service**: Contains business logic.  
* **Repository**: Manages data access operations.  
* **Model**: Defines data structures.  
* **Middleware**: Handles authorization operations.  
* **pkg/jwt**: Manages token operations.  

## Frontend  

A frontend interface has also been developed for this project using HTML, CSS, and JavaScript. You can view the live demo at the following link:  

🔗 [Live Demo (Render)](https://todo-project-69kz.onrender.com)  

## Contribution  

If you would like to contribute, please fork the repository and submit a pull request. It is recommended to open an issue before making major changes.  

## License  

MIT License  
