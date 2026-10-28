# ipx (https://github.com/unjs/ipx) has no official image. This one mirrors
# the imgproxy container used by the benchmark: a glibc base (sharp's
# prebuilt libvips is tuned for glibc; musl's allocator hurts it), one
# process listening on 8080, config via environment variables.
ARG NODE_VERSION=22
FROM node:${NODE_VERSION}-bookworm-slim

ARG IPX_VERSION=4.0.0
RUN npm install -g --omit=dev "ipx@${IPX_VERSION}" \
  && npm cache clean --force

ENV NODE_ENV=production \
    HOST=0.0.0.0 \
    PORT=8080

# Filesystem storage root (IPX_FS_DIR); unused by the benchmark, which only
# fetches over HTTP from the allowlisted origin (IPX_HTTP_DOMAINS).
WORKDIR /srv
USER node
EXPOSE 8080

CMD ["ipx", "serve"]
