# Base image
FROM node:18-alpine AS base

# Install dependencies
WORKDIR /app
COPY package*.json ./
RUN yarn install

# Copy TypeScript code and compile
COPY . .
RUN npm run build

# Expose port and start the server
EXPOSE 8000
CMD ["npm", "start"]