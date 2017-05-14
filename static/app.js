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
var announce = function(hash) {
	$.ajax("/api/announce?hash=" + hash);
}

var loadData = function() {
	$.ajax("/api/stats").success(function(data) {
		//alert(data["UL"]);
		$("#totalUL").html(data["UL"].fileSize(1) + "/s");
		$("#totalDL").html(data["DL"].fileSize(1) + "/s");

	});
	if(tTable) {
		tTable.DataTable().ajax.reload(null, false);
		return;
	}
	var statsRender = function(data, type, row, meta) {
		return row["Done"].fileSize(1) + "/" + row["Size"].fileSize(1);
	}

	var nameRender = function(data, type, row, meta) {
		percent = Math.round((row['Done'] * 100) / row['Size']);
		//console.log(row["Activity"]);
		var data = row['Name'];

		data += '<div class="pull-right clear">';

		switch (row['Status']) {
			case "Seeding":
			data += '<div class="progress progress-badge label-badge pull-right"><span class="label  progress-center">Seeding: ' + row['Uploaded'].fileSize(1) + '</span><div role="progressbar" class="progress-bar progress-bar-success progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="' + percent + '" style="width:' + percent + '%"></div></div>';
				break
			case "Downloading":
				//data += '<span class="label label-badge label-primary">Downloading: ' + percent + '%</span> ';
				data += '<div class="progress progress-badge label-badge pull-right"><span class="label  progress-center">Downloading: ' + row['Done'].fileSize(1) + ' (' + percent + '%)</span><div role="progressbar" class="progress-bar progress-bar-primary progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="' + percent + '" style="width:' + percent + '%"></div></div>';
				break
			case "Stopped":
				data += '<div class="progress progress-badge label-badge pull-right"><span class="label  progress-center">Stopped: ' + row['Done'].fileSize(1) + ' (' + percent + '%)</span><div role="progressbar" class="progress-bar progress-bar-warning progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="' + percent + '" style="width:' + percent + '%"></div></div>';

				//data += '<span class="label label-badge label-warning">Stopped</span> ';
				break
			case "Error":
				data += '<div class="progress progress-badge label-badge pull-right"><span class="label  progress-center">Error</span><div role="progressbar" class="progress-bar progress-bar-danger progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="' + percent + '" style="width:' + percent + '%"></div></div>';
				break
			default:
			data += '<div class="progress progress-badge label-badge pull-right"><span class="label  progress-center">' + row["Status"] + '</span><div role="progressbar" class="progress-bar progress-bar-danger progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="' + percent + '" style="width:' + percent + '%"></div></div>';
		}
		data += '<span class=" label label-badge label-default label-silver"><span class="glyphicon glyphicon-cloud-upload"></span> ' + row['UL'].fileSize(1) + '/s</span> ';
		data += '<span class=" label label-badge label-default label-silver"><span class="glyphicon glyphicon-cloud-download"></span> ' + row['DL'].fileSize(1) + '/s</span> ';

		data += '<br style="clear: both;" />';
		data += "</div>";



		//data += '<span  onclick="announce(\'' + row["Hash"] + '\');" class="btn btn-danger btn-xs glyphicon glyphicon-refresh"></span> ';


		data += '</div>';

		return data;

	}
	tTable = $('#list').dataTable( {
		"ajax": "/api/list",
		"sAjaxDataProp": "",
		"sDom": '<"pull-left" i><"pull-right" l><t><"center" p>',
		"select": true,
		"rowId": 'Hash',
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
setInterval(loadData, 1000);
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
$("#tfDel").click(function() {
	var data = tTable.DataTable().rows( { selected: true } ).data();
	$(data).each(function(x) {
		$.ajax("/api/remove?hash=" + data[x]["Hash"])
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

$("#tfAdd").click(function() {
	$("#modalAdd").modal();
	$("#modalAdd .submit").click(function() {
		$("#modalAdd form").submit();
		$("#modalAdd").modal('hide');
	});

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
$("#tfRefresh").click(loadData);
//copypaste from http://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable
Object.defineProperty(Number.prototype,'fileSize',{value:function(a,b,c,d){
 return (a=a?[1e3,'k','B']:[1024,'K','iB'],b=Math,c=b.log,
 d=c(this)/c(a[0])|0,this/b.pow(a[0],d)).toFixed(2)
 +' '+(d?(a[1]+'MGTPEZY')[--d]+a[2]:'Bytes');
},writable:false,enumerable:false});
