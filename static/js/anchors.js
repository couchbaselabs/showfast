var hash = window.location.hash.substring(1);
if (hash) {
	$("#" + hash).button('toggle');
	$("div[id*='cont_" + hash + "']").css("visibility", "visible");
} else {
	$(".metric").css("visibility", "visible");
}
