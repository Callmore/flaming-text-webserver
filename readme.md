# Webserver for APIS for use in Neos

A clone of the flaming text webserver.

A collection of APIS for use within neos.

## APIS

### Flaming Text

Generator/WebServer that generates flaming text, intended for use with NeosVR. Using the command line flag `-text` allows for the generation of any text without starting the webserver, and writes the output to a folder named `out` if it exists (Will probably create that folder automaticly in the future).

### Slot machine

Grabs two images and a short description and returns it so that a slot machine can be generated!

## Dependencies

- github.com/joho/godotenv - DotEnv implementation
- github.com/ojrac/opensimplex-go - OpenSimplex implementation for go
