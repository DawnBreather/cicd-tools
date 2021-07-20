FROM mcr.microsoft.com/dotnet/core/sdk:3.1

COPY --from=dawnbreather/deploy_kaniko /kaniko/executor /kaniko/executor
COPY --from=dawnbreather/deploy_kaniko /kaniko/.docker/config.json /root/.docker/config.json

COPY --from=dawnbreather/buildtools /usr/bin/envsubst /usr/bin/envsubst
COPY --from=dawnbreather/buildtools /usr/bin/envmake /usr/bin/envmake
COPY --from=dawnbreather/buildtools /usr/bin/setsubst /usr/bin/setsubst
COPY --from=dawnbreather/buildtools /usr/bin/set2secret /usr/bin/set2secret
COPY --from=dawnbreather/buildtools /usr/bin/k8sdeploy /usr/bin/k8sdeploy

COPY --from=dawnbreather/deploy_invoker /usr/bin/deploy-invoker /usr/bin/deploy-invoker

RUN curl https://amazon-ecr-credential-helper-releases.s3.us-east-2.amazonaws.com/0.5.0/linux-amd64/docker-credential-ecr-login --output /usr/bin/docker-credential-ecr-login \
    && chmod +x /usr/bin/docker-credential-ecr-login

RUN apt update \
    && apt install curl unzip -y \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf ./aws \
    && rm -rf /var/lib/apt/lists/*

RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl \
    && curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 \
    && chmod +x get_helm.sh && ./get_helm.sh

RUN rm -rf /var/lock
