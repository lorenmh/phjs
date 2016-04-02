var webpage = require('webpage')
;

var JQUERY = 'https://code.jquery.com/jquery-1.12.0.min.js'
;

module.exports = function(path, next) {
  var page = webpage.create();

  page.open(path, function(status) {
    if (status !== 'success') {
      next('could not connect to ' + path);
      page.close();
      return;
    }

    var title = page.evaluate(function() {
      return document.title;
    });

    next(null, title);

    page.close();
  });
};
