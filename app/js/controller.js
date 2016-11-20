function MainDashboard($scope, $http, $routeParams, $location) {
	'use strict';

	$http.get('/api/v1/metrics').success(function(data) {
		var metrics = [];
		if (data.length) {
			metrics = data;
		} else {
			$scope.error = true;
			return;
		}

		$http.get('/api/v1/clusters').success(function(data) {
			for (var i = 0, l = metrics.length; i < l; i++ ) {
				metrics[i].cluster = jlinq.from(data)
					.starts('name', metrics[i].cluster)
					.ends('name', metrics[i].cluster)
					.select()[0];
			}
		});

		$http.get('/api/v1/timelines').success(function(data) {
			$scope.metrics = [];
			var j = 0;
			for (var i = 0, l = metrics.length; i < l; i++ ) {
				var id = metrics[i].id;
				if (id in data) {
					$scope.metrics[j] = metrics[i];
					$scope.metrics[j].chartData = [{"key": id, "values": data[id]}];
					$scope.metrics[j].link = id.replace(".", "_");
					j++;
				}
			}
		});

		var format = d3.format(',');
		$scope.valueFormatFunction = function(){
			return function(d) {
				return format(d);
			};
		};

		$scope.$on('elementClick.directive', function(event, data) {
			var build = data.data[0],
				metric = event.targetScope.id,
				a = $("#run_"  + metric);
			a.attr("href", "/#/runs/" + metric + "/" + build);
			a[0].click();
		});
	});
}

function MenuRouter($scope, $routeParams, $location) {
	$scope.activeOS = $routeParams.os;
	$scope.activeComponent = $routeParams.component;
	$scope.activeCategory = $routeParams.category;

	$scope.setActiveOS = function (os) {
		$location.path("/timeline/" + os + "/" + $scope.activeComponent + "/" + $scope.activeCategory);
	};

	$scope.setActiveComponent = function (component) {
		$location.path("/timeline/" + $scope.activeOS + "/" + component + "/" + $scope.components[component].categories[0].id);
	};

	$scope.setActiveCategory = function (category) {
		$location.path("/timeline/" + $scope.activeOS + "/" + $scope.activeComponent + "/" + category);
	};

	DefineComponents($scope, $location);
	DefineCategories($scope, $location);
	DefineOS($scope, $location);
}

function DefineOS($scope, $location) {
	$scope.oses = ["Linux", "Windows"];

	$scope.byOS = function(entry) {
		var entryOS = entry.cluster.os;
		switch($scope.activeOS) {
			case "Linux":
				return entryOS.substring(0, 7) !== "Windows";
			case "Windows":
				return entryOS.substring(0, 7) === "Windows";
		}
	};
}

