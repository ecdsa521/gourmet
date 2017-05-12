var tTable;
$(document).ready(function() {

	loadData();
} );
$("#search").on("input propertychange paste", function() {
	if(tTable) {

		tTable.fnFilter( $(this).val() );
		tTable.fnDraw();
	}
	return false;
})
$.ajaxSetup({
    cache: false
});

var loadData = function() {
	if(tTable) {
		tTable.DataTable().ajax.reload(null, false);
		return;
	}
	tTable = $('#list').dataTable( {
		"ajax": "/api/list",
		"sAjaxDataProp": "",
		"sDom": 'lrtip',

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
