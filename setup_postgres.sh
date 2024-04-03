#!/bin/bash

# Configuration variables
CONTAINER_NAME="user_management"
POSTGRES_USER=""
POSTGRES_PASSWORD=""
POSTGRES_DB="user_management_db"
TABLE_USERS_CREATION_SQL="CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email VARCHAR(255) UNIQUE NOT NULL, password_hash CHAR(60) NOT NULL, created_at TIMESTAMP DEFAULT NOW());"
TABLE_CODES_CREATION_SQL="CREATE TABLE IF NOT EXISTS invitation_codes (id SERIAL PRIMARY KEY, code VARCHAR(255) UNIQUE NOT NULL, used BOOLEAN DEFAULT false, email VARCHAR(255) NOT NULL, expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '2 minutes');"
TABLE_SESSIONS_CREATION_SQL="CREATE TABLE IF NOT EXISTS sessions (id SERIAL PRIMARY KEY, user_id INT NOT NULL, token TEXT NOT NULL UNIQUE, created_at TIMESTAMP DEFAULT NOW(), expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '2 hours', FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);"

# Check if container already exists
if [ $(docker ps -a -f name=^/${CONTAINER_NAME}$ --format '{{.Names}}') == $CONTAINER_NAME ]; then
    echo "Container $CONTAINER_NAME already exists. Stopping and removing it."
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
fi

# Pull the PostgreSQL Docker image
echo "Pulling the PostgreSQL Docker image..."
docker pull postgres

# Run the PostgreSQL container
echo "Running the PostgreSQL Docker container..."
docker run --name $CONTAINER_NAME -e POSTGRES_USER=$POSTGRES_USER -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -e POSTGRES_DB=$POSTGRES_DB -p 5432:5432 -d postgres

# Wait for PostgreSQL to start
echo "Waiting for PostgreSQL to start..."
sleep 10

# Execute the SQL to create the users table
echo "Creating the users table in the $POSTGRES_DB database..."
docker exec -it $CONTAINER_NAME psql -U $POSTGRES_USER -d $POSTGRES_DB -c "$TABLE_USERS_CREATION_SQL"

# Execute the SQL to create the invitation_codes table
echo "Creating the invitation_codes table in the $POSTGRES_DB database..."
docker exec -it $CONTAINER_NAME psql -U $POSTGRES_USER -d $POSTGRES_DB -c "$TABLE_CODES_CREATION_SQL"

# Execute the SQL to create the sessions table
echo "Creating the sessions table in the $POSTGRES_DB database..."
docker exec -it $CONTAINER_NAME psql -U $POSTGRES_USER -d $POSTGRES_DB -c "$TABLE_SESSIONS_CREATION_SQL"


echo "PostgreSQL setup completed successfully."

