# Hello World

This document summarizes how to build, deploy and run the helloworld
program with a GCP IoT Kit. Once you are finished, you will have a display
blinking "Hello gopher!" with a blue backlight.

![Hello gopher!](https://i.imgur.com/TKI5Iz1.gif)

## Building

First, go get the sample program.

    $ go get -d github.com/GoogleCloudPlatform/golang-samples/iotkit/helloworld

Then, build a binary targeting the GCP IoT Kit board, BeagleBone Green
Wireless, by running the command below.

    $ GOOS=linux GOARCH=arm GOARM=7 go build github.com/GoogleCloudPlatform/golang-samples/iotkit/helloworld

The command below will generate a binary that can be executed on your board.

## Setting up the board

Download [Debian 8.5 2016-05-13 4GB SD LXQT](https://beagleboard.org/latest-images)
and write it to a microSD card.

Put the microSD card into the SDcard slot and power the board to boot.
Once the board boots, you will be able to see an unsecured Wifi network
in the format of "BeagleBoneXXXX".

Join to the "BeagleBoneXXXX" network from your laptop and visit
[192.168.8.1:3000/ide.html](http://192.168.8.1:3000/ide.html) to launch the
Cloud9 IDE.

Use "Files > Upload Local Files" to launch the uploader dialog.
Drag and drop the "helloworld" binary and wait until the upload is finished.

From the Cloud9 IDE terminal, run the following command to make the
binary executable.

    $ chmod +x ./helloworld

## Connect the Grove LCD RGB Backlight

Grove LCD RGB Backlight is an I2C device that requires 5V.
Connect the device to the board using the I2C and the 5V pinouts as shown on
the screenshot below. If you use 3.3V, the characters won't be displayed.

![Grove LCD RGB Backlight connection](https://i.imgur.com/8dnySQn.jpg)

## Running the program

Once the display is connected, you can return back to your Cloud9 IDE
terminal to start the `helloworld` program.

    $ ./helloworld

It will turn on the display, set the backlight color to blue and blink
"Hello gopher!".
