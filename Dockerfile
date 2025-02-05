FROM alpine:3.16.0

USER nobody
# work somewhere where we can write
COPY tfsec /usr/bin/tfsec
# set the default entrypoint -- when this container is run, use this command
ENTRYPOINT [ "tfsec" ]
# as we specified an entrypoint, this is appended as an argument (i.e., `tfsec --help`)
CMD [ "--help" ]
