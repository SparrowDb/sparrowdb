var isEmpty = function(str) {
    return (typeof(str) === 'undefined' || str.length === 0 || !str.trim());
};

var Storage = {
    set: function(key, value) {
        localStorage.setItem(key, value);
    },
    get: function(key) {
        return localStorage.getItem(key);
    },
    remove: function(key) {
        localStorage.removeItem(key);
    }
}

var app = angular.module("sparrowUI", ["ngRoute", "toaster", "ngAnimate"]);
app.config(function($routeProvider, $locationProvider) {
    $routeProvider
        .when("/", {
            templateUrl: "/main.html" ,
            controller: 'mainController'
        })
        .when("/db", {
            templateUrl: "/database.html",
            controller: 'dbController'
        })
        .when("/db/upload", {
            templateUrl: "database.html",
            controller: 'dbController'
        })
        .when("/script", {
            templateUrl: "/script.html",
            controller: 'scriptController'
        })
        .when("/logout", {
            controller: 'logoutController',
            templateUrl: "main.html" ,
        });
});

app.factory('sparrow', function($location) {
    var self = {};
    self.currentDb = null;
    self.currentUser = null;
    self.client = null;

    self.createClient = function(info) {
        self.currentUser = info.username;
        self.client = new SparrowDb({
            host: info.host,
            token: info.token
        });
    };

    self.getClient = function() {
        if (self.client == null) {
            window.location.href = '/login.html';
            return;
        }
        return self.client;
    };

    self.checkError = function(xhr, cb) {
        var message = 'Could not retrieve information';
        if (xhr.status == 401) {
            message = 'Not authorized';
        } else if (xhr.status == 0) {
            message = 'Connection lost';
        }

        if (cb !== undefined) {
            cb(message);
        } else {
            bootbox.alert(message, function() {
                Storage.remove('sparrow-lgn');
                window.location.href = '/login.html';
            });
        }
    };

    return self;
});


app.controller('logoutController', function($scope, $location, sparrow, $rootScope) {
    Storage.remove('sparrow-lgn');
    window.location.href = '/login.html';
});

app.controller('mainController', function($scope, $location, sparrow, $rootScope, $timeout) {
    $scope.dbData = { name: '', params: {} };

    info = Storage.get('sparrow-lgn');
    if (info == null) {
        window.location.href = '/login.html';
    }
    info = JSON.parse(info);
    sparrow.createClient(info);
    $('#username').html(info.username);

    function updateDbTable() {
        sparrow.getClient().showDatabases().success(function(r) {
            $rootScope.$apply(function() {
                $scope.databases = r.content._all;
            });
        }).error(function(xhr) {
            sparrow.checkError(xhr);
        });;
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
        if (isEmpty($scope.dbData.name)) {
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



app.controller('dbController', function($scope, $location, sparrow, $rootScope, toaster) {
    $scope.currentDb = sparrow.currentDb;
    $scope.uploadData = {};
    $scope.searchData = { key: '' };

    var updateInfo = function(cb) {
        sparrow.getClient().infoDatabase(sparrow.currentDb)
            .success(function(info) {
                sparrow.getClient().scriptList()
                    .success(function(infoScript) {
                        if (cb !== undefined) {
                            cb();
                        }
                        $scope.$apply(function() {
                            $scope.scripts = infoScript.content.scripts;
                            $scope.info = info.content;
                        });
                    }).error(function(xhr) {
                        sparrow.checkError(xhr, function() {
                            toaster.pop('error', 'Script', 'Could not get script list');
                        });
                    });
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    toaster.pop('error', 'Information', 'Could not refresh database information');
                });
            });
    }
    updateInfo();

    $scope.refresh = function() {
        updateInfo(function() {
            toaster.pop('info', 'Success', 'Information refreshed');
        });
    }

    $scope.uploadImage = function() {
        var fileUpload = angular.element(document.querySelector('#frmFile'));

        var options = {};
        options.script = $scope.uploadData.script || '';
        options.upsert = $scope.uploadData.upsert || false;

        if (isEmpty($scope.uploadData.key)) {
            bootbox.alert('Invalid image name');
            return;
        }

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
                    console.log(r)
                    $scope.keys = r.content.keys;
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
        if (isEmpty($scope.searchData.key)) {
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

var applogin = angular.module("sparrowLogin", []);
applogin.controller('loginController', function($scope, $location, $rootScope) {
    $scope.loginData = { host: '127.0.0.1:8081', username: 'sparrow', password: 'sparrow' };
    $scope.error = '';

    if (Storage.get('sparrow-lgn') != null) {
        window.location.href = '/';
    }

    $scope.doLogin = function() {
        var sparrow = new SparrowDb({ host: $scope.loginData.host });
        sparrow.login($scope.loginData.username, $scope.loginData.password)
            .success(function(r) {
                r.host = $scope.loginData.host;
                r.username = $scope.loginData.username;
                Storage.set('sparrow-lgn', JSON.stringify(r));
                window.location.href = '/';
            }).error(function(xhr) {
                $scope.$apply(function() {
                    bootbox.alert((xhr.status == 401) ? 'Invalid user and/or password' : 'Connection error');
                });
            });
    };
});

app.controller('scriptController', function($scope, $location, sparrow, $rootScope, toaster) {
    $scope.scripts = [];
    $scope.currentScript = '';
    $scope.isNewScript = true;

    var editor = ace.edit("editor");
    editor.$blockScrolling = Infinity;
    editor.getSession().setMode("ace/mode/lua");

    function updateInfo() {
        sparrow.getClient().scriptList()
            .success(function(r) {
                $scope.$apply(function() {
                    $scope.scripts = r.content.scripts;
                });
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    toaster.pop('error', 'Script', 'Could not get script list');
                });
            });
    }
    updateInfo();

    $scope.getScript = function(name) {
        $scope.isNewScript = false;
        $scope.currentScript = name;
        angular.element('#modalScriptEditor').modal('show');
        sparrow.getClient().scriptInfo(name)
            .success(function(r) {
                $scope.$apply(function() {
                    editor.getSession().setValue(r.content.script);
                });
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    bootbox.alert('Could not get script info.\n' + xhr.responseJSON.error.join("\n"));
                });
            });
    }

    $scope.newScript = function() {
        $scope.currentScript = '';
        $scope.isNewScript = true;
        angular.element('#modalScriptEditor').modal('show');
        editor.getSession().setValue('');
    }

    $scope.saveScript = function() {
        if (isEmpty($scope.currentScript)) {
            bootbox.alert('Insert a valid script name')
            return;
        }

        sparrow.getClient().saveScript($scope.currentScript, editor.getValue())
            .success(function(r) {
                updateInfo();
                toaster.pop('info', 'Script', 'Script saved !');
            }).error(function(xhr) {
                sparrow.checkError(xhr, function() {
                    toaster.pop('error', 'Script', 'Could not saved script');
                });
            });
    }

    $scope.deleteScript = function() {
        bootbox.confirm('Delete ' + $scope.currentScript + ' ?', function(r) {
            if (r == false) return;
            sparrow.getClient().deleteScript($scope.currentScript)
                .success(function(r) {
                    angular.element('#modalScriptEditor').modal('hide');
                    toaster.pop('info', 'Script', 'Script ' + $scope.currentScript + ' removed');
                    editor.getSession().setValue('');
                    updateInfo();
                }).error(function(xhr) {
                    sparrow.checkError(xhr, function() {
                        toaster.pop('error', 'Script', 'Could not delete script ' + $scope.currentScript);
                    });
                });
        });
    }
});