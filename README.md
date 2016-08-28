Tailla Chat [![Build Status](https://drone.io/github.com/elcct/taillachat/status.png)](https://drone.io/github.com/elcct/taillachat/latest)
===========

Tailla Chat https://tailla.com/ is an anonymous chat and picture sharing web application written in Go.

# Dependencies

Tailla Chat requires `Go`.

Instructions below have been tested on `Ubuntu 16.04`.

## Installation

To install Tailla Chat issue:

```
go get -u github.com/elcct/taillachat
```

In your `GOPATH` directory you can create `config.json` file:

```
{
	"secret": "secret",
	"public_path": "./src/github.com/elcct/taillachat/public",
	"template_path": "./src/github.com/elcct/taillachat/views"
}
```

Finally, you can run:

```
$GOPATH/bin/taillachat
```

That should output something like:

```
2014/06/19 15:31:15.386961 Starting Goji on [::]:8000
```

And it means you can now direct your browser to `localhost:8000`

## Icons

Tailla Chat uses Skype Emoticons via https://github.com/kof/emoticons/

Check license: https://github.com/kof/emoticons/blob/master/LICENSE
