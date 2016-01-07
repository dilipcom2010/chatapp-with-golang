if (window["WebSocket"]) 
{
	conn = new WebSocket("ws://127.0.0.1:8080/chat");

	conn.onclose = function(evt) {
	//appendLog($("<div><b>Connection closed.</b></div>"))
	document.getElementById("MsgArea").innerHTML="bye"
}
conn.onmessage = function(evt) {
//appendLog($("<div/>").text(evt.data))
//document.getElementById("MsgArea").innerHTML= evt.data
var message = evt.data
var received = message.split(",");
//var received = JSON.parse(message);
//document.write(typeOf(message))
//document.write(received[0]+"<br>"+received[1]+"<br>"+received.slice(2))
var receiving_box_id = "chat_box_of_user"+received[1];
var temp = document.getElementById(receiving_box_id).innerHTML
document.getElementById(receiving_box_id).innerHTML = temp+"<br>"+received.slice(2)
}

} 
else 
{
	appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
}

function sendmsg()
{
	//document.write("socketttttt");
	if (!conn)
	{
		document.write(conn.send)
	}
	var msg = document.getElementById("sendmsg").value
	//var to = document.getElementById("sendmsg").getAttribute("to")
	var letter = document.getElementById("sendmsg")
	var to = letter.getAttribute("to")
	var from = letter.getAttribute("from")
	//document.write(to, from)
	document.getElementById("sendmsg").value=''
	var message = [];
	message[0] = to
	message[1] = from
	message[2] = msg
	//document.write(message[0]+"<br>"+message[1]+"<br>"+message[2])
	conn.send(message)
}