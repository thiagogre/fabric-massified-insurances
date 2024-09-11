#!/bin/bash

USERNAME="$1"
PASSWORD="$2"
ORG="${3:-"org1"}"
ABAC_ROLE="${4:-"null"}" # Attribute-based access control [insurer, claimValidator]

export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/${ORG}.example.com/

register() {
    fabric-ca-client register \
        --id.name "${USERNAME}" \
        --id.secret "${PASSWORD}" \
        --id.type client \
        --id.affiliation "${ORG}" \
        --id.attrs "abac.role=${ABAC_ROLE}:ecert" \
        --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem"
}

enroll() {
    # TODO: FIX THIS DUPLICATION
    declare -A ca_ports
    ca_ports["org1"]="7054"
    ca_ports["org2"]="8054"
    ca_ports["org3"]="11054"
    local ca_port=${ca_ports[$ORG]}

    fabric-ca-client enroll \
        -u "https://${USERNAME}:${PASSWORD}@localhost:${ca_port}" \
        --caname "ca-${ORG}" \
        -M "${PWD}/organizations/peerOrganizations/${ORG}.example.com/users/${USERNAME}@${ORG}.example.com/msp" \
        --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem"
}

{
    register
    enroll
    cp "${PWD}/organizations/peerOrganizations/${ORG}.example.com/msp/config.yaml" \
        "${PWD}/organizations/peerOrganizations/${ORG}.example.com/users/${USERNAME}@${ORG}.example.com/msp/config.yaml"
}

echo "Created and enrolled identity for user: ${USERNAME}"
