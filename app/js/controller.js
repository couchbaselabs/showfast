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
						return entryOS.substring(0, 7) !== "Windows";
				}
			};
		});

		$http.get('/all_timelines').success(function(data) {
			for (var i = 0, l = $scope.metrics.length; i < l; i++ ) {
				var id = $scope.metrics[i].id;
				$scope.metrics[i].chartData = [{"key": id, "values": data[id]}];
				$scope.metrics[i].link = id.replace(".", "_");
			}
		});

		$scope.categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "bandr", "title": "Backup"
		}, {
			"id": "beam", "title": "beam.smp"
		}, {
			"id": "kv", "title": "KV"
		}, {
			"id": "subdoc", "title": "Sub Doc"
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
		}

		];

		$scope.reb_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "empty", "title": "Empty"
		}, {
			"id": "kv", "title": "KV"
		}, {
			"id": "views", "title": "Views"
		}, {
			"id": "xdcr", "title": "XDCR"
		}, {
			"id": "failover", "title": "Failover"
		}, {
			"id": "failover_views", "title": "Failover+Views"
		}];

		$scope.idx_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "init", "title": "Initial"
		}, {
			"id": "incr", "title": "Incremental"
		}];

		$scope.secondary_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "lat", "title": "Latency"
		}, {
			"id": "thr", "title": "Throughput"
		}, {
			"id": "init", "title": "Initial"
		}, {
			"id": "incr", "title": "Incremental"
		}, {
			"id": "fdb", "title": "fdb"
		}, {
			"id": "memdb", "title": "MOI"
		}];
		
        $scope.n1ql_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "lat", "title": "Latency"
		}, {
            "id": "thr", "title": "Throughput"
		}, {
            "id": "vs", "title": "CE vs. EE"
		}, {
            "id": "part", "title": "Partitioned"
		}, {
			"id": "moi", "title": "MOI"
		}, {
			"id": "array", "title": "Array"
		}, {
			"id": "CI", "title": "Covering"
		}, {
			"id": "_Q", "title": "Query"
		},  {
            "id": "wl", "title": "Archive"
		}];
		
		$scope.xdcr_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "ongoing", "title": "Ongoing"
		}, {
			"id": "init", "title": "Initial"
		}];

		$scope.backup_restore_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "backup", "title": "Backup"
		}, {
			"id": "restore", "title": "Restore"
		}];

		$scope.beam_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "kv", "title": "KV"
		}, {
			"id": "query", "title": "Views"
		}, {
			"id": "xdcr", "title": "XDCR"
		}];

		$scope.kv_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "latency", "title": "Latency"
		}, {
			"id": "observe", "title": "Observe"
		}, {
			"id": "warmup", "title": "Warmup"
		}, {
			"id": "fragmentation", "title": "Memory"
		}, {
			"id": "bgfetcher", "title": "BgFetcher"
		}, {
			"id": "drain", "title": "Flusher"
		}, {
			"id": "max_ops", "title": "Max Throughput"
		}];

		$scope.query_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "lat", "title": "Bulk Latency"
		}, {
			"id": "dev", "title": "Latency by Query Type"
		}, {
			"id": "thr", "title": "Throughput"
		}];

		$scope.spatial_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "lat", "title": "Latency"
		}, {
			"id": "index", "title": "Index Creation"
		}, {
			"id": "thr", "title": "Throughput"
		}, {
			"id": "reb", "title": "Rebalance"
		}];

		$scope.spatial_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "lat", "title": "Latency"
		}, {
			"id": "index", "title": "Index Creation"
		}, {
			"id": "thr", "title": "Throughput"
		}, {
			"id": "reb", "title": "Rebalance"
		}];

		$scope.subdoc_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "latency", "title": "Latency"
		}, {
			"id": "bandwidth", "title": "Bandwidth"
		}];

		$scope.fts_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "latency", "title": "Latency"
		}, {
			"id": "throughput", "title": "Throughput"
		},{
			"id": "index", "title": "Index"
		},{
            "id": "elastic", "title": "ElasticSearch"
        }];

        $scope.ycsb_categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "workloada", "title": "Workload A"
		}, {
			"id": "workloade", "title": "Workload E"
		}];

		$scope.selectedCategory = $.cookie("selectedCategory") || "all";
		$scope.selectedRebCategory = $.cookie("selectedRebCategory") || "all";
		$scope.selectedIdxCategory = $.cookie("selectedIdxCategory") || "all";
		$scope.selectedXdcrCategory = $.cookie("selectedXdcrCategory") || "all";
		$scope.selectedBeamCategory = $.cookie("selectedBeamCategory") || "all";
		$scope.selectedBackupRestoreCategory = $.cookie("selectedBackupRestoreCategory") || "all";
		$scope.selectedKVCategory = $.cookie("selectedKVCategory") || "all";
		$scope.selectedQueryCategory = $.cookie("selectedQueryCategory") || "all";
		$scope.selectedQueryCategory = $.cookie("selectedSpatialCategory") || "all";
		$scope.selectedN1QLCategory = $.cookie("selectedN1QLCategory") || "all";
		$scope.selectedSecondaryCategory = $.cookie("selectedSecondaryCategory") || "all";
		$scope.selectedSubdocCategory = $.cookie("selectedSubdocCategory") || "all";
		$scope.selectedFtsCategory = $.cookie("selectedFtsCategory") || "all";
		$scope.selectedYcsbCategory = $.cookie("selectedYcsbCategory") || "all";

		$scope.setSelectedCategory = function (value) {
			$scope.selectedCategory = value;
			$.cookie("selectedCategory", value);
		};

		$scope.setSelectedRebCategory = function (value) {
			$scope.selectedRebCategory = value;
			$.cookie("selectedRebCategory", value);
		};

		$scope.setSelectedIdxCategory = function (value) {
			$scope.selectedIdxCategory = value;
			$.cookie("selectedIdxCategory", value);
		};

		$scope.setSelectedXdcrCategory = function (value) {
			$scope.selectedXdcrCategory = value;
			$.cookie("selectedXdcrCategory", value);
		};

		$scope.setSelectedBackupRestoreCategory = function (value) {
			$scope.selectedBackupRestoreCategory = value;
			$.cookie("selectedBackupRestoreCategory", value);
		};

		$scope.setSelectedBeamCategory = function (value) {
			$scope.selectedBeamCategory = value;
			$.cookie("selectedBeamCategory", value);
		};

		$scope.setSelectedKVCategory = function (value) {
			$scope.selectedKVCategory = value;
			$.cookie("selectedKVCategory", value);
		};

		$scope.setSelectedQueryCategory = function (value) {
			$scope.selectedQueryCategory = value;
			$.cookie("selectedQueryCategory", value);
		};
		
		$scope.setSelectedSpatialCategory = function (value) {
			$scope.selectedSpatialCategory = value;
			$.cookie("selectedSpatialCategory", value);
		};
		
		$scope.setSelectedSecondaryCategory = function (value) {
			$scope.selectedSecondaryCategory = value;
			$.cookie("selectedSecondaryCategory", value);
		};
		
		$scope.setSelectedN1QLCategory = function (value) {
			$scope.selectedN1QLCategory = value;
			$.cookie("selectedN1QLCategory", value);
		};

		$scope.setSelectedSubdocCategory = function (value) {
			$scope.selectedSubdocCategory = value;
			$.cookie("selectedSubdocCategory", value);
		};

		$scope.setSelectedFtsCategory = function (value) {
			$scope.selectedFtsCategory = value;
			$.cookie("selectedFtsCategory", value);
		};

		$scope.setSelectedYcsbCategory = function (value) {
			$scope.selectedYcsbCategory = value;
			$.cookie("selectedYcsbCategory", value);
		};

		$scope.byCategory = function(entry) {
			var selectedCategory = $scope.selectedCategory,
				entryCategory = entry.id.substring(0, selectedCategory.length);

			switch(selectedCategory) {
				case "all":
					return true;
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
				case "bandr":
					if (entryCategory === selectedCategory) {
						return byBackupRestoreCategory(entry);
					}
					break;
				case "beam":
					if (entryCategory === selectedCategory) {
						return byBeamCategory(entry);
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
				case "spatial":
					if (entryCategory === selectedCategory) {
						return bySpatialCategory(entry);
					}
					break;
				case "n1ql":
					if (entryCategory === selectedCategory) {
						return byN1QLCategory(entry);
					}
					break;
				case "xdcr":
					if (entryCategory === selectedCategory) {
						return byXdcrCategory(entry);
					}
					break;
				case "subdoc":
					if (entryCategory === selectedCategory) {
						return bySubdocCategory(entry);
					}
					break;
				case "fts":
					if (entryCategory === selectedCategory) {
						return byFtsCategory(entry);
					}
					break;
				case "ycsb":
					if (entryCategory === selectedCategory) {
						return byYcsbCategory(entry);
					}
					break;	
				case entryCategory:
					return true;
				default:
					return false;
			}
		};

		var byIdxCategory = function(entry) {
			var selectedIdxCategory = $scope.selectedIdxCategory;

			if (selectedIdxCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedIdxCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};
		
		var bySecondaryCategory = function(entry) {
			var selectedSecondaryCategory = $scope.selectedSecondaryCategory;

			if (selectedSecondaryCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedSecondaryCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};

		var byRebCategory = function(entry) {
			var selectedRebCategory = $scope.selectedRebCategory;

			if (selectedRebCategory === "all") {
				return true;
			} else {
				switch(selectedRebCategory) {
					case "failover_views":
						if (entry.id.indexOf("failover") !== -1 &&
								entry.id.indexOf("views") !== -1) {
							return true;
						}
						break;
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
			}
		};

		var byXdcrCategory = function(entry) {
			var selectedXdcrCategory = $scope.selectedXdcrCategory;

			switch(selectedXdcrCategory) {
					case "all":
						return true;
					case "init":
						if (entry.id.indexOf("init") !== -1) {
							return true;
						}
						break;
					case "ongoing":
						if (entry.id.indexOf("init") === -1) {
							return true;
						}
						break;
					default:
						return false;
			}
		};

		var byBackupRestoreCategory = function(entry) {
			var selectedBackupRestoreCategory = $scope.selectedBackupRestoreCategory;

			if (selectedBackupRestoreCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedBackupRestoreCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};

		var byBeamCategory = function(entry) {
			var selectedBeamCategory = $scope.selectedBeamCategory;

			if (selectedBeamCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedBeamCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};

		var byKVCategory = function(entry) {
			var selectedKVCategory = $scope.selectedKVCategory;

			if (selectedKVCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedKVCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};

		var byQueryCategory = function(entry) {
			var selectedQueryCategory = $scope.selectedQueryCategory;

			if (selectedQueryCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedQueryCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};
		
		var bySpatialCategory = function(entry) {
			var selectedSpatialCategory = $scope.selectedSpatialCategory;

			if (selectedSpatialCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedSpatialCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};
		
        var byN1QLCategory = function(entry) {
			var selectedN1QLCategory = $scope.selectedN1QLCategory;

            if (selectedN1QLCategory === "all") {
                if (entry.id.indexOf("wl") == -1) {
                    return true;
                } else {
                    return false;
                }
            } else if (selectedN1QLCategory === "part") {
                if (entry.id.indexOf(selectedN1QLCategory) !== -1 ||
                    entry.id.indexOf("thr_Q2") !== -1 ||
                    entry.id.indexOf("thr_Q3") !== -1) {
                    return true;
                } else {
                    return false;
                }

            } else {
                if (entry.id.indexOf(selectedN1QLCategory) !== -1 &&
                    entry.id.indexOf("part") == -1) {
                    return true;
                } else {
                    return false;
                }
            }
		};

		var bySubdocCategory = function(entry) {
			var selectedSubdocCategory = $scope.selectedSubdocCategory;

			if (selectedSubdocCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedSubdocCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};

		var byFtsCategory = function(entry) {
			var selectedFtsCategory = $scope.selectedFtsCategory;

			if (selectedFtsCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedFtsCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};

		var byYcsbCategory = function(entry) {
			var selectedYcsbCategory = $scope.selectedYcsbCategory;

			if (selectedYcsbCategory === "all") {
				return true;
			} else {
				if (entry.id.indexOf(selectedYcsbCategory) !== -1) {
					return true;
				} else {
					return false;
				}
			}
		};

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

function Feed($scope, $http) {
	$http.get('/all_feed_records').success(function(data) {
		$scope.records = data;
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
