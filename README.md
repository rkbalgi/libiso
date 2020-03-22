# A library of utilities related to payments, crypto, ISO8583 etc 

## UPDATE (03/22/2020)
1. There will be no further development on paysim, this repo will be strictly be used as a module/library
2. If you're interested in a ISO8583 simulator, please check out [ISO WebSim](https://github.com/rkbalgi/isosim)

## UPDATE (06/16/2019)
1. Folks developing on Windows please see this - https://github.com/rkbalgi/go/wiki/Building-on-Windows
2. Doesn't follow standard go style (coding conventions etc) - WIP
2. This has not be subject to any kind of targeted tests (performance or otherwise), so use this with a bit of caution - It's at the moment perhaps best suited for simulators


# Paysim
An open ISO8583 Simulator

<ul>
<li>The application is built using go and GTK+2 bindings made available at github.com/mattn/go-gtk (Thanks a ton Yashuhiro Matsumoto!)</li>
<li>The entire source code is available at https://github.com/rkbalgi/go</li>
<li>The interesting packages would be github.com/rkbalgi/go/execs/paysim, github.com/rkbalgi/go/paysim and github.com/rkbalgi/iso8583</li>
<li>There are loads of other interesting things available in other packages – like a minimalist implementation of a Thales HSM for basic commands (A6, MS, M6 and the like)
</li>
</ul>

You can read more about paysim here - https://github.com/rkbalgi/go/wiki/Paysim
