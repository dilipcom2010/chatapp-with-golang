$(document).ready(function() {
	$(function(){
		sizingcontainer();
	});
	$(window).resize(function() {
		sizingcontainer();
	});
});


function sizingcontainer()
{
	var ih = $('body').height();
	$(".users").css("height", ih-50+"px")
	$(".messages").css("height", ih-50+"px")
	$(".options").css("height", ih-50+"px")
}

