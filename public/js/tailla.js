/* 0.0.1 (C)2013 Tailla.com */
function Tailla(socket_url) {
	this.messages = [];
	this.socket = new SockJS(socket_url);	
	this.audioFormat = null;
	this.audioEnabled = true;

	this.closeConfirmation = false;

	this.initUi();
	this.initAudio();
	this.initSocketEvents();
	this.initUiEvents();
}

Tailla.prototype.initUi = function() {
	this.regions = [
		'All UK',
		'London',
 		'Channel Islands',
		'East Anglia (Eastern England)',
		'East Midlands',		
		'Northern Ireland',
		'North East',
		'North West',
		'Scotland',
		'South East',
		'South West',
		'Wales',
		'West Midlands',
		'Yorkshire & Humberside',		
	];

	$(this.regions).each(function(k, v) {		
		$('#regions').append($('<option></option>').attr('value', v).text(v));
	});	

	this.field = $('#field');
	this.send = $('#send');
	this.start = $('#start-chat');
	this.content = $('#content');
	this.disconnect = $('#disconnect');
	this.restart = $('#restart-chat');
	this.picture = $('#picture');
	this.sound = $('#sound');	
	this.chat = $('#chat');
	this.bsod = $('#diedtodeath');

	this.title = document.title;
	this.alertId = null;
}

Tailla.prototype.insertTerminated = function() {
	this.messages.push('<tr><td class="col-md-1 text-muted"><small>' + moment().format('hh:mm A') + '</small></td><td class="col-md-1"></td><td><span class="label label-danger">Conversation completed.</span></td></tr>');
	var html = '';
	for(var i=0; i<this.messages.length; i++) {			
		html += this.messages[i];
	}

	var table = $('#messages');

	table.html(html);		
}

Tailla.prototype.insertMessage = function(who, what, self) {		
	var s = self ? 'label-primary' : 'label-success';
	var r = self ? '' : ' class="success"';
	what = $.emoticons.replace(what);
	this.messages.push('<tr' + r + '><td class="col-md-1 text-muted"><small>' + moment().format('hh:mm A') + '</small></td><td class="col-md-1"><span class="label ' + s + '">' + who + '</span></td><td>' + what + '</td></tr>');
	var html = '';
	for(var i=0; i<this.messages.length; i++) {			
		html += this.messages[i];
	}

	var table = $('#messages');
	table.html(html);	
}

Tailla.prototype.insertPicture = function(who, picture, self) {
		var _picture = '<img class="img-responsive" src="' + picture + '"/>';
		this.insertMessage(who, _picture, self);	
}

Tailla.prototype.insertTyping = function(flag) {
	var table = $('#actions');
	if(flag == "true") {
		var html = '<tr class="success"><td class="col-md-1"></td><td class="col-md-1"><span class="label label-success">Alien</span></td><td><img src="/assets/img/cursor.gif"/></td></tr>';
		table.append(html);
	} else {
		table.html('');
	} 
}

Tailla.prototype.insertPictureSending= function(flag) {
	var table = $('#actions');
	if(flag == "true") {
		var html = '<tr class="success"><td class="col-md-1"></td><td class="col-md-1"><span class="label label-success">Alien</span></td><td>is sending a picture <img src="/images/cursor.gif"/></td></tr>';
		table.append(html);
	} else {
		table.html(''); 
	}
}

Tailla.prototype.alertUser = function(message) {
	this.alertMessage = message;

	if(!this.alertId) {
		this.alertId = setInterval(function() {
			document.title = document.title == this.alertMessage ? this.title : this.alertMessage;
		}.bind(this), 1000);
		window.onmousemove = function() {
			clearInterval(this.alertId);
			document.title = this.title;
			this.alertId = window.onmousemove = null;			
		}.bind(this);
	}    
}

Tailla.prototype.terminate = function() {
	this.insertTerminated();

	$('#chat-speaker').slideUp('fast', function() {
		$('#chat-restart').slideDown();	  		
	});

	this.closeConfirmation = false;
	this.playCompleted();
}

Tailla.prototype.initSocketEvents = function() {

	this.socket.onopen = function() {
		//console.log("Connection Opened");
	}

	this.socket.onmessage = function(e) {		
		var message = JSON.parse(e.data)

		switch(message.event) {
			case 'online':
				$('#connection-status').html(message.r + ' ready to chat / ' + message.c + ' currently chatting');
				break;
			case 'join':				
				$('#waiting').slideUp('fast', function() {			
					$('#chat').slideDown();
				});	
				break;				
			case 'message':
					if(message.data) {            
						this.insertMessage('Alien', this.replaceURLs(message.data), false);
						this.playMessage();
						this.alertUser('New message!');
					} else {
						console.log("Error:", message);
					}
				break;
			case 'picturebefore':
				this.insertPictureSending(message.data);
				break;
			case 'picture': 
				this.insertPictureSending(false);
				this.insertPicture('Alien', message.data, false);
				this.playMessage();
				this.alertUser('New message!');			
				break;
			case 'typing':
				this.insertTyping(message.data)
				break;
			case 'exit':
				this.terminate();
				break;
		}

	}.bind(this);

	this.socket.onclose = function() {
		this.terminate();
		$('#welcome').hide();
		$('#chat').hide();
		$('#chat-restart').hide();
		$('#chat-speaker').hide();
		$('#waiting').hide();
		$('#messages').html(''); 
		$('#diedtodeath').show();		
	}.bind(this);
}

