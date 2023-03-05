FROM ubuntu:focal as base

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moskow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    lsb-release \
    ca-certificates \
    libssl-dev \
    libldap-common \
    gnupg \
    openssl \
    ldap-utils \
    libldap-2.4-2 \
    libldap-dev

RUN curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'

RUN apt-get update && apt-get install -y --no-install-recommends \
    sudo postgresql-14 \
    build-essential \
    cmake \
    clang \
    gdb \
    libpam0g-dev \
    libpq5 \
    libpq-dev \
    postgresql-common \
    postgresql-server-dev-14

COPY ./CMakeLists.txt .
COPY ./cmake ./cmake
COPY ./sources ./sources
COPY ./third_party ./third_party

RUN cmake -S $(pwd) -Bbuild -DCMAKE_BUILD_TYPE=Release -DCC_FLAGS="-Wextra -Wstrict-aliasing" -DUSE_SCRAM=YES
RUN make -Cbuild -j8

ENTRYPOINT ["./build/sources/odyssey", "/etc/odyssey/odyssey.conf", "--verbose", "--console", "--log_to_stdout"]