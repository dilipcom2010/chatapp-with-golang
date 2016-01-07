$(document).ready(function() {
	$(".usr-option").click(function(){
		$(".pic-popup").toggle();
	});
	$(function(){
		sizingcontainer();
	});
	$(window).resize(function() {
		sizingcontainer();
		//document.write($('body').width())
		var miw = $(".chat-panel").height();
		$(".receiving-panel").css("height", miw-120+"px")
		$(".sending-wrapper").css("width", $('.receiving-panel').innerWidth()-50+"px")
	});
	
});


function sizingcontainer()
{
	var iw = $('body').width();
	//var iw = Math.max( $(window).width(), window.innerWidth);
	var ih = $('body').height();
	//var ih = $('body').innerHeight();
	//document.write(iw)
	$(".chat-with").css("width", iw-300+"px")
	$(".users-panel").css("height", ih-50+"px")
	$(".message-panel").css("width", iw-300+"px")
	$(".message-panel").css("height", ih-50+"px")
	$(".option-panel").css("height", ih-50+"px")
	$("#all-user").css("max-height", ih-123+"px")
	$("#groups").css("max-height", ih-123+"px")
	$("#online").css("max-height", ih-123+"px")

	/*var miw = $(".chat-panel").height();
	$(".receiving-panel").css("height", miw-120+"px")
	$(".sending-wrapper").css("width", $('.receiving-panel').innerWidth()-50+"px")*/
}


function chatting(from, to, name, pic, dp)
{
	//document.write(dp)
	//document.write(name);
	//document.getElementById("chat-with").innerHTML = '<span class="chatting-with">'+name+'</spn>';
	
	document.getElementById("chat-with").innerHTML = '<div class="chatting"><div class="img-box"><img src="'+pic+'"></div><div class="name"><span class="with">'+name+'</span><br><span class="status">status</span></div></div>';

	//document.getElementById("sendmsg").setAttribute("to", to);

	var receiving_box_id = "chat_box_of_user"+to;
	var receiving_box = document.getElementById(receiving_box_id);
	
	//var divs = document.getElementsByClassName('chat-panel');
	//for(var i=0; i < divs.length; i++) 
	//{ 
	//	divs[i].style.display = 'none';
	//}
	
	$(".chat-panel").css("display", "none");

	if (receiving_box == null)
	{
		createChatPanel(from, to, dp)
		var log = $("#"+receiving_box_id+" .receiving-panel");
		loadmsg(receiving_box_id, log, from, to);
			
		var scrolling = log[0];
		scrolling.scrollTop = scrolling.scrollHeight - scrolling.clientHeight;
		/*var doScroll = scrolling.scrollTop == (scrolling.scrollHeight - scrolling.clientHeight);
		if(doScroll)
		{
			//document.write(scrolling.scrollTop);
			scrolling.scrollTop = 100;
			//$("#"+receiving_box_id+" .receiving-panel").scrollTop(100);	
			//document.write(scrolling.scrollTop);
		}*/

	}
	else
	{
		receiving_box.style.display="block";
		var rsvd_msg_count = "rsv_msg_cnt_of"+to;
		$("#"+rsvd_msg_count).text('');
	}


	/*$(function(){
	if (receiving_box == null)
	{
		var miw = $(".chat-panel").height();
		//var z = $('.sending-panel').width();
		//document.write(z);
		var receiving_panel='<div class="receiving-panel" id="'+receiving_box_id+'" style="height: '+(miw-120)+'px"></div><div class="sending-panel"><div class="separator"><div class="user-img"><img src="static/images/profile-pic/dilip.jpg"></div></div><div class="separator sending-wrapper"><div class="sending-box"><textarea id="sendmsg" from="'+from+'" to="'+to+'" onkeyup="Javascript: if (event.keyCode==13) sendmsg(this);"></textarea></div><div class="sending-labels"><button class="send-btn">SEND</button></div></div></div>';
     	           



    	var new_div = document.createElement('div');
    	new_div.innerHTML = receiving_panel
			$(new_div).addClass("chat-panel")
					.appendTo($(".message-panel")) //main div

	}
});*/
	//var miw = $(".chat-panel").height();
	//$(".receiving-panel").css("height", miw-120+"px")
	//$(".sending-wrapper").css("width", $('.sending-panel').width()-50+"px")
}

function createChatPanel(from, to, dp)
{
	//document.write(dp)
	var receiving_box_id = "chat_box_of_user"+to;
	var receiving_box = document.getElementById(receiving_box_id);

	var message_panel = document.getElementsByClassName("message-panel")[0]
	var chat_panel = document.createElement("div");
	var receiving_panel = document.createElement("div");
	//var loading_button = document.createElement("button");
	var sending_panel = document.createElement("div");
	var separator1 = document.createElement("div");
	var separator2 = document.createElement("div");
	var user_img = document.createElement("div");
	var img = document.createElement("img");
	var sending_box = document.createElement("div");
	var textarea = document.createElement("textarea");
	var sending_labels = document.createElement("div");
	var send_btn = document.createElement("button");


		
	$(chat_panel).addClass("chat-panel")
			.attr("id", receiving_box_id)
			.appendTo($(message_panel))
	$(receiving_panel).addClass("receiving-panel")
			.css("height", $(chat_panel).height()-120+"px")
			.attr("onscroll", "Javascript: if (this.scrollTop == 0) loadmsg('"+receiving_box_id+"', this, "+from+", "+to+");")
			.attr("offset", 0)
			.attr("load", 1)
			.appendTo($(chat_panel))



	//$(loading_button).addClass("load-btn")
	//		.appendTo($(receiving_panel))




	$(sending_panel).addClass("sending-panel")
			.appendTo($(chat_panel))
	$(separator1).addClass("separator")
			.appendTo($(sending_panel))
	$(user_img).addClass("user-img")
			.appendTo($(separator1))
	$(img).attr("src", dp)
			.appendTo($(user_img))
	$(separator2).addClass("separator")
			.addClass("sending-wrapper")
			.css("width", $(receiving_panel).width()-50+"px")
			.appendTo($(sending_panel))
	$(sending_box).addClass("sending-box")
			.appendTo($(separator2))
	$(textarea).addClass("sendmsg")
			.attr("to", to)
			.attr("from", from)
			.attr("onkeyup", "Javascript: if (event.keyCode==13) sendmsg(this);")
			.appendTo($(sending_box))
	$(sending_labels).addClass("sending-labels")
			.appendTo($(separator2))
	$(send_btn).addClass("send-btn")
			.text("SEND")
			.appendTo($(sending_labels))
}


function loadmsg(rbox_id, e, u1, u2)
{
	var offset = $(e).attr("offset");
	var limit = offset + 20;
	var load = $(e).attr("load");
	//document.write(offset)
	if(load == 1)
	{
		$.ajax({url: "http://127.0.0.1:8080/loadchat?user1="+u1+"&user2="+u2+"&offset="+offset+"&limit="+limit, success: function(result){
        	    $('<br>'+result+'</br>').prependTo("#"+rbox_id+" .receiving-panel");
            	$(e).attr("offset", limit);
            	if(result == "No more messages")
            	{
            		$(e).attr("load", 0);
            	}
    	}});
    }
}