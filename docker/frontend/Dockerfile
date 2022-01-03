ARG NODE_VERSION=16-alpine

FROM node:${NODE_VERSION}

LABEL maintainer="Neekity"

VOLUME /data

WORKDIR /data

ARG NPM_REGISTRY
ENV NPM_REGISTRY ${NPM_REGISTRY}

EXPOSE 8881

CMD ["/bin/sh", "-c", "yarn install && yarn run serve"]