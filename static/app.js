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

		const template = ({icon, state, status, color, percent, dl, ul, name, trackers}) => `
		${name}
		<div class="hidden">
			${trackers}
			${state}
		</div>
		<div class="pull-right clear">
			<div class="progress progress-badge label-badge pull-right">
				<span class="label progress-center"><span class="${icon}"></span> ${status} (${percent}%)</span>
				<div role="progressbar" class="progress-bar progress-bar-${color} progress-bar-bg" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${percent}" style="width:${percent}%"></div>
			</div>
			<span class="label label-badge label-default label-silver"><span class="glyphicon glyphicon-cloud-download"></span> ${dl}/s</span>
			<span class="label label-badge label-default label-silver"><span class="glyphicon glyphicon-cloud-upload"></span> ${ul}/s</span>

		</div>`;
		var trackerList = "";
		for(i in row["Trackers"]) {
			trackerList += row["Trackers"][i];
		}
		var data = {
			"percent": percent,
			"trackers": trackerList,
			"ul": row['UL'].fileSize(1),
			"dl": row['DL'].fileSize(1),
			"name": row["Name"],
			"state": row["Status"]
		}

		switch (row['Status']) {
			case "Seeding":
				data["status"] = row['Uploaded'].fileSize(1);
				data["color"] = "success";
				data["icon"] = "glyphicon glyphicon-leaf";
				break
			case "Downloading":
				data["status"] = row['Done'].fileSize(1);
				data["color"] = "primary";
				data["icon"] = "glyphicon glyphicon-flash";
				break
			case "Stopped":
				data["status"] = row['Done'].fileSize(1);
				data["color"] = "warning";
				data["icon"] = "glyphicon glyphicon-off";
				break
			case "Error":
				data["status"] = "Error";
				data["color"] = "danger";
				data["icon"] = "glyphicon glyphicon-warning-icon";
				break
			default:
				data["status"] = row['Status'];
				data["color"] = "warning";
		}



		return [data].map(template).join("");

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

var loadStats = function() {
	$.ajax("/api/stats").success(function(data) {
		//alert(data["UL"]);
		$("#totalUL").html(data["UL"].fileSize(1) + "/s");
		$("#totalDL").html(data["DL"].fileSize(1) + "/s");
		$("#totalPeers").html(data["Peers"]);
		$("#totalSeeds").html(data["Seeds"]);
		$("#totalTrackers").html(data["TrackersNo"]);

		var trackerListData = "";
		var statesListData = "";

		const labelTpl = ({label, count}) => `
			<li class="list-group-item">${label} <span class="badge pull-right">${count}</span></li>
		`;

		for(i in data["TrackersMap"]) {
			trackerListData += [{count: data["TrackersMap"][i], label: i }].map(labelTpl).join("");
		}
		for(i in data["States"]) {
			statesListData += [{count: data["States"][i], label: i}].map(labelTpl).join("");
		}
		$("#trackerList").html(trackerListData);
		$("#statesList").html(statesListData);
	});

}
setInterval(loadData, 1000);
setInterval(loadStats, 1000);
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
