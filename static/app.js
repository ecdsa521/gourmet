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
		"ajax": "/api/list?v",
		"sAjaxDataProp": "",
		"sDom": 'lrtip',
		"select": true,
		"autoWidth": false,
		"columns.defaultContent": "",
		"columns": [
			{ "width": "85%", "data": "Name", render: nameRender },
			{ "width": "10%", "data": "Stats", render: statsRender },


		],
		"rowCallback": function( row, data, index ) {
			$(row).attr("data-hash", data["Hash"]);
			//$(row).click(getDetails);
		}

	} );


}
//setInterval(loadData, 1000);
var getDetails = function(e) {
	if($(this).hasClass("selected")) {
		$("#list tr").removeClass("selected");
	} else {
		$("#list tr").removeClass("selected");
		$(this).addClass("selected");
	}


}
$("#tfStart").click(function() {
	var data = tTable.DataTable().rows( { selected: true } ).data();
	$(data).each(function(x) {
		$.ajax("/api/start/" + data[x]["Hash"])
	})

	return false;
});

$("#tfMagnet").click(function() {
	$("#modalMagnet").modal();
	$("#modalMagnet .submit").click(function() {
		$("#modalMagnet form").submit();
		$("#modalMagnet").modal('hide');
	});

	return false;
});
