function init() {
	// Connect to the server's WebSocket
    var serverSock = new WebSocket("ws://" + window.location.host + "/sock/");

    serverSock.onmessage = function(message) {
		console.log("Someone sent: ", message);
		var jsonMessage = JSON.parse(message);


		if(jsonMessage.Event == "chatMessage") {
			// Add the chat message to the output box
			var chatOutput = document.getElementById("chat_output");
			chatOutput.innerHTML += jsonMessage.Data.User + ": " + (jsonMessage.Data.message).replace(/[<>]/g, '') + "<br>";

			// Scroll to bottom of textbox
			chatOutput.scrollTop = chatOutput.scrollHeight;
		} else if(jsonMessage.Event == "screenUpdate") {

		}
	};



	var stage = new createjs.Stage("mainCanvas");

	var circle = new createjs.Shape();

	// To cache an object: DisplayObject.cache()

	circle.graphics.beginFill("DeepSkyBlue").drawCircle(0, 0, 25);
	circle.x = 400;
	circle.y = 250;
	stage.addChild(circle);

	stage.update();

	// The ammount the player should move every time a movement key is pressed
	var moveStep = 10;

	// Keypress listener
	var listener = new window.keypress.Listener();
	listener.register_many([
	{
	    "keys"       : "w",
	    "on_keydown" : function() {
            console.log("You pressed w");
            serverSock.send(JSON.stringify({
				Event: "w down"
			}));
        },
        "on_keyup"   : function(e) {
            console.log("And now you've released w.");
            serverSock.send(JSON.stringify({
				Event: "w up"
			}));
        }
	},
	{
		"keys"       : "a",
	    "on_keydown" : function() {
            console.log("You pressed a");
            serverSock.send(JSON.stringify({
				Event: "a down"
			}));
        },
        "on_keyup"   : function(e) {
            console.log("And now you've released a.");
            serverSock.send(JSON.stringify({
				Event: "a up"
			}));
        }
	},
	{
		"keys"       : "s",
	    "on_keydown" : function() {
            console.log("You pressed s");
            serverSock.send(JSON.stringify({
				Event: "s down"
			}));
        },
        "on_keyup"   : function(e) {
            console.log("And now you've released s.");
            serverSock.send(JSON.stringify({
				Event: "s up"
			}));
        }
	},
	{
		"keys"       : "d",
	    "on_keydown" : function() {
            console.log("You pressed d");
            serverSock.send(JSON.stringify({
				Event: "d down"
			}));
        },
        "on_keyup"   : function(e) {
            console.log("And now you've released d.");
            serverSock.send(JSON.stringify({
				Event: "d up"
			}));
        }
	}]);


	// Get the chat input box
	var chatInput = document.getElementById('chat_input');
	// Stop listening for keyboard events for the canvas when the chat box is focussed
	chatInput.addEventListener("focus", chatInputFocussed);
	function chatInputFocussed() {
		listener.stop_listening();
	} // end chatInputFocussed()
	// Start listening again when it loses focus
	chatInput.addEventListener("blur", chatInputFocusLost);
	function chatInputFocusLost() {
		listener.listen();
	} // end chatInputFocusLost()


	// Set the framerate for the ticker
	createjs.Ticker.setFPS(30);
	// Update stage will render next frame
    createjs.Ticker.addEventListener("tick", update);


    // Text chat input onkeydown event
	document.getElementById("chat_input").onkeydown = function(e) {
		// If the enter key is pressed
		if((e.keyCode || e.charCode) === 13) {
			// Get the input text
			var chatInputBox = document.getElementById("chat_input");

			// Send the chat message
			serverSock.send(JSON.stringify({
				Event: "chatMessage",
				Data: chatInputBox.value
			}));
			//.send(chatInputBox.value);

			// Add the chat message to the output box
			var chatOutput = document.getElementById("chat_output");
			chatOutput.innerHTML += "You: " + (chatInputBox.value).replace(/[<>]/g, "") + "<br>";

			// Scroll to bottom of textbox
			chatOutput.scrollTop = chatOutput.scrollHeight;

			// Reset the chat input box
			chatInputBox.value = "";
		} // end if
	};


    function update() {
		// Will cause the circle to wrap back
		if(circle.x > stage.canvas.width) { 
			circle.x = stage.canvas.width;
		} // end if
		if(circle.x < 0) {
			circle.x = 0;
		} // end if
		if(circle.y > stage.canvas.height) { 
			circle.y = stage.canvas.height;
		} // end if
		if(circle.y < 0) {
			circle.y = 0;
		} // end if

		stage.update();

	} // end update()

} // end init()

