# Marbles App

Graph equations and run marbles on them. Based on [desmos.com](https://desmos.com). Uses [goki/gi](https://github.com/goki/gi) for graphics, and [Knetic/govaluate](https://github.com/Knetic/govaluate) for evaluating equations.  

## Install

To install run 
``` bash
$ go install github.com/kplat1/marbles@latest
```
Once you have done this you should be able to launch marbles by just doing
```bash
$ marbles
```
If the install does not work, check the [GoKi Install Page](https://github.com/goki/gi/wiki/Install) and make sure you have installed the prerequisites. If the widgets example doesn't work, then marbles won't work. 

If there is a new version of marbles released, just run this command again to update:
``` bash
$ go install github.com/kplat1/marbles@latest
```

