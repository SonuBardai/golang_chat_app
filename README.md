# Golang Chat App

This is a chat application built using React with TypeScript on the frontend and Golang on the backend. The chat server is implemented using the standard websockets library.

## Repository Structure
This repository contains two directories:
- `client`: Contains the React app code for the frontend
- `server`: Contains the Golang code for the chat server

## Features
- Real-time chat using websockets
- Frontend built with React and TypeScript
- Backend built with Golang

## Getting Started
To get started with this project, follow these steps:

1. Clone the repository
```
git clone https://github.com/SonuBardai/golang_chat_app.git
```

2. Install dependencies
```
cd golang_chat_app/client
npm install
```

3. Start the server
```
cd ../server
go run .
```

4. Start the frontend
```
cd ../client
npm start
```

5. Open your browser and navigate to `http://localhost:3000` to start chatting!
