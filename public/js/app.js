getJson('/api/metrics', function(d) {});

function exit() {
	getJson('/api/stop', function(d) { });
}

function reload() {
	addRounds();
	for (var i = 0; i < phases; i++) {
		reloadPhase(i);
	}
}

function addRounds() {
	getJson('/api/metrics', function(d) {
		$('#rps').text(d.Rounds);
		rchart.flow({
			columns: [
				['d', d.Rounds]
			],
			length: 1,
		});

		$('#fps').text(d.Fights);
		fchart.flow({
			columns: [
				['d', d.Fights]
			],
			length: 1,
		});
	});
}

function reloadCharts() {
	for (var i = 0; i < phases; i++) {
		reloadChart(i);
	}
}

var chartData = {};

function reloadChart(phase) {
	if (charts[phase].data()[0].values[0].value == chartData[phase].Passed && charts[phase].data()[1].values[0].value == chartData[phase].Failed) {
		return;
	}
	charts[phase].load({
		columns: [
			['Passed', chartData[phase].Passed],
			['Failed', chartData[phase].Failed],
		]
	});
}

function reloadPhase(phase) {
	getJson('/api/phase/' + phase, function(d) {
		chartData[phase] = {
			Passed: d.Passed,
			Failed: d.Failed
		};

		$('#p' + phase + 'total').text(d.Total);
		$('#p' + phase + 'failed').text(d.Failed);
		$('#p' + phase + 'passed').text(d.Passed);

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
			$('#p' + phase + i + 'score').text(d.Top[i].Score);
			if (d.Top[i].Passed) {
				$('#p' + phase + i + 'passed').html('<i class="fa fa-check"></i></td>');
			} else {
				$('#p' + phase + i + 'passed').html('<i class="fa fa-times"></i></td>');
			}
			$('#p' + phase + i + 'show').attr('href', '/api/warrior/' + id);
		}
	});
}

function toID(c) {
	var s = '';
	for (var i = 0; i < c.length - 1; i++) {
		s += c[i] + ',';
	}
	s += c[c.length - 1];
	return s;
}

function getJson(url, success) {
	var xmlhttp;

	if (window.XMLHttpRequest) {
		xmlhttp = new XMLHttpRequest();
	} else {
		xmlhttp = new ActiveXObject('Microsoft.XMLHTTP');
	}

	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == XMLHttpRequest.DONE) {
			if (xmlhttp.status == 200) {
				success(JSON.parse(xmlhttp.responseText));
			}
		}
	};

	xmlhttp.open('GET', url, true);
	xmlhttp.send();
}
