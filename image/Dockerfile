FROM fedora

COPY whitelist /opt/dataverse-broker/whitelist
ADD dataverse-broker /opt/dataverse-broker/dataverse-broker
CMD /opt/dataverse-broker/dataverse-broker --port 8080 -v 4 --logtostderr
