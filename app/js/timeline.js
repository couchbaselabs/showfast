angular
    .module('showfast', ['ngRoute', 'nvd3ChartDirectives'])
    .config([
        '$routeProvider', function($routeProvider) {
            $routeProvider.
                when('/timeline/:os/:component/:category/:subCategory', {templateUrl: '/static/timeline.html', controller: MenuRouter, reloadOnSearch: false}).
                when('/runs/:metric/:build', {templateUrl: '/static/runs.html', controller: RunList}).
                when('/cloudTimeline/:os/:component/:category/:subCategory', {templateUrl: '/static/cloud_timeline.html', controller: CloudMenuRouter, reloadOnSearch: false}).
                otherwise({redirectTo: 'timeline/Linux/kv/max_ops/all'});
        }
    ])
    .controller('MainDashboard', MainDashboard)
    .directive('ngAnchor', function anchorDirective($location, $anchorScroll) {
        // Based on Ben Nadel's blog post:
        // https://www.bennadel.com/blog/2869-using-anchor-tags-and-url-fragment-links-in-angularjs.htm

        return ({
            link: link,
            restrict: "A"
        });

        function link(scope, element, attributes) {
            attributes.$observe("ngAnchor", configureHref);
            scope.$on("$locationChangeSuccess", configureHref);

            function configureHref() {
                var fragment = (attributes.ngAnchor || "");

                if (fragment.charAt(0) === "#") {
                    fragment = fragment.slice(1);
                }

                var routeValue = ("#" + $location.url().split("#").shift());
                var fragmentValue = ("#" + fragment);

                attributes.$set("href", (routeValue + fragmentValue));

                if ($location.hash() === fragment) {
                    element.addClass("active-anchor");
                    $anchorScroll();
                } else {
                    element.removeClass("active-anchor");
                }
            }
        }
    });

function MainDashboard($scope, $http, $routeParams) {
    var format = d3.format(',');
    $scope.valueFormatFunction = function() {
        return function(d) {
            return format(d);
        };
    };

    $scope.$on('elementClick.directive', function(event, data) {
        var build = data.data[0],
            metric = event.targetScope.id,
            a = $("#run_" + metric);
        a.attr("href", "/#/runs/" + metric + "/" + build);
        a[0].click();
    });
}

function MenuRouter($scope, $http, $routeParams, $location) {
    $scope.activeOS = $routeParams.os;
    $scope.activeComponent = $routeParams.component;
    $scope.activeCategory = $routeParams.category;
    $scope.activeSubCategory = $routeParams.subCategory;
    $scope.testType = "All";
    DefineMenu($scope, $http);

    var url = '/api/v1/metrics/' + $scope.activeComponent + "/" + $scope.activeCategory + "/" + $scope.activeSubCategory;
    $http.get(url).success(function(metrics) {
        $scope.metrics = metrics;

        $.each(metrics, function(i, metric) {
            $http.get('/api/v1/timeline/' + metric.id).success(function(data) {
                $scope.metrics[i].chartData = [{key: metric.id, values: data}];
            });
        });
    });

    $scope.setActiveOS = function(os) {
        var url = "/timeline/" + os + "/" + $scope.activeComponent + "/" + $scope.activeCategory + "/" + $scope.activeSubCategory;
        $location.url(url);
    };

    $scope.setTestType = function(testType) {
        $scope.testType = testType;
    };

    $scope.setActiveComponent = function(component) {
        var category = $scope.components[component].categories[0];
        if (category.subCategories.length > 0) {
            $scope.activeSubCategory = category.subCategories[0];
        } else {
            $scope.activeSubCategory = "all";
        }
        var url = "/timeline/" + $scope.activeOS + "/" + component + "/" + category.id + "/" + $scope.activeSubCategory;
        $location.url(url);
    };

    $scope.setActiveCategory = function(category) {
        var categories = $scope.components[$scope.activeComponent].categories;
        for (var i in categories) {
            if (category === categories[i].id) {
                if (categories[i].subCategories.length > 0) {
                    $scope.activeSubCategory = categories[i].subCategories[0];
                } else {
                    $scope.activeSubCategory = "all";
                }
                break;
            }
        }
        var url = "/timeline/" + $scope.activeOS + "/" + $scope.activeComponent + "/" + category + "/" + $scope.activeSubCategory;
        $location.url(url);
    };

    $scope.setActiveSubCategory = function(subCategory) {
        var url = "/timeline/" + $scope.activeOS + "/" + $scope.activeComponent + "/" + $scope.activeCategory + "/" + subCategory;
        $location.url(url);
    };

    DefineFilters($scope);
}

