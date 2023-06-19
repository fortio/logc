FROM scratch
COPY logc /usr/bin/logc
ENTRYPOINT ["/usr/bin/logc"]
CMD ["help"]
