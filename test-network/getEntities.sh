#!/bin/bash

ORG=${1:-org1}

ls -p -w 1 ./organizations/peerOrganizations/${ORG}.example.com/users
