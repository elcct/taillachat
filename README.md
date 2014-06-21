Tailla Chat
===========

Tailla Chat http://tailla.com/ is an anonymous chat and picture sharing web application written in Go.

Tailla Chat has been built using `Default Project` - http://defaultproject.com/

It consists of 4 core components:

- Goji - A web microframework for Golang - http://goji.io/
- Gorilla web toolkit sessions - cookie (and filesystem) sessions - http://www.gorillatoolkit.org/pkg/sessions
- SockJS - WebSocket emulation - https://github.com/sockjs
- mgo - MongoDB driver for the Go language - http://labix.org/mgo

Note: MongoDB is currently not used for anything, but it will be to store messages and perhaps user accounts.

# Dependencies

Tailla Chat requires `Go`, `MongoDB` and few other tools installed.

Instructions below have been tested on `Ubuntu 14.04`.

## Installation

If you don't have `Go` installed, follow installation instructions described here: http://golang.org/doc/install

Then install remaining dependecies:

```
sudo apt-get install make git mercurial subversion bzr
```

MongoDB:

```
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv 7F0CEB10
sudo echo 'deb http://downloads-distro.mongodb.org/repo/debian-sysvinit dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list
sudo apt-get update
sudo apt-get install mongodb-org
```

No go to your GOPATH location and run:

```
go get github.com/elcct/taillachat
```

And then:

```
go install github.com/elcct/taillachat
```

In your GOPATH directory you can create `config.json` file:

```
{
	"secret": "secret",
	"public_path": "./src/github.com/elcct/taillachat/public",
	"template_path": "./src/github.com/elcct/taillachat/views",	
	"database": {
		"hosts": "localhost",
		"database": "defaultproject"
	}
}
```

Finally, you can run:

```
./bin/taillachat
```

That should output something like:

```
2014/06/19 15:31:15.386961 Starting Goji on [::]:8000
```

And it means you can now direct your browser to `localhost:8000`

## Icons

Tailla Chat uses Skype Emoticons via https://github.com/kof/emoticons/

Check license: https://github.com/kof/emoticons/blob/master/LICENSE

