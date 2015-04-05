$(document).ready(function() {
    $(".navbar-toggle").click(function() {
		$(".navbar-toggle .icon-bar-box").toggleClass("rotate-icon-bar");
		$(".show-users").css({"transition":"width 0.7s"});
		$(".container-user").toggleClass("show-users");
		$(".container-chat").toggleClass("expand-container-chat");
	});
	
	
	
	
	$(function(){
		$(".container-user").css("height", window.innerHeight-100+"px");
		$(".container-chat").css("height", window.innerHeight-100+"px");
		$(".container-options").css("height", window.innerHeight-100+"px");
		$(window).resize(function() {
			$(".container-user").css("height", window.innerHeight-100+"px");
			$(".container-chat").css("height", window.innerHeight-100+"px");
			$(".container-options").css("height", window.innerHeight-100+"px");
		});
	});


	$(function(){
		var h = $(".container-chat").height();
		$(".msg-receiving-area").height(h-140);
		$(".msg-sending-area").height(120);
		var p = $(".msg-area").width();
		$(".msg-area textarea").width(p-56);
		//$(".container-options").css("height", window.innerHeight-100+"px");
		$(window).resize(function() {
			var h = $(".container-chat").height();
			$(".msg-receiving-area").height(h-140);
			$(".msg-sending-area").height(120);
			var p = $(".msg-area").width();
			$(".msg-area textarea").width(p-56);
		});
	});
	
});


function show_hide_user(x, y, z)
{
	document.getElementById(x).style.display='block';
	document.getElementById(y).style.display="none";
	document.getElementById(z).style.display="none";
}
function start_chat(id)
{
	//document.write(id);
	document.getElementById("sendmsg").setAttribute("to", id);
	var receiving_box_id = "chat_box_of_user"+id;
	var receiving_box = document.getElementById(receiving_box_id);
	var divs = document.getElementsByClassName('msg-receiving-box');
	for(var i=0; i < divs.length; i++) 
	{ 
		divs[i].style.display = 'none';
	}

	if (receiving_box == null)
	{
		//document.write("");
		var new_div =document.createElement('div');
		//$("<div></div>").attr('id','new').appendTo('body');
		//$(".msg-receiving-box").height($(".msg-receiving-area").height());
		$(new_div).addClass("msg-receiving-box")
				.attr('id',receiving_box_id).height($(".msg-receiving-area").height())
				.appendTo($("#MsgArea")) //main div
			//.delay(2500)
	}
	else
	{
		receiving_box.style.display = 'block';
	}
}


//wo kahte h na love me enough to let me go....yaar main utna payar nhi kar sakta