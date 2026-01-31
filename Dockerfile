# Simple runtime image for Snake Game Services
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

# Create bin directory
RUN mkdir -p ./bin

# Copy pre-built binaries
COPY bin/lobby_server ./bin/lobby_server
COPY bin/matching_server ./bin/matching_server
COPY bin/room_server ./bin/room_server
COPY bin/leaderboard_server ./bin/leaderboard_server
COPY bin/game_server ./bin/game_server
COPY bin/friends_server ./bin/friends_server
COPY bin/gateway_server ./bin/gateway_server

# Make binaries executable
RUN chmod +x ./bin/*

# Expose ports
EXPOSE 8080 50051 50052 50053 50054 50055 50056

CMD ["./bin/gateway_server"]