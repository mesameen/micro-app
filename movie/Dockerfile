FROM alpine:latest

# Create and set work directory
WORKDIR /app

# Copy our binary and configs to the container
COPY main ./
COPY configs/default.yaml ./

# Make sure the binary is executable
RUN chmod +x ./main

# EXPOSE the port for accepting incoming requests
EXPOSE 8093

# Execute our service
CMD [ "/app/main" ]
