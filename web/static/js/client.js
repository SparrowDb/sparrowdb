var Client = (function() {
    var self = {};
    var client = null;

    self.setHost = function(host) {
        self.host = host;
    };

    self.login = function(user, password) {
        client = new SparrowDb({
            host: self.host
        });
        client.login(user, password).success(function(r) {
            self.changePage('database');
        }).error(function(xhr) {
            if (xhr.status == 404) {
                self.changePage('database');
                return;
            }
            alert('Invalid user and/or password')
        });
    };

    self.getClient = function() {
        return client;
    };

    self.check = function() {
        client.ping();
    };

    self.changePage = function(addr) {
        $('#_menu').css('display', (addr == 'login') ? 'none' : 'block')
        $.ajax({
            url: '/static/' + addr + '.html',
            type: 'GET',
            dataType: 'text',
            success: function(r) {
                $('#_content').html(r);
            },
            error: function() {
                alert('Could not change page');
            }
        });
    };

    self.showInfoModal = function(dbname) {
        client.infoDatabase(dbname).success(function(r) {
            $('#modalInfoDb').modal('show')
            $('#modalInfoDbLabel').html(dbname);

            $('#modalInfoDb').on('shown.bs.modal', function() {
                $('#tblInfoCfg').find('tbody').empty();
                Object.keys(r.content.config).forEach(function(key) {
                    $('#tblInfoCfg').find('tbody')
                        .append('<tr><td>' + key + '</td><td>' + r.content.config[key] + '</td></tr>');
                });
                $('#tblInfoStatistic').find('tbody').empty();
                Object.keys(r.content.statistics).forEach(function(key) {
                    $('#tblInfoStatistic').find('tbody')
                        .append('<tr><td>' + key + '</td><td>' + r.content.statistics[key] + '</td></tr>');
                });
            });

            console.log(r)
        });
    }

    return self;
})();