# This Dockerfile creates an image that includes all the
# dependencies necessary to build Jswap from sources.

FROM golang:1.22-bookworm

RUN apt-get update &&\
    apt-get install -y unzip nsis

# Install the EnVar plugin for NSIS
RUN wget https://nsis.sourceforge.io/mediawiki/images/7/7f/EnVar_plugin.zip &&\
    unzip EnVar_plugin.zip -d EnVar_plugin && \
    mv EnVar_plugin/Plugins/amd64-unicode/EnVar.dll /usr/share/nsis/Plugins/amd64-unicode/EnVar.dll &&\
    mv EnVar_plugin/Plugins/x86-ansi/EnVar.dll /usr/share/nsis/Plugins/x86-ansi/EnVar.dll &&\
    mv EnVar_plugin/Plugins/x86-unicode/EnVar.dll /usr/share/nsis/Plugins/x86-unicode/EnVar.dll &&\
    rm EnVar_plugin.zip &&\
    rm -r EnVar_plugin

WORKDIR /app

ENTRYPOINT ["make"]
CMD ["all"]
