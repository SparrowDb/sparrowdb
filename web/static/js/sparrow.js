 function SparrowDb(config) {
     this.config = config || {};

     this.doRequest = function(url, method, data) {
         var options = {
             url: 'http://' + this.config.host + '/' + url,
             type: method
         }

         if (data && data.constructor.name === 'FormData') {
             options.data = data;
             options.processData = false;
             options.contentType = false;
         } else {
             options.data = JSON.stringify(data);
             options.dataType = 'json';
         }

         if (this.config.token != null) {
             var self = this;
             options.beforeSend = function(xhr, settings) {
                 xhr.setRequestHeader('Authorization', 'Bearer ' + self.config.token);
             }
         }

         var resp = $.ajax(options);
         return {
             success: resp.done,
             error: resp.fail
         }
     };

     this.login = function(user, password) {
         var resp = this.doRequest('user/login', 'POST', { username: user, password: password })
         var self = this;
         resp.success(function(r) {
             self.token = r.token;
         });
         return resp
     };

     this.ping = function() {
         return this.doRequest('ping', 'GET');
     };

     this.createDatabase = function(dbname, options) {
         return this.doRequest('api/' + dbname, 'PUT', options);
     };

     this.dropDatabase = function(dbname) {
         return this.doRequest('api/' + dbname, 'DELETE');
     };

     this.showDatabases = function() {
         return this.infoDatabase('_all');
     };

     this.infoDatabase = function(dbname) {
         return this.doRequest('api/' + dbname, 'GET');
     };

     this.imageInfo = function(dbname, key) {
         return this.doRequest('api/' + dbname + '/' + key, 'GET');
     };

     this.uploadImage = function(dbname, key, inputId, options) {
         options = options || {};
         var data = new FormData();

         if (options.upsert) {
             data.append('upsert', true);
         }

         if (options.script) {
             data.append('script', options.script);
         }
         data.append('uploadfile', inputId[0].files[0]);
         data.append('dbname', dbname);
         data.append('key', key);
         return this.doRequest('api/' + dbname + '/' + key, 'PUT', data);
     };

     this.deleteImage = function(dbname, key) {
         return this.doRequest('api/' + dbname + '/' + key, 'DELETE');
     };

     this.keys = function(dbname) {
         return this.doRequest('api/' + dbname + '/_keys', 'GET');
     };

     this.scriptList = function() {
         return this.doRequest('script/_all', 'GET');
     };

     this.scriptInfo = function(name) {
         return this.doRequest('script/' + name, 'GET');
     };

     this.saveScript = function(name, script) {
         return this.doRequest('script/' + name, 'POST', { content: script });
     };

     this.deleteScript = function(name) {
         return this.doRequest('script/' + name, 'DELETE');
     };
 };