# Use the official Node.js image as the base image
FROM node:18-alpine AS base

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json (if available)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application files
COPY . .

# Build the Next.js app
RUN npm run build

# Production image
FROM node:18-alpine AS production

# Set the working directory inside the container
WORKDIR /app

# Install only production dependencies
COPY package*.json ./
RUN npm install --only=production

# Copy the built app from the build stage
COPY --from=base /app/.next ./.next
COPY --from=base /app/public ./public
COPY --from=base /app/package.json ./package.json

# Expose port 3000
EXPOSE 3000

# Start the Next.js app
CMD ["npm", "run", "start"]