# Would You Rather - Programmer Edition


## Live Demo

[https://whatwouldyourather.netlify.app/](https://whatwouldyourather.netlify.app/)

## Features

- Would you rather game logic.
- Get a Random situation pair from server
- See the percentage of people who picked the different situation
- Replay the game
- Select the color theme of your choice.


## Tech Stack

- Frontend: React, TypeScript, TailwindCSS, Vercel
- Backend: Go, SQLite, Docker, Fly.io

## Installation & Running Locally

### Prerequisites

- [Node.js](https://nodejs.org/en) (LTS recommended)
- [Go](https://go.dev/doc/install)
- [Git](https://git-scm.com/downloads)

### Setup

#### Frontend 

```bash 
git clone https://github.com/karlsb/WouldYouRather.git
```

```bash
cd WouldYouRather/WouldYouRatherClient
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

## Deployment

- **Frontend:** Vercel builds `WouldYouRatherClient` on every push. `vercel.json` proxies `/api/*` to the backend.
- **Backend:** pushes to `main` that change `WouldYouRatherBackend/` are tested and deployed to Fly.io by `.github/workflows/deploy-backend.yml`. To deploy manually:

```bash
cd WouldYouRatherBackend
fly deploy --ha=false
```

The SQLite database lives on the Fly volume `wyr_data` at `/data/wouldyourather.db`. On first boot it's seeded from `build-database/wouldyourather.db`. Later deploys never overwrite it.