Tailla.prototype.emit = function(e, data) {
	data['event'] = e;
	var msg = JSON.stringify(data);	
	this.socket.send(msg)
}

Tailla.prototype.sendReady = function(region) {
	this.emit('ready', { region: region });	
}

Tailla.prototype.sendMessage = function(message) {
	this.emit('send', { message: message });	
}

Tailla.prototype.sendPicture = function(data) {
	this.emit('picture', { data: data });
}

Tailla.prototype.sendIsTyping = function(flag) {
	this.emit('typing', { typing: flag });
}

Tailla.prototype.sendExit = function() {
	this.emit('exit', { });	
}

Tailla.prototype.initUiEvents = function() {
	this.isTyping = null;
	this.doneTyping = null;

	$(window).on('beforeunload', function() {
		if(this.closeConfirmation) {
			return 'Are you sure you want to leave?';
		}
	}.bind(this));

	this.sound.click(function() {
		if(this.audioEnabled) {
			this.sound.html('<span class="glyphicon glyphicon-volume-off"></span>');
		} else {
			this.sound.html('<span class="glyphicon glyphicon-volume-up"></span>');
		}
		this.audioEnabled = !this.audioEnabled;

	}.bind(this));

	this.start.click(function() {		
		var region = $('#regions').val();

		this.sendReady(region);

		$('#welcome').hide();
		$('#waiting').show();

		this.closeConfirmation = true;
	}.bind(this));

	this.restart.click(function() {		
		var region = $('#regions').val();
		
		this.sendReady(region);
		
		$('#chat').hide();
	  	$('#chat-restart').hide();
		$('#chat-speaker').show();
		$('#waiting').show();
		$('#messages').html('');
		this.messages = [];		
	}.bind(this));	
 
    this.send.click(function() {    	
        var text = $.trim(this.field.val());
        if(text == '') return;

        screen_text = this.replaceURLs(sanitize(text).entityEncode());

        this.insertMessage('You', screen_text, true);
        this.sendMessage(text);
        this.field.val('');
    }.bind(this));

	this.field.keypress(function( event ) {				
		if(event.which != 13 ) {
			if(!this.isTyping) {			
				this.isTyping = true;
				this.sendIsTyping(true);
			}

			if(this.doneTyping) {
				window.clearTimeout(this.doneTyping);
			}

			this.doneTyping = window.setTimeout(function() {
				this.doneTyping = null;
				this.isTyping = null;
				this.sendIsTyping(false);
			}.bind(this), 1000);
		}

  		if(event.which == 13 ) {
			if(this.doneTyping) {
				window.clearTimeout(this.doneTyping);
				this.sendIsTyping(false);
				this.isTyping = null;
			}
  			this.send.click();
     		event.preventDefault();
  		}
  	}.bind(this));

  	this.picture.change(function(e) {  		
		var file = e.originalEvent.target.files[0];
		var reader = new FileReader();
		reader.onload = function(_e) {	
			var picture = _e.target.result;		
			this.insertPicture('You', picture, true);				        
	        this.sendPicture(picture);
    	}.bind(this);
    	reader.readAsDataURL(file);
  	}.bind(this));

  	this.disconnect.click(function() {  		
  		this.sendExit();
  		this.insertTerminated();
  		$('#chat-speaker').slideUp('fast', function() {
  			$('#chat-restart').slideDown();	
  		});
  		this.closeConfirmation = false;
  	}.bind(this));	
}

Tailla.prototype.play = function(filename) {
	if(this.audioFormat && this.audioEnabled) {
		new Audio('/assets/audio/' + filename + '.' + this.audioFormat).play();	
	}
}

Tailla.prototype.playMessage = function() {
	this.play('message');
}

Tailla.prototype.playCompleted = function() {
	this.play('completed');
}

Tailla.prototype.initAudio = function() {
	var format = null;

    var audio = document.createElement('audio');
    var ogg = !!(audio.canPlayType && audio.canPlayType('audio/ogg; codecs="vorbis"').replace(/no/, ''));
    if (ogg) format = 'ogg';
    var mp3 = !!(audio.canPlayType && audio.canPlayType('audio/mpeg;').replace(/no/, ''));
    if (mp3) format = 'mp3';

    this.audioFormat = format;
}

Tailla.prototype.replaceURLs = function(text) {
    var exp = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
    return text.replace(exp, '<a href="$1" target="_blank">$1</a>'); 
}