#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

# imports
. scripts/utils.sh

# export CORE_PEER_TLS_ENABLED=true
# export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
# export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
# export PEER0_ORG2_CA=${PWD}/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
# export PEER0_ORG3_CA=${PWD}/organizations/peerOrganizations/org3.example.com/tlsca/tlsca.org3.example.com-cert.pem
# export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
# export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

# Set environment variables for the peer org
setGlobalsForPeer() {
  local USING_ORG=""
  local PEER_PORT=$2
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  infoln "Using organization ${USING_ORG}"

  export CORE_PEER_LOCALMSPID="Org${USING_ORG}MSP"
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
  elif [ $USING_ORG -eq 3 ]; then
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
  fi

  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org${USING_ORG}.example.com/users/Admin@org${USING_ORG}.example.com/msp
  export CORE_PEER_ADDRESS=localhost:${PEER_PORT}

  echo ---------
  echo $CORE_PEER_LOCALMSPID
  echo $CORE_PEER_TLS_ROOTCERT_FILE
  echo $CORE_PEER_MSPCONFIGPATH
  echo $CORE_PEER_ADDRESS
  echo ---------

  if [ $USING_ORG -gt 3 ]; then
    errorln "ORG Unknown"
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

# # Set environment variables for use in the CLI container
# setGlobalsCLI() {
#   setGlobals $1

#   local USING_ORG=""
#   if [ -z "$OVERRIDE_ORG" ]; then
#     USING_ORG=$1
#   else
#     USING_ORG="${OVERRIDE_ORG}"
#   fi
#   if [ $USING_ORG -eq 1 ]; then
#     export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
#   elif [ $USING_ORG -eq 2 ]; then
#     export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
#   elif [ $USING_ORG -eq 3 ]; then
#     export CORE_PEER_ADDRESS=peer0.org3.example.com:11051
#   else
#     errorln "ORG Unknown"
#   fi
# }

# # parsePeerConnectionParameters $@
# # Helper function that sets the peer connection parameters for a chaincode
# # operation
# parsePeerConnectionParameters() {
#   PEER_CONN_PARMS=()
#   PEERS=""
#   while [ "$#" -gt 0 ]; do
#     setGlobals $1
#     PEER="peer0.org$1"
#     ## Set peer addresses
#     if [ -z "$PEERS" ]
#     then
# 	PEERS="$PEER"
#     else
# 	PEERS="$PEERS $PEER"
#     fi
#     PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" --peerAddresses $CORE_PEER_ADDRESS)
#     ## Set path to TLS certificate
#     CA=PEER0_ORG$1_CA
#     TLSINFO=(--tlsRootCertFiles "${!CA}")
#     PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" "${TLSINFO[@]}")
#     # shift by one to get to the next organization
#     shift
#   done
# }

# verifyResult() {
#   if [ $1 -ne 0 ]; then
#     fatalln "$2"
#   fi
# }
