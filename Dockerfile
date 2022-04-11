FROM ubuntu
RUN mkdir /app
WORKDIR /app
COPY bin/traceview-amd64-linux ./
COPY ./content ./content
RUN rm -rf ./content/node_modules
EXPOSE 1969/udp
EXPOSE 8969/tcp
ENTRYPOINT ["./traceview-amd64-linux"]