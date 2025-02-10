# Would You Rather - Programmer Edition


## Live Demo

[https://whatwouldyourather.netlify.app/](https://whatwouldyourather.netlify.app/)

## Features

- Would you rather game logic.
- Get a Random situation pair from server
- See the percentage of people who picked the different situation
- Replay the game
- Select the color theme of your choise.


## Tech Stack

- Fronend: React, TypeScript, TailwindCSS, Netlify
- Backend: Go, Sqlite, Docker, DigitalOcean

## Installation & Running Locally

### Prerequisites

[Node.js](https://nodejs.org/en) (LTS recommended)
[Go](https://go.dev/doc/install)
[Git](https://git-scm.com/downloads)

### Setup

#### Frontend 

```bash 
git clone https://github.com/karlsb/WouldYouRather.git
```

```bash
cd WouldYouRather/WouldYouRatherFrontend
```

```bash
npm install
```

#### Backend

Navigate to server directory

```bash
cd WouldYouRather/WouldYouRatherBackend
```

Create a .env (it can be empty) file inside the WouldYouRather/WouldYouRatherBackend directory.

You can specify the PORT you want to run the server on in the .env file in the following way:

PORT=8081 


Build the server:

```bash
go build -o main
```

### Start Frontend

```bash
npm run dev
```

### Start Backend

```bash
./main
```
