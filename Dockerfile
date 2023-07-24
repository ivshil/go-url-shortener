# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go app inside the container
RUN go build -o main .

# Expose the port on which the Go app will listen
EXPOSE 1337

# Set environment variables for the PostgreSQL connection
ENV PGHOST=postgres
ENV PGPORT=5432
ENV PGUSER=admin
ENV PGPASSWORD=url_short
ENV PGDATABASE=postgres_url_short_db

# Run the Go app
CMD ["./main"]