$(document).ready(function() {

	loadData();
} );
$.ajaxSetup({
    cache: false
});

var loadData = function() {
	$('#list').DataTable( {
		"ajax": "/api/list",
		"sAjaxDataProp": "",
		"destroy": true,
		"columns": [
			{ "data": "Name" },
			{ "data": "Path" },
			{ "data": "Hash" },
			{ "data": "Size" },
			{ "data": "Done" },
			{ "data": "Seeds" },
			{ "data": "Peers" }
		]
	} );

}
