var app = angular.module("sparrowUI", ["ngRoute"]);
app.config(function($routeProvider) {  
    $routeProvider
        .when("/", {
            templateUrl: "main.html" ,
            controller: 'mainController'
        })
        .when("/login", {
            templateUrl: "login.html" ,
            controller: 'loginController'
        })
        .when("/db", {
            templateUrl: "database.html",
            controller: 'dbController'
        })
        .when("/db/upload", {
            templateUrl: "database.html",
            controller: 'dbController'
        });
});

app.factory('sparrow', function($location) {
    var self = {};
    self.token = '';
    self.currentDb = null;
    self.currentUser = 'sparrow';

    self.connect = function(_host) {
        self.client = new SparrowDb({
            host: _host
        });
    };

    self.getClient = function() {
        if (self.client == null) {
            $location.path("/login");
            return;
        }
        return self.client;
    };

    self.checkError = function(xhr, cb) {
        if (xhr.status == 0) {
            bootbox.alert("Lost connection with server");
        } else {
            cb(xhr);
        }
    };

    return self;
});

app.controller('loginController', function($scope, $location, sparrow, $rootScope) {
    $scope.loginData = { host: '127.0.0.1:8081', username: 'sparrow', password: 'sparrow' };
    $scope.error = '';

    $scope.doLogin = function() {
        sparrow.connect($scope.loginData.host);

        sparrow.getClient().login($scope.loginData.username, $scope.loginData.password)
            .success(function(r) {
                sparrow.token = r.token;
                sparrow.currentUser = $scope.loginData.username;
                $rootScope.$apply(function() {
                    $location.path("/");
                });
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    $scope.$apply(function() {
                        $scope.error = 'Invalid user and/or password';
                    });
                });
            });
    };
});

app.controller('mainController', function($scope, $location, sparrow, $rootScope) {
    $scope.dbData = { name: '', params: {} };

    function updateDbTable() {
        sparrow.getClient().showDatabases().success(function(r) {
            $scope.$apply(function() {
                $scope.databases = r.content._all;
            });
        });
    }
    updateDbTable();

    $scope.dbInfo = function(db) {
        sparrow.currentDb = db;
        $location.path("/db");
    };

    $scope.dbDrop = function(db) {
        bootbox.confirm('Drop ' + db + ' ?', function(r) {
            if (r == false) return;
            sparrow.getClient().dropDatabase(db)
                .success(function(r) {
                    bootbox.alert('Database dropped')
                    updateDbTable();
                }).error(function(xhr) {
                    sparrow.checkError(xhr, function() {
                        if (xhr.status == 404) {
                            $rootScope.$apply(function() {
                                $location.path("/");
                            });
                            return;
                        }
                        bootbox.alert('Could not drop database');
                    });
                });
        });
    };

    $scope.createDb = function() {
        if ($scope.dbData.name == '') {
            bootbox.alert('Insert a valid database name')
            return;
        };

        sparrow.getClient().createDatabase($scope.dbData.name, $scope.dbData.params)
            .success(function(r) {
                bootbox.alert('Database created')
                angular.element('#modalCreateDb').modal('hide');
                updateDbTable();
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    if (xhr.status == 404) {
                        $rootScope.$apply(function() {
                            $location.path("/");
                        });
                        return;
                    }
                    bootbox.alert('Could not create database');
                });
            });
    };

    $scope.addParam = function() {
        var p = $scope.input.dbparam.toLowerCase();
        var v = $scope.input.dbvalue.toLowerCase();
        if (p == '' || v == '') return;

        switch (p) {
            case 'read_only':
            case 'generate_token':
                v = (p == 'true') ? true : false;
                break;
            case 'max_cache_size':
            case 'bloomfilter_fpp':
                v = parseFloat(v);
                break;
        }

        $scope.dbData.params[p] = v;
    };

    $scope.removeParam = function(key) {
        delete $scope.dbData.params[key];
    };
});


app.controller('dbController', function($scope, $location, sparrow, $rootScope) {
    $scope.currentDb = sparrow.currentDb;
    $scope.uploadData = {};
    $scope.searchData = { key: '' };

    var updateInfo = function() {
        sparrow.getClient().infoDatabase(sparrow.currentDb)
            .success(function(r) {
                $scope.$apply(function() {
                    $scope.info = r.content;
                });
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    $location.path("/");
                });
            });
    }
    updateInfo();

    $scope.refresh = function() {
        updateInfo();
    }

    $scope.uploadImage = function() {
        var fileUpload = angular.element(document.querySelector('#frmFile'));

        var options = {};
        options.script = $scope.uploadData.script || '';
        options.upsert = $scope.uploadData.upsert || false;

        sparrow.getClient().uploadImage(
                sparrow.currentDb,
                $scope.uploadData.key,
                fileUpload,
                options
            )
            .success(function(r) {
                bootbox.alert('Image ' + $scope.uploadData.key + ' sent to ' + sparrow.currentDb);
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    bootbox.alert('Could not send image.\n' + xhr.responseJSON.error.join("\n"));
                });
            });
    }

    $scope.showKeys = function() {
        sparrow.getClient().keys(sparrow.currentDb)
            .success(function(r) {
                $scope.$apply(function() {
                    $scope.keys = r;
                });
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    $scope.$apply(function() {
                        $scope.imgInfo = {};
                    });
                    bootbox.alert('Could not get image list.\n' + xhr.responseJSON.error.join("\n"));
                });
            });
    }

    $scope.imageInfo = function() {
        if ($scope.searchData.key == '') {
            bootbox.alert('Insert a key');
            return;
        }

        sparrow.getClient().imageInfo(
                sparrow.currentDb,
                $scope.searchData.key
            )
            .success(function(r) {
                $scope.$apply(function() {
                    $scope.imgInfo = r.content;
                });
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    $scope.$apply(function() {
                        $scope.imgInfo = {};
                    });
                    bootbox.alert('Could not get image info.\n' + xhr.responseJSON.error.join("\n"));
                });
            });
    }
});