Tailla Chat
===========

Tailla Chat https://tailla.com/ is an anonymous chat and picture sharing web application written in Go.

# Dependencies

Tailla Chat requires `Go` 1.7.

Instructions below have been tested on `Ubuntu 16.04`.

## Installation

To install Tailla Chat issue:

```
go get -u github.com/elcct/taillachat
```

Then you start it:

```
TAILLA_PUBLIC_PATH=$GOPATH/src/github.com/elcct/taillachat/public TAILLA_TEMPLATE_PATH=$GOPATH/src/github.com/elcct/taillachat/views $GOPATH/bin/taillachat
```

You can now direct your browser to `localhost:8000`

## Icons

Tailla Chat uses Skype Emoticons via https://github.com/kof/emoticons/

Check license: https://github.com/kof/emoticons/blob/master/LICENSE
