# Notes on using the Samsung Galaxy s7 (and similar) Keyboard Cover with termux and ssh

I have been a diehard hardware keyboard user on Android phones
starting from the HTC Dream, Desire and finally the Galaxy S Relay
which seemed to be the end of the line...

Recently I found out that Samsung is making a Keyboard Cover add-on
for the newer Galaxy S series phones including the S7 and the S8 and
picked one of these up to replace my ancient Relay.

My main use case for the keyboard is ssh. On the S7 I use https://termux.com/ and
this setup is workable out of the box with some annoyances and
limitations:

## Ctrl and ESC
The keyboard cover has no way to generate Ctrl shifted characters or
ESC. Termux allows working around this by using the volume up/down
buttons or on-screen buttons which is workable but not perfect. 

Termux also allows configuring the Back key to send ESC by putting the
following lines in .termux/termux.properties:

    bell-character=ignore
    back-key=escape

This takes care of ESC. For Ctrl I decided to use the keyboard's "$"
as a Ctrl prefix. This could be done by modifying Termux and maybe I
will try that later but for now I decided to handle it on the Linux
side. At first I attempted to do it with the tmux configuration in
ctrlhack.tmux

This mostly works but it doesn't allow sending the tmux escape
sequence itself using the $ key. So I wrote a small pty wrapper to do
the translation that which be found in ctrlhack.go

Since it's written in Go and uses only pure Go libraries it is trivial
to cross-compile to an Arm binary that will run on Android:

    go get -d .
    GOARCH=arm go build -o ctrlhack-arm

It will launch bash or the executable specified in the first
argument. In the child shell (or process) you can type "$c" for Ctrl-C
and enter other Ctrl shifted characters similarly. "$$" will emit a
real "$". 

