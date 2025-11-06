# Configuration

## Environment Variables

Create a `.env` file in this directory (`config/.env`) with the following required variables:

```env
DATABASE_URL=postgres://postgres:postgres@localhost/blogdb?sslmode=disable
PORT=8080
```

### Variables

- **DATABASE_URL** (required): PostgreSQL connection string
- **PORT** (required): Server port number

### Optional Variables

- **JWT_SECRET**: Secret key for JWT token signing (defaults to a development value if not set)

## Usage

The application will automatically load `config/.env` on startup. If the file is not found, it will try to load `.env` from the project root.

For Docker deployments, you can specify a custom env file path using the `--env-path` flag:

```bash
./blog-api --env-path /path/to/.env
```

