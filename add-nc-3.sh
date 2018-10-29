#!/bin/bash

# setup client 3
echo '###################'
echo 'setting up client 3'
echo '###################'
rsync -rtuv setup-client.sh jbe@nc-3:/home/jbe/
ssh -t jbe@nc-3 'sudo apt update && sudo apt upgrade -y && sudo sh /home/jbe/setup-client.sh'
