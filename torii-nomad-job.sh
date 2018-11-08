#!/bin/bash

# setup server 1
echo '###################'
echo 'setting up server 1'
echo '###################'
rsync -rtuv setup-server.sh jbe@ns-1:/home/jbe/
rsync -rtuv ./vault jbe@ns-1:/home/jbe
rsync -rtuv ./nomad jbe@ns-1:/home/jbe
ssh -t jbe@ns-1 'sudo apt update && sudo apt upgrade -y && sudo sh /home/jbe/setup-server.sh'


# setup client 1
echo '###################'
echo 'setting up client 1'
echo '###################'
rsync -rtuv setup-client.sh jbe@nc-1:/home/jbe/
ssh -t jbe@nc-1 'sudo apt update && sudo apt upgrade -y && sudo sh /home/jbe/setup-client.sh'

# setup client 2
echo '###################'
echo 'setting up client 2'
echo '###################'
rsync -rtuv setup-client.sh jbe@nc-2:/home/jbe/
ssh -t jbe@nc-2 'sudo apt update && sudo apt upgrade -y && sudo sh /home/jbe/setup-client.sh'
