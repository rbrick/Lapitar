FROM golang:1.8.3

# Instal required dependencies
RUN apt-get update -qq && apt-get install -y build-essential pkg-config libosmesa6-dev libglu1-mesa-dev

# Install glide
RUN curl https://glide.sh/get | sh

# Set the working directory
WORKDIR "/go/src/github.com/FrozenOrb/lapitar"

# Add everything that is in this folder to the SRC folder
COPY . "/go/src/github.com/FrozenOrb/lapitar"

# Install the dependencies
RUN glide install

EXPOSE 8088

CMD ["go", "run" , "lapitar/main.go", "server"]