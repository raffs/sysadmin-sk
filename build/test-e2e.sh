#!/usr/bin/env bash
#
# This file is part of the Sysadmin Sidekick Toolkit (Sysadmin-SK) (https://github.com/raffs/sysadmin-sk).
# Copyright (c) 2019 Rafael Oliveira Silva
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 3.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
# General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.
##

[[ ! $CREATE_CONTAINER ]] && CREATE_CONTAINER=1   # define variable is not already defined elsewhere

if [[ $CREATE_CONTAINER -eq 1 ]]
then
  RUN_CONTAINER=$(docker run -d -it --rm -p 5000:5000 raffs/moto:latest)

  if [[ $RUN_CONTAINER = "" ]]
  then
    echo "ERROR: it seems there's is an error when executing container"
    echo "ERROR: Please check above messages for details"
    exit 1
  fi
fi

# tests starts here
echo "E2E Testing" && (
  # why parentheses '()' you ask?
  # Because everything will be executed in sub-shell, to ensure (at least try :D)
  # execute the cleanup code on bottom of the script.

  sqs_cmd='aws --endpoint http://localhost:5000 sqs'

  echo "Running SQS Move Test" && {
    sourceQueue=$($sqs_cmd create-queue --queue-name moveTestSourceQueue | jq .QueueUrl | sed 's/\"//g')
    targetQueue=$($sqs_cmd create-queue --queue-name moveTestTargetQueue | jq .QueueUrl | sed 's/\"//g')

    for _ in `seq 1 10`
    do
      $sqs_cmd send-message --queue-url "${sourceQueue}" --message-body "Message ID: ${RANDOM}"
    done

    echo "Running SQS Move Test: DONE"
  }
)

if [[ $CREATE_CONTAINER -eq 1 ]]
then
  echo "Cleaning up" && {
    if ! docker rm -f "${RUN_CONTAINER}"; then
      echo "Error when trying to remove the container"
      exit 1
    fi
  }
fi

echo "that's all folks" && exit 0