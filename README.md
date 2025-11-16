# Samuel's Forecasts

A full-stack web application for creating, displaying, and tracking forecasts. Built as a learning project to improve forecasting skills and software development expertise.

Visit [samuelsforecasts.com](https://samuelsforecasts.com) to see the live application.

## Tech Stack

### Backend
- **Language**: Go 1.22.5
- **Database**: PostgreSQL
- **Key Libraries**:
  - `pgx/v5` - PostgreSQL driver
  - `golang-jwt/jwt` - Authentication
  - `golang.org/x/crypto` - Password hashing

### Frontend
- **Framework**: React 18.2
- **Build Tool**: Vite
- **UI Library**: Material-UI (MUI)
- **Routing**: React Router v6
- **Charts**: Chart.js, Recharts
- **Styling**: Emotion
- **Markdown**: marked-react
- **Math Rendering**: KaTeX, react-latex

## Project Structure

```
forecasting_project/
├── backend/
│   ├── cmd/                    # Command-line tools
│   ├── internal/
│   │   ├── auth/              # Authentication logic
│   │   ├── cache/             # Caching layer
│   │   ├── database/          # Database connection
│   │   ├── handlers/          # HTTP handlers
│   │   ├── middleware/        # HTTP middleware
│   │   ├── models/            # Data models
│   │   ├── repository/        # Data access layer
│   │   ├── routes/            # Route definitions
│   │   └── services/          # Business logic
│   ├── main.go                # Application entry point
│   └── go.mod                 # Go dependencies
│
└── frontend/
    ├── public/                # Static assets
    ├── src/
    │   ├── components/        # React components
    │   ├── context/           # React context (AuthContext)
    │   ├── pages/             # Page components
    │   ├── services/          # API service layer
    │   ├── utils/             # Utility functions
    │   ├── App.js             # Main application component
    │   └── index.js           # Application entry point
    ├── package.json           # Node dependencies
    └── vite.config.js         # Vite configuration
```

## Architecture

The application follows a clean architecture pattern:

### Backend
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic separation
- **Handler Layer**: HTTP request handling
- **Middleware**: CORS, request logging, authentication
- **Caching**: In-memory cache for performance optimization

### Frontend
- **Component-Based**: Modular React components
- **Context API**: Global state management (authentication)
- **Protected Routes**: Route guards for authenticated pages
- **Theme**: Material-UI theming for consistent design

## Features

- User authentication (login/register)
- Create and manage forecasts
- Track forecast scores
- Visualize forecasts with charts
- Admin dashboard
- FAQ page
- Blog functionality
- Markdown and LaTeX support for rich content

## Getting Started

### Prerequisites

- Go 1.22.5 or higher
- Node.js and npm
- PostgreSQL database

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   Create a `.env` file with:
   ```
   DB_CONNECTION_STRING=your_postgresql_connection_string
   ```

4. Run the backend server:
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Set up environment variables (if needed):
   Create a `.env` file with your API endpoint configuration

4. Run the development server:
   ```bash
   npm run dev
   ```

   The application will be available at `http://localhost:5173`

5. Build for production:
   ```bash
   npm run build
   ```

## Testing

### Backend Tests
Run Go tests:
```bash
cd backend
go test ./...
```

### Frontend Tests
Run React tests:
```bash
cd frontend
npm test
```

## API Endpoints

The backend provides RESTful API endpoints for:
- `/forecasts` - Forecast management
- `/users` - User management
- `/scores` - Score tracking
- `/news` - News/blog content

## Deployment

The application is deployed at [samuelsforecasts.com](https://samuelsforecasts.com)

- **Frontend**: Deployed on Vercel
- **Backend**: Google Cloud Platform (App Engine)

## Development

This project serves as a learning platform for:
- Full-stack web development
- Go backend development
- React frontend development
- PostgreSQL database design
- RESTful API design
- Authentication and authorization
- Deployment and DevOps

## License

Personal project by Samuel Svensson 