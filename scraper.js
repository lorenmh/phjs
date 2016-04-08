/* global phantom */

var webserver = require('webserver'),
    system = require('system'),
    extract = require('./extract'),
    server = webserver.create()
;

var port = system.args[1] || 5000;

  console.log('yo');
  console.log('ho');
  console.log('yo');
  console.log('ho');
  console.log('and');
  console.log('a');
  console.log('bucket of rum');

function write(error, data) {
  if (error) {
    system.stderr.writeLine(error);
  } else {
    system.stdout.writeLine(JSON.stringify(data));
  }
}

server.listen(port, function(req, res) {
  console.log('done did listen');
  extract(req.post, write);
  res.statusCode = 200;
  res.close();
});
