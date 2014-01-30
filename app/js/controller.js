/*jshint jquery: true, browser: true*/
/*global jlinq, d3: true*/


function MetricList($scope, $http) {
	'use strict';

	$http.get('/all_metrics').success(function(data) {
		if (data.length) {
			$scope.metrics = data;
		} else {
			$scope.error = true;
			return;
		}

		$http.get('/all_clusters').success(function(data) {
			for (var i = 0, l = $scope.metrics.length; i < l; i++ ) {
				$scope.metrics[i].cluster = jlinq.from(data)
					.starts('Name', $scope.metrics[i].cluster)
					.ends('Name', $scope.metrics[i].cluster)
					.select()[0];
			}
		});

		$http.get('/all_timelines').success(function(data) {
			for (var i = 0, l = $scope.metrics.length; i < l; i++ ) {
				var id = $scope.metrics[i].id;
				$scope.metrics[i].chartData = [{"values": data[id]}];
			}
		});

		$scope.categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "comp", "title": "Compaction"
		}, {
			"id": "beam", "title": "beam.smp"
		}, {
			"id": "index", "title": "Indexing"
		}, {
			"id": "reb", "title": "Rebalance"
		}, {
			"id": "kv", "title": "KV"
		}, {
			"id": "query", "title": "Query"
		}, {
			"id": "xdcr", "title": "XDCR"
		}];

		$scope.selectedCategory = $.cookie("selectedCategory") || "all";

		$scope.setSelectedCategory = function (value) {
			$scope.selectedCategory = value;
			$.cookie("selectedCategory", value);
		};

		$scope.byCategory = function(entry) {
			var selectedCategory = $scope.selectedCategory,
				entryCategory = entry.id.substring(0, selectedCategory.length);

			switch(selectedCategory) {
				case "all":
					return true;
				case entryCategory:
					return true;
				default:
					return false;
			}
		};

		$scope.oses = ["All", "Windows", "Linux"];

		$scope.selectedOS = $.cookie("selectedOS") || "All";

		$scope.setSelectedOS = function (value) {
			$scope.selectedOS = value;
			$.cookie("selectedOS", value);
		};

		$scope.byOS = function(entry) {
			var entryOS = entry.cluster.OS,
				selectedOS = $scope.selectedOS;

			switch(selectedOS) {
				case "All":
					return true;
				case "Windows":
					return entryOS.substring(0, 7) === "Windows";
				case "Linux":
					return !(entryOS.substring(0, 7) === "Windows");
			}
		};

		var format = d3.format(',');
		$scope.valueFormatFunction = function(){
			return function(d) {
				return format(d);
			};
		};

		$scope.$on('barClick', function(event, data) {
			var build = data.point[0],
				metric = event.targetScope.id.substring(6),
				a = $("#run_"  + metric);
			a.attr("href", "/#/runs/" + metric + "/" + build);
			a[0].click();
		});
	});
}

function RunList($scope, $routeParams, $http) {
	$http({method: 'GET', url: '/all_runs', params: $routeParams})
	.success(function(data) {
		$scope.runs = data;
	});
}

function AdminList($scope, $http) {
	$scope.deleteBenchmark = function(benchmark) {
		$http({method: 'POST', url: '/delete', data: {id: benchmark.id}})
		.success(function(data) {
			var index = $scope.benchmarks.indexOf(benchmark);
			$scope.benchmarks.splice(index, 1);
		});
	};

	$scope.reverseObsolete = function(id) {
		$http({method: 'POST', url: '/reverse_obsolete', data: {id: id}});
	};

	$http.get('/all_benchmarks').success(function(data) {
		$scope.benchmarks = data;
	});
}

function GetComparison($scope, $http) {
	var params = {
		'baseline': $scope.selectedBaseline,
		'target': $scope.selectedTarget
	};
	$http({method: 'GET', url: '/get_comparison', params: params})
	.success(function(data) {
		$scope.metrics = data;
	});
}

function ReleaseList($scope, $http) {
	$http.get('/all_releases').success(function(data) {
		$scope.baselines = data;
		$scope.targets = data;

		$scope.selectedBaseline = $.cookie('selectedBaseline') || data[0];
		$scope.selectedTarget = $.cookie('selectedTarget') || data[0];

		$scope.setSelectedBaseline = function (value) {
			$scope.selectedBaseline = value;
			$.cookie('selectedBaseline', value, {expires: 60});
			GetComparison($scope, $http);
		};

		$scope.setSelectedTarget = function (value) {
			$scope.selectedTarget = value;
			$.cookie('selectedTarget', value, {expires: 60});
			GetComparison($scope, $http);
		};

		GetComparison($scope, $http);
	});
}