function CloudMenuRouter($scope, $http, $routeParams, $location) {
    $scope.activeOS = $routeParams.os;
    $scope.activeComponent = $routeParams.component;
    $scope.activeCategory = $routeParams.category;
    $scope.activeSubCategory = $routeParams.subCategory;
    $scope.activeProvider = "Provisioned";
    $scope.testType = "All";

    $http.get('/static/cloud_menu.json').success(function(menu) {
        $scope.components = menu.components;
        $scope.oses = menu.oses;
        $scope.bucket_collection = menu.bucket_collection;
        $scope.providers = menu.providers;
    });

    var url = '/api/v1/metrics/' + $scope.activeComponent + "/" + $scope.activeCategory + "/" + $scope.activeSubCategory;
    $http.get(url).success(function(metrics) {
        $scope.metrics = metrics;

        $.each(metrics, function(i, metric) {
            $http.get('/api/v1/timeline/' + metric.id).success(function(data) {
                $scope.metrics[i].chartData = [{key: metric.id, values: data}];
            });
        });
    });

    $scope.setActiveOS = function(os) {
        var url = "/cloudtimeline/" + os + "/" + $scope.activeComponent + "/" + $scope.activeCategory + "/" + $scope.activeSubCategory;
        $location.url(url);
    };

    $scope.setTestType = function(testType) {
        $scope.testType = testType;
    };

    $scope.setActiveProvider = function(provider) {
        $scope.activeProvider = provider;
    };
    
    $scope.setActiveComponent = function(component) {
        var category = $scope.components[component].categories[0];
        if (category.subCategories.length > 0) {
            $scope.activeSubCategory = category.subCategories[0];
        } else {
            $scope.activeSubCategory = "all";
        }
        var url = "/cloudTimeline/" + $scope.activeOS + "/" + component + "/" + category.id + "/" + $scope.activeSubCategory;
        $location.url(url);
    };

    $scope.setActiveCategory = function(category) {
        var categories = $scope.components[$scope.activeComponent].categories;
        for (var i in categories) {
            if (category === categories[i].id) {
                if (categories[i].subCategories.length > 0) {
                    $scope.activeSubCategory = categories[i].subCategories[0];
                } else {
                    $scope.activeSubCategory = "all";
                }
                break;
            }
        }
        var url = "/cloudTimeline/" + $scope.activeOS + "/" + $scope.activeComponent + "/" + category + "/" + $scope.activeSubCategory;
        $location.url(url);
    };

    $scope.setActiveSubCategory = function(subCategory) {
        var url = "/cloudTimeline/" + $scope.activeOS + "/" + $scope.activeComponent + "/" + $scope.activeCategory + "/" + subCategory;
        $location.url(url);
    };
    DefineFilters($scope);
}

function DefineMenu($scope, $http) {
    $http.get('/static/menu.json').success(function(menu) {
        $scope.components = menu.components;
        $scope.oses = menu.oses;
        $scope.bucket_collection = menu.bucket_collection;
    });
}

function DefineFilters($scope) {
    $scope.byOS = function(metric) {
        switch($scope.activeOS) {
            case "Linux":
                return metric.cluster.os.indexOf("Windows") === -1;
            case "Windows":
                return metric.cluster.os.indexOf("Windows") === 0;
        }
    };
    $scope.byTestType = function(metric) {
        switch($scope.testType) {
            case "Collection":
                return (metric.title.indexOf("c=") != -1 || metric.title.indexOf("collection") != -1 || (metric.title.indexOf("Collection") != -1));
            case "Bucket":
                return metric.title.indexOf("c=") === -1 && metric.title.indexOf("collection") === -1 && metric.title.indexOf("Collection") === -1;
            case "All":
                return true;
        }
    };
    $scope.byProvider = function(metric) {
        switch($scope.activeProvider) {
            case "Simulated":
                return metric.cluster.name.toLowerCase().indexOf("capella") === -1;
            case "Provisioned":
                return metric.cluster.name.toLowerCase().indexOf("capella") != -1;
            case "Serverless":
                return false;
        }
    };
}

function RunList($scope, $routeParams, $http) {
    var url = '/api/v1/runs/' + $routeParams.metric + '/' + $routeParams.build;
    $http.get(url).success(function(data) {
        $scope.runs = data;
    });
}
