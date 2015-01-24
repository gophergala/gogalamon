

function init() {
	var stage = new createjs.Stage("mainCanvas");

	var circle = new createjs.Shape();

	circle.graphics.beginFill("DeepSkyBlue").drawCircle(0, 0, 50);
	circle.x = 400;
	circle.y = 250;
	stage.addChild(circle);


	stage.update();

	var moveStep = 10;

	// Keypress listener
	var listener = new window.keypress.Listener();

	listener.simple_combo("d", function() {
		circle.x += moveStep;
	});
	listener.simple_combo("a", function() {
		circle.x -= moveStep;
	});
	listener.simple_combo("w", function() {
		circle.y -= moveStep;
	});
	listener.simple_combo("s", function() {
		circle.y += moveStep;
	});


	createjs.Ticker.setFPS(30);
	//Update stage will render next frame
    createjs.Ticker.addEventListener("tick", update);


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

