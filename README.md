# Catomat

Webserver to run on a raspberry pi and control a stepper motor to move it to a
given degree right or left.

## Build
To cross compile for ARM, run:

    GOOS=linux GOARCH=arm go build

