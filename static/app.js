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
		return row["Done"].fileSize(1) + "/" + row["Size"].fileSize(1);
	}

	var nameRender = function(data, type, row, meta) {
		percent = Math.round((row['Done'] * 100) / row['Size']);

		var data = row['Name'];
		data += '<div class="pull-right clear">';
		data += '<span class="label label-badge label-default label-silver"><span class="glyphicon glyphicon-cloud-upload"></span> ' + row['UL'].fileSize(1) + '/s</span> ';
		data += '<span class="label label-badge label-default label-silver"><span class="glyphicon glyphicon-cloud-download"></span> ' + row['DL'].fileSize(1) + '/s</span> ';
		data += '<span class="label label-badge label-success">' + row['Seeds'] + ' seeds</span> ';
		data += '<span class="label label-badge label-primary">' + row['Peers'] + ' peers</span> ';
		data += '<span class="label label-badge label-warning">' + row['Done'].fileSize(1) + ' done</span> ';
data += '<span class="label label-badge label-warning">' + row['Uploaded'].fileSize(1) + ' uploaded</span> ';
		data += '</div><br style="clear: both;" />';

		data += '<div class="pull-left">';
		data += '</div>';
		data += '<div class="progress"><div role="progressbar" class="progress-bar progress-bar-success progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="' + percent + '" style="min-width: 32px; width:' + percent + '%">' + percent + '%</div></div>';

		return data;

	}
	tTable = $('#list').dataTable( {
		"ajax": "/api/list",
		"sAjaxDataProp": "",
		"sDom": 'lrtip',
		"select": true,
		"autoWidth": false,
		"columns.defaultContent": "",
		"columns": [
			{ "width": "100%", "data": "Name", render: nameRender },


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
		$.ajax("/api/start?hash=" + data[x]["Hash"])
	})

	return false;
});
$("#tfStop").click(function() {
	var data = tTable.DataTable().rows( { selected: true } ).data();
	$(data).each(function(x) {
		$.ajax("/api/stop?hash=" + data[x]["Hash"])
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

//copypaste from http://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable
Object.defineProperty(Number.prototype,'fileSize',{value:function(a,b,c,d){
 return (a=a?[1e3,'k','B']:[1024,'K','iB'],b=Math,c=b.log,
 d=c(this)/c(a[0])|0,this/b.pow(a[0],d)).toFixed(2)
 +' '+(d?(a[1]+'MGTPEZY')[--d]+a[2]:'Bytes');
},writable:false,enumerable:false});
