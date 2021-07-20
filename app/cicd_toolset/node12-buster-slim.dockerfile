FROM node:12-buster-slim

COPY --from=dawnbreather/deploy_kaniko /kaniko/executor /kaniko/executor
COPY --from=dawnbreather/deploy_kaniko /kaniko/.docker/config.json /root/.docker/config.json

COPY --from=dawnbreather/envsubst /usr/bin/envsubst /usr/bin/envsubst
COPY --from=dawnbreather/envsubst /usr/bin/envmake /usr/bin/envmake
COPY --from=dawnbreather/envsubst /usr/bin/setsubst /usr/bin/setsubst
COPY --from=dawnbreather/envsubst /usr/bin/set2secret /usr/bin/set2secret

COPY --from=dawnbreather/deploy_invoker /usr/bin/deploy-invoker /usr/bin/deploy-invoker

RUN apt update \
    && apt install curl unzip -y \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf ./aws \
    && rm -rf /var/lib/apt/lists/*

RUN curl https://amazon-ecr-credential-helper-releases.s3.us-east-2.amazonaws.com/0.5.0/linux-amd64/docker-credential-ecr-login --output /usr/bin/docker-credential-ecr-login \
    && chmod +x /usr/bin/docker-credential-ecr-login

RUN rm -rf /var/lock
