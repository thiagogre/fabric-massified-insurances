#!/bin/bash

#
#   This script has only been tested with Debian and GNOME Terminal
#   Only the function 'run_command_in_new_tab' relies on it
#

display_message() {
    local RED='\033[31m'
    local GREEN='\033[32m'
    local YELLOW='\033[33m'
    local BLUE='\033[34m'
    local RESET='\033[0m'

    local message_type="$1"
    local message="$2"

    # Define a single associative array containing color and label as a single string
    declare -A message_info
    message_info["ERROR"]="$RED [ERROR]"
    message_info["INFO"]="$YELLOW [INFO]"
    message_info["SUCCESS"]="$GREEN [SUCCESS]"

    local color_label=${message_info[$message_type]}
    # Extract everything before the first space
    local color="${color_label%% *}"
    # Extract everything after the first space
    local label="${color_label#* }"

    echo -e "${color}${label} ${message}${RESET}"
}

run_command() {
    local command="$1"
    display_message "INFO" "Running command: $command"
    eval $command
    if [ $? -ne 0 ]; then
        display_message "ERROR" "Command failed: $command"
        exit 1
    fi
}

run_command_in_new_tab() {
    local command="$1"
    local tab_name="$2"
    gnome-terminal --tab --title="$tab_name" -- bash -c "$command; exec bash"
}

create_identities() {
    local org="$1"

    local users=("BackendClient" "TestUser")
    pwd
    for user in "${users[@]}"; do
        run_command "./registerEnrollIdentity.sh $user $org"
        display_message "SUCCESS" "Registered and enrolled: $user at $org"
    done
}

# Bring down the network, then bring it up with certificate authority, a channel and CouchDB
run_command "./network.sh down"

run_command "./network.sh up createChannel -c mychannel -ca -r 10 -d 3 -verbose -s couchdb"

run_command_in_new_tab "./monitordocker.sh" "Monitor Docker"

# Deploy the chaincode
run_command "./network.sh deployCC -ccn basic -ccp ./chaincode-go -ccl go"

display_message "INFO" "Creating identities..."
create_identities org1
create_identities org2
create_identities org3

run_command "./getEntities.sh org1"
run_command "./getEntities.sh org2"
run_command "./getEntities.sh org3"

run_command_in_new_tab "cd rest-api-go && go run main.go" "Rest API"

run_command_in_new_tab "cd frontend-react && yarn dev" "Frontend"
