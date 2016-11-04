/*jshint jquery: true, browser: true*/
/*global jlinq, d3: true*/


function Timeline($scope, $http) {
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

			$scope.oses = ["Linux", "Windows"];
			$scope.selectedOS = $.cookie("selectedOSV1") || "Linux";
			$scope.setSelectedOS = function (value) {
				$scope.selectedOS = value;
				$.cookie("selectedOSV1", value);
			};
			$scope.byOS = function(entry) {
				var entryOS = entry.cluster.os,
					selectedOS = $scope.selectedOS;
				switch(selectedOS) {
					case "Linux":
						return entryOS.substring(0, 7) !== "Windows";
					case "Windows":
						return entryOS.substring(0, 7) === "Windows";
				}
			};
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

		DefineCategories($scope);
		DefineSubCategories($scope);

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

function DefineCategories($scope) {
	$scope.categories = [{
		"id": "kv", "title": "KV"
	}, {
		"id": "reb", "title": "Rebalance"
	}, {
		"id": "index", "title": "View Indexing"
	}, {
		"id": "query", "title": "View Query"
	}, {
		"id": "n1ql", "title": "N1QL"
	}, {
		"id": "secondary", "title": "2i"
	}, {
		"id": "xdcr", "title": "XDCR"
	}, {
		"id": "fts", "title": "FTS"
	}, {
		"id": "ycsb", "title": "YCSB"
	}, {
		"id": "tools", "title": "Tools"
	}];

	$scope.selectedCategory = $.cookie("selectedCategoryV1") || $scope.categories[0].id;

	$scope.setSelectedCategory = function (value) {
		$scope.selectedCategory = value;
		$.cookie("selectedCategoryV1", value);
	};
}

function DefineSubCategories($scope) {
	$scope.rebCategories = [{
		"id": "empty", "title": "Empty"
	}, {
		"id": "kv", "title": "KV"
	}, {
		"id": "views", "title": "Views"
	}, {
		"id": "xdcr", "title": "XDCR"
	}, {
		"id": "failover", "title": "Failover"
	}];

	$scope.idxCategories = [{
		"id": "init", "title": "Initial"
	}, {
		"id": "incr", "title": "Incremental"
	}];

	$scope.secondaryCategories = [{
		"id": "fdb_lat", "title": "Latency FDB"
	}, {
		"id": "fdb_thr", "title": "Throughput FDB"
	}, {
		"id": "fdb_init", "title": "Initial FDB"
	}, {
		"id": "fdb_incr", "title": "Incremental FDB"
	}, {
		"id": "fdb_standalone", "title": "Standalone FDB"
	}, {
		"id": "moi_lat", "title": "Latency MOI"
	}, {
		"id": "moi_thr", "title": "Throughput MOI"
	}, {
		"id": "moi_init", "title": "Initial MOI"
	}, {
		"id": "moi_incr", "title": "Incremental MOI"
	}, {
		"id": "fdb", "title": "FDB"
	}, {
		"id": "moi", "title": "MOI"
	}];

	$scope.n1qlCategories = [{
		"id": "Q1_Q3_thr", "title": "Q1-Q3 Throughput"
	}, {
		"id": "Q1_Q3_lat", "title": "Q1-Q3 Latency"
	}, {
		"id": "Q5_Q7_thr", "title": "Q5-Q7 Throughput"
	}, {
		"id": "CI_thr", "title": "Covering Indexes"
	}, {
		"id": "array", "title": "Array Indexing"
	}, {
		"id": "join_unnest", "title": "JOIN & UNNEST"
	}, {
		"id": "dml", "title": "DML"
	}];

	$scope.xdcrCategories = [{
		"id": "init", "title": "Initial"
	}, {
		"id": "reb", "title": "Initial+Rebalance"
	}, {
		"id": "ongoing", "title": "Ongoing"
	}, {
		"id": "lww", "title": "LWW"
	}];

	$scope.toolsCategories = [{
		"id": "backup", "title": "Backup"
	}, {
		"id": "restore", "title": "Restore"
	}, {
		"id": "import", "title": "Import"
	}, {
		"id": "export", "title": "Export"
	}];

	$scope.kvCategories = [{
		"id": "max_ops", "title": "Max Throughput"
	}, {
		"id": "latency", "title": "Latency"
	}, {
		"id": "_io_", "title": "Storage"
	}, {
		"id": "observe", "title": "Observe"
	}, {
		"id": "subdoc", "title": "Sub Doc"
	}, {
		"id": "warmup", "title": "Warmup"
	}, {
		"id": "fragmentation", "title": "Fragmentation"
	}, {
		"id": "compact", "title": "Compaction"
	}];

	$scope.queryCategories = [{
		"id": "lat", "title": "Bulk Latency"
	}, {
		"id": "dev", "title": "Latency by Query Type"
	}, {
		"id": "thr", "title": "Throughput"
	}];

	$scope.ftsCategories = [{
		"id": "latency", "title": "Latency"
	}, {
		"id": "throughput", "title": "Throughput"
	}, {
		"id": "index", "title": "Index"
	}, {
		"id": "latency3", "title": "3 node Latency"
	}, {
		"id": "throughput3", "title": "3 node Throughput "
	}, {
		"id": "index3", "title": "3 node Index"
	}, {
		"id": "elastic", "title": "ElasticSearch"
	}];

	$scope.ycsbCategories = [{
		"id": "workloada", "title": "Workload A"
	}, {
		"id": "workloadc", "title": "Workload C"
	}, {
		"id": "workloade", "title": "Workload E"
	}];

	$scope.selectedRebCategory = $.cookie("selectedRebCategoryV1") || $scope.rebCategories[0].id;
	$scope.selectedIdxCategory = $.cookie("selectedIdxCategoryV1") || $scope.idxCategories[0].id;
	$scope.selectedXDCRCategory = $.cookie("selectedXDCRCategoryV1") || $scope.xdcrCategories[0].id;
	$scope.selectedToolsCategory = $.cookie("selectedToolsCategoryV1") || $scope.toolsCategories[0].id;
	$scope.selectedKVCategory = $.cookie("selectedKVCategoryV1") || $scope.kvCategories[0].id;
	$scope.selectedQueryCategory = $.cookie("selectedQueryCategoryV1") || $scope.queryCategories[0].id;
	$scope.selectedN1QLCategory = $.cookie("selectedN1QLCategoryV1") || $scope.n1qlCategories[0].id;
	$scope.selectedSecondaryCategory = $.cookie("selectedSecondaryCategoryV1") || $scope.secondaryCategories[0].id;
	$scope.selectedFTSCategory = $.cookie("selectedFTSCategoryV1") || $scope.ftsCategories[0].id;
	$scope.selectedYCSBCategory = $.cookie("selectedYCSBCategoryV1") || $scope.ycsbCategories[0].id;

	$scope.setSelectedRebCategory = function (value) {
		$scope.selectedRebCategory = value;
		$.cookie("selectedRebCategoryV1", value);
	};

	$scope.setSelectedIdxCategory = function (value) {
		$scope.selectedIdxCategory = value;
		$.cookie("selectedIdxCategoryV1", value);
	};

	$scope.setSelectedXDCRCategory = function (value) {
		$scope.selectedXDCRCategory = value;
		$.cookie("selectedXDCRCategoryV1", value);
	};

	$scope.setSelectedToolsCategory = function (value) {
		$scope.selectedToolsCategory = value;
		$.cookie("selectedToolsCategoryV1", value);
	};

	$scope.setSelectedKVCategory = function (value) {
		$scope.selectedKVCategory = value;
		$.cookie("selectedKVCategoryV1", value);
	};

	$scope.setSelectedQueryCategory = function (value) {
		$scope.selectedQueryCategory = value;
		$.cookie("selectedQueryCategoryV1", value);
	};

	$scope.setSelectedSecondaryCategory = function (value) {
		$scope.selectedSecondaryCategory = value;
		$.cookie("selectedSecondaryCategoryV1", value);
	};

	$scope.setSelectedN1QLCategory = function (value) {
		$scope.selectedN1QLCategory = value;
		$.cookie("selectedN1QLCategoryV1", value);
	};

	$scope.setSelectedFTSCategory = function (value) {
		$scope.selectedFTSCategory = value;
		$.cookie("selectedFTSCategoryV1", value);
	};

	$scope.setSelectedYCSBCategory = function (value) {
		$scope.selectedYCSBCategory = value;
		$.cookie("selectedYCSBCategory", value);
	};

	var byIdxCategory = function(entry) {
		var selectedIdxCategory = $scope.selectedIdxCategory;

		return entry.id.indexOf(selectedIdxCategory) !== -1;
	};

	$scope.byCategory = function(entry) {
		var selectedCategory = $scope.selectedCategory,
			entryCategory = entry.id.substring(0, selectedCategory.length);

		switch(selectedCategory) {
			case "index":
				if (entryCategory === selectedCategory) {
					return byIdxCategory(entry);
				}
				break;
			case "secondary":
				if (entryCategory === selectedCategory) {
					return bySecondaryCategory(entry);
				}
				break;
			case "tools":
				if (entryCategory === selectedCategory) {
					return byToolsCategory(entry);
				}
				break;
			case "reb":
				if (entryCategory === selectedCategory) {
					return byRebCategory(entry);
				}
				break;
			case "kv":
				if (entryCategory === selectedCategory) {
					return byKVCategory(entry);
				}
				break;
			case "query":
				if (entryCategory === selectedCategory) {
					return byQueryCategory(entry);
				}
				break;
			case "n1ql":
				if (entryCategory === selectedCategory) {
					return byN1QLCategory(entry);
				}
				break;
			case "xdcr":
				if (entryCategory === selectedCategory) {
					return byXDCRCategory(entry);
				}
				break;
			case "fts":
				if (entryCategory === selectedCategory) {
					return byFTSCategory(entry);
				}
				break;
			case "ycsb":
				if (entryCategory === selectedCategory) {
					return byYCSBCategory(entry);
				}
				break;
			case entryCategory:
				return true;
			default:
				return false;
		}
	};

	var bySecondaryCategory = function(entry) {
		var selectedSecondaryCategory = $scope.selectedSecondaryCategory;

		switch(selectedSecondaryCategory) {
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
				return entry.id.indexOf(selectedSecondaryCategory) !== -1;
		}
	};

	var byRebCategory = function(entry) {
		var selectedRebCategory = $scope.selectedRebCategory;

		switch(selectedRebCategory) {
			case "failover":
				if (entry.id.indexOf("failover") !== -1 &&
						entry.id.indexOf("views") === -1) {
					return true;
				}
				break;
			case "empty":
				if (entry.id.indexOf("0_kv_") !== -1) {
					return true;
				}
				break;
			case "views":
				if (entry.id.indexOf("views") !== -1 &&
						entry.id.indexOf("failover") === -1) {
					return true;
				}
				break;
			case "xdcr":
				if (entry.id.indexOf("xdcr") !== -1) {
					return true;
				}
				break;
			case "kv":
				if (entry.id.indexOf("0_kv_") === -1 &&
						entry.id.indexOf("views") === -1 &&
						entry.id.indexOf("failover") === -1 &&
						entry.id.indexOf("xdcr") === -1) {
					return true;
				}
				break;
			default:
				return false;
		}
	};

	var byXDCRCategory = function(entry) {
		var selectedXDCRCategory = $scope.selectedXDCRCategory;

		switch(selectedXDCRCategory) {
			case "lww":
				if (entry.id.indexOf("lww") !== -1) {
					return true;
				}
				break;
			case "init":
				if (entry.id.indexOf("init") !== -1 &&
						entry.id.indexOf("lww_init") === -1) {
					return true;
				}
				break;
			case "reb":
				if (entry.id.indexOf("reb") !== -1) {
					return true;
				}
				break;
			case "ongoing":
				if (entry.id.indexOf("init") === -1 &&
						entry.id.indexOf("reb") === -1 &&
							entry.id.indexOf("lww") === -1) {
					return true;
				}
				break;
			default:
				return false;
		}
	};

	var byToolsCategory = function(entry) {
		var selectedToolsCategory = $scope.selectedToolsCategory;

		return entry.id.indexOf(selectedToolsCategory) !== -1;
	};

	var byKVCategory = function(entry) {
		var selectedKVCategory = $scope.selectedKVCategory;

		return entry.id.indexOf(selectedKVCategory) !== -1;
	};

	var byQueryCategory = function(entry) {
		var selectedQueryCategory = $scope.selectedQueryCategory;

		return entry.id.indexOf(selectedQueryCategory) !== -1;
	};

	var byN1QLCategory = function(entry) {
		var selectedN1QLCategory = $scope.selectedN1QLCategory;

		switch(selectedN1QLCategory) {
			case "Q1_Q3_thr":
				if (entry.id.indexOf("_thr_") !== -1 &&
					entry.id.indexOf("_array_") === -1 && (
					entry.id.indexOf("_Q1") !== -1 ||
					entry.id.indexOf("_Q2_") !== -1 ||
					entry.id.indexOf("_Q3_") !== -1)) {
					return true;
				}
				break;
			case "Q1_Q3_lat":
				if (entry.id.indexOf("_lat_") !== -1 && (
					entry.id.indexOf("_Q1") !== -1 ||
					entry.id.indexOf("_Q2_") !== -1 ||
					entry.id.indexOf("_Q3_") !== -1)) {
					return true;
				}
				break;
			case "Q5_Q7_thr":
				if (entry.id.indexOf("_thr_") !== -1 &&
					entry.id.indexOf("_array_") === -1 && (
					entry.id.indexOf("_Q5") !== -1 ||
					entry.id.indexOf("_Q6_") !== -1 ||
					entry.id.indexOf("_Q7_") !== -1)) {
					return true;
				}
				break;
			case "join_unnest":
				if (entry.id.indexOf("_JOIN_") !== -1 ||
					entry.id.indexOf("_UNNEST_") !== -1) {
					return true;
				}
				break;
			case "CI_thr":
				if (entry.id.indexOf("_thr_") !== -1 &&
					entry.id.indexOf("_CI") !== -1) {
					return true;
				}
				break;
			case "CI_lat":
				if (entry.id.indexOf("_lat_") !== -1 &&
					entry.id.indexOf("_CI") !== -1) {
					return true;
				}
				break;
			case "array":
				if (entry.id.indexOf("_array_") !== -1) {
					return true;
				}
				break;
			case "dml":
				if (entry.id.indexOf("_IN") !== -1 ||
					entry.id.indexOf("_UP") !== -1) {
					return true;
				}
				break;
			default:
				return entry.id.indexOf(selectedN1QLCategory) !== -1;
		}
	};

	var byFTSCategory = function(entry) {
		var selectedFTSCategory = $scope.selectedFTSCategory;
		switch(selectedFTSCategory) {
			case "elastic":
				if (entry.id.indexOf("elastic") !== -1) {
					return true;
				}
				break;

			case "latency":
				if (entry.id.indexOf("latency") !== -1 &&
					entry.id.indexOf("3nodes") === -1) {
					return true;
				}
				break;

			case "throughput":
				if (entry.id.indexOf("throughput") !== -1 &&
					entry.id.indexOf("3nodes") === -1) {
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
					entry.id.indexOf("3nodes") !== -1) {
					return true;
				}
				break;

			case "throughput3":
				if (entry.id.indexOf("throughput") !== -1 &&
					entry.id.indexOf("3nodes") !== -1) {
					return true;
				}
				break;

			case "index3":
				if (entry.id.indexOf("index") !== -1 &&
					entry.id.indexOf("3nodes") !== -1) {
					return true;
				}
				 break;
		}
		return false;
	};

	var byYCSBCategory = function(entry) {
		var selectedYCSBCategory = $scope.selectedYCSBCategory;

		return entry.id.indexOf(selectedYCSBCategory) !== -1;
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

	$scope.reverseObsolete = function(id) {
		$http({method: 'PATCH', url: '/api/v1/benchmarks/' + id});
	};

	$http.get('/api/v1/benchmarks').success(function(data) {
		$scope.benchmarks = data;

		$scope.categories = [{
			"id": "kv", "title": "KV"
		}, {
			"id": "reb", "title": "Rebalance"
		}, {
			"id": "index", "title": "View Indexing"
		}, {
			"id": "query", "title": "View Query"
		}, {
			"id": "n1ql", "title": "N1QL"
		}, {
			"id": "secondary", "title": "2i"
		}, {
			"id": "xdcr", "title": "XDCR"
		}, {
			"id": "fts", "title": "FTS"
		}, {
			"id": "ycsb", "title": "YCSB"
		}, {
			"id": "tools", "title": "Tools"
		}];

		$scope.selectedCategory = $scope.categories[0].id;

		$scope.setSelectedCategory = function (value) {
			$scope.selectedCategory = value;
		};

		$scope.byCategory = function(entry) {
			if (entry.metric.indexOf($scope.selectedCategory) === 0){
				return true;
			}
			return false;
		};
	});
}
