FROM debian:stretch

ENV CLOUD_SDK_VERSION=235.0.0
ENV INSTALL_COMPONENTS="google-cloud-sdk-app-engine-go google-cloud-sdk-datastore-emulator"
ENV GOLANG_VERSION 1.11.5

RUN apt-get update -qqy && apt-get install -qqy \
        curl \
        g++ \
    		gcc \
    		libc6-dev \
    		make \
    		pkg-config \
        python-dev \
        python-setuptools \
        apt-transport-https \
        lsb-release \
        openssh-client \
        git \
        gnupg \
    && easy_install -U pip && \
    pip install -U crcmod

RUN export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
    echo "deb https://packages.cloud.google.com/apt ${CLOUD_SDK_REPO} main" > /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    apt-get update && apt-get install -y google-cloud-sdk=${CLOUD_SDK_VERSION}-0 ${INSTALL_COMPONENTS} && \
    gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image && \
    gcloud --version

RUN set -eux; \
    	url="https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz"; \
      \
    	curl --location --output go.tgz "$url"; \
    	tar -C /usr/local -xzf go.tgz; \
    	rm go.tgz; \
    	\
    	export PATH="/usr/local/go/bin:$PATH"; \
    	go version

RUN pip install -U grpcio

ENV PATH /usr/local/go/bin:$PATH

VOLUME ["/root/.config"]
