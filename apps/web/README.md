# apps/web

CustardFin Web Client (Next.js App Router).

## Purpose
- Media discovery catalog and video streaming frontend.
- React Server Components (RSC) for metadata/catalog fetching directly from Go API.
- Client component video player utilizing `hls.js`.
- Strictly interfaces with Go API (`NEXT_PUBLIC_API_URL`); zero direct DB connections.

## Development
```bash
npm install
npm run dev
```