function DefineComponents($scope, $location) {
	$scope.components = {
		kv: {
			title: "KV",
			categories: [{
				id: "max_ops", title: "Max Throughput"
			}, {
				id: "latency", title: "Latency"
			}, {
				id: "_io_", title: "Storage"
			}, {
				id: "observe", title: "Observe"
			}, {
				id: "subdoc", title: "Sub Doc"
			}, {
				id: "warmup", title: "Warmup"
			}, {
				id: "fragmentation", title: "Fragmentation"
			}, {
				id: "compact", title: "Compaction"
			}, {
				id: "dcp", title: "DCP"
			}]
		},
		reb: {
			title: "Rebalance",
			categories: [{
				id: "empty", title: "Empty"
			}, {
				id: "kv", title: "KV"
			}, {
				id: "views", title: "Views"
			}, {
				id: "xdcr", title: "XDCR"
			}, {
				id: "failover", title: "Failover"
			}]
		},
		index: {
			title: "View Indexing",
			categories: [{
				id: "init", title: "Initial"
			}, {
				id: "incr", title: "Incremental"
			}]
		},
		query: {
			title: "View Query",
			categories: [{
				id: "lat", title: "Bulk Latency"
			}, {
				id: "by_type", title: "Latency by Query Type"
			}, {
				id: "throughput", title: "Throughput"
			}]
		},
		n1ql: {
			title: "N1QL",
			categories: [{
				id: "Q1_Q3_thr", title: "Q1-Q3 Throughput"
			}, {
				id: "Q1_Q3_lat", title: "Q1-Q3 Latency"
			}, {
				id: "Q5_Q7_thr", title: "Q5-Q7 Throughput"
			}, {
				id: "CI_thr", title: "Covering Indexes"
			}, {
				id: "array", title: "Array Indexing"
			}, {
				id: "join_unnest", title: "JOIN & UNNEST"
			}, {
				id: "dml", title: "DML"
			}]
		},
		secondary: {
			title: "2i",
			categories: [{
				id: "fdb_lat", title: "Latency FDB"
			}, {
				id: "fdb_thr", title: "Throughput FDB"
			}, {
				id: "fdb_init", title: "Initial FDB"
			}, {
				id: "fdb_incr", title: "Incremental FDB"
			}, {
				id: "fdb_standalone", title: "Standalone FDB"
			}, {
				id: "moi_lat", title: "Latency MOI"
			}, {
				id: "moi_thr", title: "Throughput MOI"
			}, {
				id: "moi_init", title: "Initial MOI"
			}, {
				id: "moi_incr", title: "Incremental MOI"
			}, {
				id: "fdb", title: "FDB"
			}, {
				id: "moi", title: "MOI"
			}]
		},
		xdcr: {
			title: "XDCR",
			categories: [{
				id: "init", title: "Initial"
			}, {
				id: "reb", title: "Initial+Rebalance"
			}, {
				id: "ongoing", title: "Ongoing"
			}, {
				id: "lww", title: "LWW"
			}]
		},
		fts: {
			title: "FTS",
			categories: [{
				id: "kvlatency", title: "KV Latency"
			}, {
				id: "kvthroughput", title: "KV Throughput"
			}, {
				id: "latency", title: "Latency"
			}, {
				id: "throughput", title: "Throughput"
			}, {
				id: "index", title: "Index"
			}, {
				id: "latency3", title: "3 Node Latency"
			}, {
				id: "throughput3", title: "3 Node Throughput "
			}, {
				id: "index3", title: "3 Node Index"
			}, {
				id: "elastic", title: "ElasticSearch"
			}]
		},
		ycsb: {
			title: "YCSB",
			categories: [{
				id: "workloada", title: "Workload A"
			}, {
				id: "workloadc", title: "Workload C"
			}, {
				id: "workloade", title: "Workload E"
			}]
		},
		tools: {
			title: "Tools",
			categories: [{
				id: "backup", title: "Backup"
			}, {
				id: "restore", title: "Restore"
			}, {
				id: "import", title: "Import"
			}, {
				id: "export", title: "Export"
			}]
		}
	};
}

