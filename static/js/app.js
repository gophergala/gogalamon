function init() {
	// Connect to the server's WebSocket
    var serverSock = new WebSocket("ws://" + window.location.host + "/sock/");

    serverSock.onmessage = function(message) {
		
		var jsonMessage = JSON.parse(message.data);


		if(jsonMessage.Event == "chatMessage") {
			// Add the chat message to the output box
			var chatOutput = document.getElementById("chat_output");
			chatOutput.innerHTML += jsonMessage.Data.User + ": " + (jsonMessage.Data.Message).replace(/[<>]/g, '') + "<br>";

			// Scroll to bottom of textbox
			chatOutput.scrollTop = chatOutput.scrollHeight;
		} else if(jsonMessage.Event == "screenUpdate") {
			viewCenter.x = jsonMessage.Data.ViewX;
			viewCenter.y = jsonMessage.Data.ViewY;

			onScreenObjects = jsonMessage.Data.Objs;
		} // end if/else if

	};


	// Init the stage
	var stage = new createjs.Stage("mainCanvas");
	// Set the stage to clear before each render
	stage.autoClear = true;

	// Init the location, in map space, of the center (and therefor our player) of our view
	var viewCenter = {
		x : null,
		y : null
	}

	// Init the array holding all the objects we're going to render on the screen
	var onScreenObjects = [];


	// Keypress listener
	var listener = new window.keypress.Listener();

	listener.register_many([
		{
		    "keys"       : "w",
		    "on_keydown" : function() {
	            serverSock.send(JSON.stringify({
					Event: "w down"
				}));
	        },
	        "on_keyup"   : function(e) {
	            serverSock.send(JSON.stringify({
					Event: "w up"
				}));
	        }
		},
		{
			"keys"       : "a",
		    "on_keydown" : function() {
	            serverSock.send(JSON.stringify({
					Event: "a down"
				}));
	        },
	        "on_keyup"   : function(e) {
	            serverSock.send(JSON.stringify({
					Event: "a up"
				}));
	        }
		},
		{
			"keys"       : "s",
		    "on_keydown" : function() {
	            serverSock.send(JSON.stringify({
					Event: "s down"
				}));
	        },
	        "on_keyup"   : function(e) {
	            serverSock.send(JSON.stringify({
					Event: "s up"
				}));
	        }
		},
		{
			"keys"       : "d",
		    "on_keydown" : function() {
	            serverSock.send(JSON.stringify({
					Event: "d down"
				}));
	        },
	        "on_keyup"   : function(e) {
	            serverSock.send(JSON.stringify({
					Event: "d up"
				}));
	    	}
	    },
	    {
			"keys"       : "f",
		    "on_keydown" : function() {
	            serverSock.send(JSON.stringify({
					Event: "f down"
				}));
	        },
	        "on_keyup"   : function(e) {
	            serverSock.send(JSON.stringify({
					Event: "f up"
				}));
	    	}
		}
	]);


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
				Message : chatInputBox.value
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

	// Sweet jesus the normal prompts are ugly
	var playerName = prompt("Please enter your player name");
	serverSock.send(JSON.stringify({
		Event   : "username",
		User    : playerName
	}));


	// Var to hold IDs of the backgrounds
	//var backgroundIds = [];


    function update() {
    	// To cache an object: DisplayObject.cache()

    	// Get the mainCanvas
    	var mainCanvas = document.getElementById("mainCanvas");



    	removeOldChildren(stage.children, onScreenObjects);



		// Place the far starfield
    	for (var i = mod(viewCenter.x * -0.1) - 512; i < mainCanvas.width; i += 512) {
    		for (var j = mod(viewCenter.y * -0.1) - 512; j < mainCanvas.height; j += 512) {
    			var starFieldFar = new createjs.Bitmap("img/starfield_far.png");
				starFieldFar.x = i;
				starFieldFar.y = j;

				//starFieldFar.name = backgroundIds.length;
				//onScreenObjects.push(starFieldFar);

				stage.addChild(starFieldFar);
    		};
    	};

    	// Place the mid starfield
    	for (var i = mod(viewCenter.x * -0.9) - 512; i < mainCanvas.width; i += 512) {
    		for (var j = mod(viewCenter.y * -0.9) - 512; j < mainCanvas.height; j += 512) {
    			var starFieldNear = new createjs.Bitmap("img/starfield_near.png");
				starFieldNear.x = i;
				starFieldNear.y = j;

				//starFieldNear.name = backgroundIds.length;
				//onScreenObjects.push(starFieldNear);

				stage.addChild(starFieldNear);
    		};
    	};

    	// Place the near starfield
    	for (var i = mod(viewCenter.x * -0.9) - 512; i < mainCanvas.width; i += 512) {
    		for (var j = mod(viewCenter.y * -0.9) - 512; j < mainCanvas.height; j += 512) {
    			var starFieldMid = new createjs.Bitmap("img/starfield_middle.png");
				starFieldMid.x = i;
				starFieldMid.y = j;

				//starFieldMid.name = backgroundIds.length;
				//onScreenObjects.push(starFieldMid);

				stage.addChild(starFieldMid);
    		};
    	};

    	// Sort the IDs for the background
    	//backgroundIds.sort();

    	// Create and place each new object we're sent
    	for(var i = 0; i < onScreenObjects.length; i++) {
    		// Get the object we want to render
    		var currentObject = onScreenObjects[i];

    		var objectBitmap;
    		var addChildBool;

    		if(stage.children.indexOf(currentObject) == -1) {
    			// Create the bitmap object
    			objectBitmap = new createjs.Bitmap("img/" + currentObject.N + ".png");

    			// Set the bitmap name to its unique id
    			objectBitmap.name = currentObject.I;

    			// Set the middle of the image (using the ship's size)
	    		objectBitmap.regX = 32;
	    		objectBitmap.regY = 32;

	    		addChildBool = true;
    		} else {
    			objectBitmap = stage.children[stage.children.indexOf(currentObject)];

    			addChildBool = false;
    		} // end if/else

    		objectBitmap.x = Math.round(currentObject.X - viewCenter.x + mainCanvas.width/2 + 32);
    		objectBitmap.y = Math.round(currentObject.Y - viewCenter.y + mainCanvas.height/2 + 32);
    		objectBitmap.rotation = currentObject.R;

    		// If the object is already on the stage, don't add it
    		if(addChildBool) {
    			stage.addChild(objectBitmap);
    		} // end if
    	} // end for

    	console.log(stage.numChildren);

		stage.update();
	} // end update()

	function removeOldChildren(currentObjects, newObjects) {
		for(var i = 0; i < currentObjects.length; i++) {
			if(newObjects.indexOf(currentObjects[i]) == -1) {
				var childToRemove = stage.getChildByName(currentObjects[i].I);

				stage.removeChild(childToRemove);
			} // end if
		} // end for

	} // end removeOldChildren()

	function mod(z) {
		z = z % 512;

		if(z < 0) {
			z += 512;
		} // end if

		return z;
	} // end mod()

} // end init()

