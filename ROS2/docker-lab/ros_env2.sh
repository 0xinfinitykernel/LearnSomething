#!/bin/bash
set -e

# setup ros environment
source  "/opt/ros/humble/setup.bash" --
#source "/home/parallels/catkin_ws/devel/setup.bash" --
#source "/opt/cartographer_ros/setup.bash" --
exec "$@"
