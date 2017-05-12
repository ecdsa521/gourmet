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
	var statsRender = function(data, type, row, meta) {
		return row["Done"] + "/" + row["Size"];
	}
	var nameRender = function(data, type, row, meta) {
		percent = Math.round((row['Done'] * 100) / row['Size']);
		return row['Name'] + '<div class="progress"><div role="progressbar" class="progress-bar progress-bar-success progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="' + percent + '" style="min-width: 32px; width:' + percent + '%">' + percent + '%</div></div>';
	}
	tTable = $('#list').dataTable( {
		"ajax": "/api/list",
		"sAjaxDataProp": "",
		"sDom": 'lrtip',
		"autoWidth": false,
		"columns.defaultContent": "",
		"columns": [
			{ "width": "85%", "data": "Name", render: nameRender },
			{ "width": "15%", "data": "Stats", render: statsRender }

		]

	} );


}
