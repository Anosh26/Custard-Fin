# Multi-stage build for CustardFin Next.js Web UI
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies based on package-lock.json
COPY apps/web/package.json apps/web/package-lock.json ./
RUN npm ci

# Copy application source code
COPY apps/web/ ./

# Disable telemetry during build
ENV NEXT_TELEMETRY_DISABLED=1

# Build production bundle
RUN npm run build

# Minimal runner image
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV NEXT_TELEMETRY_DISABLED=1

# Copy built artifacts and production dependencies
COPY --from=builder /app/package.json /app/package-lock.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3000

CMD ["npm", "start"]
