/* global phantom */

var webserver = require('webserver'),
    system = require('system'),
    extract = require('./extract'),
    server = webserver.create()
;

var port = system.args[1] || 5000;

function write(error, data) {
  if (error) {
    system.stderr.writeLine(error);
  } else {
    system.stdout.writeLine(JSON.stringify(data));
  }
}

server.listen(port, function(req, res) {
  extract(req.post, write);
  res.statusCode = 200;
  res.write('');
  res.close();
});
