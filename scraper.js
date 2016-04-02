/* jshint ignore: start */

var system = require('system'),
    s = require('./scrapes/google.js')
;

var path
;

// while (path = system.stdin.readLine()) {
//   extract(path, function(data) {
//     system.stdout.writeLine(data);
//   });
// }

s('http://lorenhoward.com', function(err, data) {
  if (err) {
    system.stderr.writeLine(err);
  } else {
    system.stdout.writeLine('title: ' + data);
  }
});

s('http://www.facebook.com', function(err, data) {
  if (err) {
    system.stderr.writeLine(err);
  } else {
    system.stdout.writeLine('title: ' + data);
  }
});

s('http://www.hiqlabs.com', function(err, data) {
  if (err) {
    system.stderr.writeLine(err);
  } else {
    system.stdout.writeLine('title: ' + data);
  }
});