function DefineCategories($scope, $location) {
	$scope.byComponentAndCategory = function(entry) {
		switch($scope.activeComponent) {
			case "kv":
				if (entry.component === $scope.activeComponent) {
					return entry.category === $scope.activeCategory;
				}
				break;
			case "reb":
				if (entry.component === $scope.activeComponent) {
					return entry.category === $scope.activeCategory;
				}
				break;
			case "index":
				if (entry.component === $scope.activeComponent) {
					return entry.id.indexOf($scope.activeCategory) !== -1;
				}
				break;
			case "query":
				if (entry.component === $scope.activeComponent) {
					return entry.category === $scope.activeCategory;
				}
				break;
			case "n1ql":
				if (entry.component === $scope.activeComponent) {
					return entry.category === $scope.activeCategory;
				}
				break;
			case "secondary":
				if (entry.component === $scope.activeComponent) {
					return bySecondaryCategory(entry);
				}
				break;
			case "xdcr":
				if (entry.component === $scope.activeComponent) {
					return entry.category === $scope.activeCategory;
				}
				break;
			case "fts":
				if (entry.component === $scope.activeComponent) {
					return byFTSCategory(entry);
				}
				break;
			case "ycsb":
				if (entry.component === $scope.activeComponent) {
					return entry.id.indexOf($scope.activeCategory) !== -1;
				}
				break;
			case "tools":
				if (entry.component === $scope.activeComponent) {
					return entry.id.indexOf($scope.activeCategory) !== -1;
				}
				break;
			default:
				return false;
		}
	};

	var bySecondaryCategory = function(entry) {
		switch($scope.activeCategory) {
			case "fdb_thr":
				if (entry.id.indexOf("_scanthr") !== -1 && entry.id.indexOf("_fdb_") !== -1) {
					return true;
				}
				break;
			case "fdb_lat":
				if (entry.id.indexOf("_scanlatency") !== -1 && entry.id.indexOf("_fdb_") !== -1) {
					return true;
				}
				break;
			case "fdb_init":
				if (entry.id.indexOf("_initial_") !== -1 && entry.id.indexOf("_fdb_") !== -1) {
					return true;
				}
				break;
			case "fdb_incr":
				if (entry.id.indexOf("_incremental_") !== -1 && entry.id.indexOf("_fdb_") !== -1) {
					return true;
				}
				break;
			case "fdb_standalone":
				if (entry.id.indexOf("_standalone_") !== -1 && entry.id.indexOf("_fdb_") !== -1) {
					return true;
				}
				break;
			case "moi_thr":
				if (entry.id.indexOf("_scanthr") !== -1 && entry.id.indexOf("_moi_") !== -1) {
					return true;
				}
				break;
			case "moi_lat":
				if (entry.id.indexOf("_scanlatency") !== -1 && entry.id.indexOf("_moi_") !== -1) {
					return true;
				}
				break;
			case "moi_init":
				if (entry.id.indexOf("_initial_") !== -1 && entry.id.indexOf("_moi_") !== -1) {
					return true;
				}
				break;
			case "moi_incr":
				if (entry.id.indexOf("_incremental_") !== -1 && entry.id.indexOf("_moi_") !== -1) {
					return true;
				}
				break;
			case "fdb":
				if (entry.id.indexOf("_fdb_") !== -1) {
					return true;
				}
				break;
			case "moi":
				if (entry.id.indexOf("_moi_") !== -1) {
					return true;
				}
				break;

			default:
				return false;
		}
	};

	var byFTSCategory = function(entry) {
		switch($scope.activeCategory) {
			case "elastic":
				if (entry.id.indexOf("elastic") !== -1) {
					return true;
				}
				break;

			case "latency":
				if (entry.id.indexOf("latency") !== -1 &&
					entry.id.indexOf("3nodes") === -1 &&
					entry.id.indexOf("kv") === -1 ){
					return true;
				}
				break;

			case "throughput":
				if (entry.id.indexOf("throughput") !== -1 &&
					entry.id.indexOf("3nodes") === -1 &&
					entry.id.indexOf("kv") === -1 ) {
					return true;
				}
				break;

			case "index":
				if (entry.id.indexOf("index") !== -1 &&
					entry.id.indexOf("3nodes") === -1) {
					return true;
				}
				break;

			case "latency3":
				if (entry.id.indexOf("latency") !== -1 &&
					entry.id.indexOf("3nodes") !== -1 &&
					entry.id.indexOf("kv") === -1 ){
					return true;
				}
				break;

			case "throughput3":
				if (entry.id.indexOf("throughput") !== -1 &&
					entry.id.indexOf("3nodes") !== -1 &&
					entry.id.indexOf("kv") === -1 ){
					return true;
				}
				break;

			case "index3":
				if (entry.id.indexOf("index") !== -1 &&
					entry.id.indexOf("3nodes") !== -1) {
					return true;
				}
				break;

			case "kvlatency":
				if (entry.id.indexOf("latency") !== -1 &&
					entry.id.indexOf("3nodes") === -1 &&
					entry.id.indexOf("kv") !== -1 ){
					return true;
				}
				break;

			case "kvthroughput":
				if (entry.id.indexOf("throughput") !== -1 &&
					entry.id.indexOf("3nodes") === -1 &&
					entry.id.indexOf("kv") !== -1 ){
					return true;
				}
				break;
			default:
				return false;
		}
	};
}

function RunList($scope, $routeParams, $http) {
	$http({method: 'GET', url: '/api/v1/runs/' + $routeParams.metric + '/' + $routeParams.build})
	.success(function(data) {
		$scope.runs = data;
	});
}

function AdminList($scope, $http) {
	$scope.deleteBenchmark = function(benchmark) {
		$http({method: 'DELETE', url: '/api/v1/benchmarks/' + benchmark.id})
		.success(function(data) {
			var index = $scope.benchmarks.indexOf(benchmark);
			$scope.benchmarks.splice(index, 1);
		});
	};

	$scope.reverseHidden = function(id) {
		$http({method: 'PATCH', url: '/api/v1/benchmarks/' + id});
	};

	$http.get('/api/v1/benchmarks').success(function(data) {
		$scope.benchmarks = data;

		$scope.components = {
			kv:        "KV",
			reb:       "Rebalance",
			index:     "View Indexing",
			query:     "View Query",
			n1ql:      "N1QL",
			secondary: "2i",
			xdcr:      "XDCR",
			fts:       "FTS",
			ycsb:      "YCSB",
			tools:     "Tools"
		};

		$scope.activeComponent = "kv";

		$scope.setActiveComponent = function (value) {
			$scope.activeComponent = value;
		};

		$scope.byComponent = function(entry) {
			if (entry.metric.indexOf($scope.activeComponent) === 0){
				return true;
			}
			return false;
		};
	});
}
