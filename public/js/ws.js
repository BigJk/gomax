var socket = new WebSocket("ws://127.0.0.1:8080/ws");

var metricData = {
	rchart: [],
	fchart: []
};

metricData.rchart = getEmptyArray(dataPoints);
metricData.fchart = getEmptyArray(dataPoints);

socket.onmessage = function(event) {
	data = JSON.parse(event.data);
	switch (data.Type) {
		case 'Metrics':
			addMetrics(data.Msg);
			break;
		case 'Phases':
			for (var i = 0; i < phases; i++) {
				updatePhase(i.toString(), data.Msg[i]);
			}
			break;
	}
};

socket.onclose = function() {
	$('#online').removeClass('blink');
	$('#online').removeClass('fa-angle-double-up');
	$('#online').addClass('fa-angle-double-down');

	swal("Stopped", "gomax was successfully stopped", "success");
};

function addMetrics(d) {
	$('#rps').text(d.Rounds);
	$('#fps').text(d.Fights);

	metricData.rchart.push(d.Rounds);
	metricData.rchart = metricData.rchart.slice(1);

	metricData.fchart.push(d.Fights);
	metricData.fchart = metricData.fchart.slice(1);

	updateValues(rchart, metricData.rchart);
	updateValues(fchart, metricData.fchart);
}

function updateValues(c, d) {
	for (var i = 0; i < d.length; i++) {
		c.datasets[0].points[i].value = d[i];
	}
	c.update();
}

function exit() {
	if (socket.readyState == socket.CLOSED) {
		sweetAlert("Error!", "Program is already stopped!", "error");
		return;
	}
	swal({
		title: "Stop the program?",
		type: "info",
		showCancelButton: true,
		closeOnConfirm: false,
		showLoaderOnConfirm: true,
	}, function() {
		call('/api/stop');
	});
}

function reloadCharts() {
	for (var i = 0; i < phases; i++) {
		reloadChart(i);
	}
}

var chartData = {};

function reloadChart(phase) {
	charts[phase].segments[0].value = chartData[phase].Failed;
	charts[phase].segments[1].value = chartData[phase].Passed;
	charts[phase].update();
}

function updatePhase(phase, d) {
	chartData[phase] = {
		Passed: d.Passed,
		Failed: d.Failed
	};

	$('#p' + phase + 'total').text(d.Total);
	$('#p' + phase + 'failed').text(d.Failed);
	$('#p' + phase + 'passed').text(d.Passed);
	$('#p' + phase + 'bs').text(d.Bestscore.toFixed(2));

	for (var i = 0; i < d.Top.length; i++) {
		if ($('#p' + phase + i)[0] === null) {
			document.getElementById('p' + phase + 'table').innerHTML += '<tr id="p' + phase + i + '"><td>' + i + '</td><td id="p' + phase + i + 'id"></td><td id="p' + phase + i + 'res"></td><td id="p' + phase + i + 'score">23.48936</td><td id="p' + phase + i + 'passed" class="text-center"><i class="fa fa-times"></i> </td><td id="p' + phase + i + 'show" class="show"><a href="/api/warrior/" target="_blank"><i class="fa fa-search"></i></a></td></tr>';
		}

		var id = toID(d.Top[i].Combination);

		if ($('#p' + phase + i + 'id').text() == id) {
			$('#p' + phase + i).css("transition", "0.5s");
			$('#p' + phase + i).css("background", "rgba(000, 000, 000, 0)");
			continue;
		} else {
			$('#p' + phase + i).css("transition", "0.5s");
			$('#p' + phase + i).css("background", "rgba(51, 195, 240, 0.05)");
		}

		$('#p' + phase + i + 'id').text(id);
		$('#p' + phase + i + 'res').text(d.Top[i].Result.Win + ' / ' + d.Top[i].Result.Equal + ' / ' + d.Top[i].Result.Lose);
		$('#p' + phase + i + 'score').text(d.Top[i].Score.toFixed(2));
		if (d.Top[i].Passed) {
			$('#p' + phase + i + 'passed').html('<i class="fa fa-check"></i></td>');
		} else {
			$('#p' + phase + i + 'passed').html('<i class="fa fa-times"></i></td>');
		}
		$('#p' + phase + i + 'show a').attr('href', '/api/warrior/' + id);
	}
}

function toID(c) {
	var s = '';
	for (var i = 0; i < c.length - 1; i++) {
		s += c[i] + ',';
	}
	s += c[c.length - 1];
	return s;
}

function call(url) {
	var xmlhttp;

	if (window.XMLHttpRequest) {
		xmlhttp = new XMLHttpRequest();
	} else {
		xmlhttp = new ActiveXObject('Microsoft.XMLHTTP');
	}

	xmlhttp.open('GET', url, true);
	xmlhttp.send();
}
