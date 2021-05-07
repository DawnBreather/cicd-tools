FROM maven:3-openjdk-16

COPY --from=dawnbreather/deploy_kaniko /kaniko/executor /kaniko/executor
COPY --from=dawnbreather/deploy_kaniko /kaniko/.docker/config.json /root/.docker/config.json

COPY --from=dawnbreather/envsubst /usr/bin/envsubst /usr/bin/envsubst
COPY --from=dawnbreather/envsubst /usr/bin/envmake /usr/bin/envmake

COPY --from=dawnbreather/deploy_invoker /usr/bin/deploy-invoker /usr/bin/deploy-invoker

RUN curl https://amazon-ecr-credential-helper-releases.s3.us-east-2.amazonaws.com/0.5.0/linux-amd64/docker-credential-ecr-login --output /usr/bin/docker-credential-ecr-login \
    && chmod +x /usr/bin/docker-credential-ecr-login

RUN rm -rf /var/lock
