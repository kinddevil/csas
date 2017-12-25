FROM alpine

# Define working directory.
WORKDIR /opt/
ADD app /opt/
# Only copy files in folder
ADD templates /opt/templates
# ADD ../templates

EXPOSE 6767
ENV env=prod

# Define default command.
CMD ["./app"]

