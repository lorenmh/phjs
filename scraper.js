/* global phantom */

var webserver = require('webserver'),
    system = require('system'),
    extract = require('./extract'),
    server = webserver.create()
;

function write(error, data) {
  if (error) {
    system.stderr.writeLine(error);
  } else {
    system.stdout.writeLine(JSON.stringify(data));
  }
}

server.listen(5000, function(req, res) {
  extract(req.post, write);
  
  res.statusCode = 200;
  res.close();
});
