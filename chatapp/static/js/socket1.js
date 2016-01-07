$(document).ready(function() {
	//setTimeout(function(){
	createsocket();
	//}, 15000);
});


function createsocket()
{
	if (window["WebSocket"]) 
	{
		conn = new WebSocket("ws://127.0.0.1:8080/chat");

		conn.onclose = function(evt) {
			//appendLog($("<div><b>Connection closed.</b></div>"))
			document.getElementById("MsgArea").innerHTML="bye"
		}
		var snd = new Audio("/static/sound/notify2.mp3");
		conn.onmessage = function(evt) {
			//appendLog($("<div/>").text(evt.data))
			//document.getElementById("MsgArea").innerHTML= evt.data
			var message = evt.data
			var received = message.split(",");
			if(received[2] == "-1")
			{
				var ol9_usr = "online"+received[1];
				document.getElementById(ol9_usr).innerHTML = ol9_usr;
			}
			else if(received[2] == "-2")
			{
				document.write("hii")
			}
			else
			{
				//var received = JSON.parse(message);
				//document.write(typeOf(message))
				//document.write(received[0]+"<br>"+received[1]+"<br>"+received.slice(2))
				var receiving_box_id = "chat_box_of_user"+received[1];
				var rsvd_msg_count = "rsv_msg_cnt_of"+received[1];
				//var temp = document.getElementById(receiving_box_id).innerHTML
				//$("<div/>").text(received.slice(2)).appendTo("#"+receiving_box_id+" .receiving-panel");
				//document.getElementById(receiving_box_id).innerHTML = received.slice(2)
				//$("#"+receiving_box_id+" .receiving-panel").text("keep it up")
				//var xxx = $("#"+receiving_box_id+" .receiving-panel");
				//$("kkeepp iitt uupp").appendTo($(xxx))
				var cc=received.slice(3)
				//<div class="msg-box"><div class="bubble-right">www.css3-generator.weebly.com<span class="time">jjjj</span></div></div>
				//document.write(cc)
				//var msg = $('<div class="msg-box"><div class="bubble-left">'+cc+'&nbsp;&nbsp<span class="time">'+received[2]+'</span></div></div>');
	

		
				//document.write(received[0], received[1]);
				var receiving_box = document.getElementById(receiving_box_id);
				
				if(receiving_box == null)
				{
					//document.write("hii");
					//document.write(received[0], received[1]);
					createChatPanel(received[0], received[1], "{{$.dp}}");
					//receiving_box.style.display = none;
					$("#"+receiving_box_id).css("display", "none");	
				}
				if($("#"+receiving_box_id).css("display") == 'none')
				{
					var msg_count = $("#"+rsvd_msg_count).text();
					if(msg_count == '')
					{
						msg_count = "0";
					}
					$("#"+rsvd_msg_count).text(parseInt(msg_count)+1);
				}

				



				var log = $("#"+receiving_box_id+" .receiving-panel");
				var scrolling = log[0];
				var doScroll = scrolling.scrollTop == (scrolling.scrollHeight - scrolling.clientHeight);
				$('<div class="msg-box"><div class="bubble-left">'+cc+'&nbsp;&nbsp<span class="time">'+received[2]+'</span></div></div>').appendTo("#"+receiving_box_id+" .receiving-panel");
				snd.play();






				if(doScroll){
					scrolling.scrollTop = scrolling.scrollHeight - scrolling.clientHeight;
				}
				else
				{
					var msg_count = $("#"+rsvd_msg_count).text();
					if(msg_count == '')
					{
						msg_count = "0";
					}
					$("#"+rsvd_msg_count).text(parseInt(msg_count)+1);
				}
				//$('<div class="user-img"><img src="static/images/profile-pic/dilip.jpg"></div><div class="received">'+cc+'</div>').appendTo("#"+receiving_box_id+" .receiving-panel");
			}
		}

	} 
	else 
	{
		appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
	}
}



function sendmsg(e)
{
	//document.write("socketttttt");
	if (!conn)
	{
		document.write(conn.send)
	}
	var msg = e.value
	//var to = document.getElementById("sendmsg").getAttribute("to")
	//var letter = document.getElementById("sendmsg")
	var to = e.getAttribute("to")
	var from = e.getAttribute("from")
	//document.write(to, from)
	e.value=''
	var message = [];
	message[0] = to
	message[1] = from
	message[2] = new Date().toString("dd/MM/yy hh:mm tt")
	message[3] = msg
	var receiving_box_id = "chat_box_of_user"+to;
	//document.write(message[0]+"<br>"+message[1]+"<br>"+message[2])
	//document.getElementById("xxx").innerHTML=message;
	//document.write(message);
	var log = $("#"+receiving_box_id+" .receiving-panel");
	var scrolling = log[0];
	var doScroll = scrolling.scrollTop == (scrolling.scrollHeight - scrolling.clientHeight);
	$('<div class="msg-box"><div class="bubble-right">'+message[3]+'&nbsp;&nbsp<span class="time">'+message[2]+'</span></div></div>').appendTo("#"+receiving_box_id+" .receiving-panel");
	//if(doScroll){
		scrolling.scrollTop = scrolling.scrollHeight - scrolling.clientHeight;
	//}
	conn.send(message)
}