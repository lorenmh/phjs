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

    var scenes = page.evaluate(function() {
      function toArray(query) {
        return Array.prototype.slice.call(query);
      }


      return (
        toArray(document.querySelectorAll('#cal-event-entry'))
          .map(function(el) {
            var title_el = el.querySelector('.home-event-blocks-title');
            var link = title_el.querySelector('a');
            var url = link.href;
            return {
              title: title_el.textContent,
              url: url,
              date: el.querySelector('.home-event-blocks-date').textContent
            };
          })
      );
    });

    scenes.forEach(function(d) { next(null, d); });

    page.close();
  });
};
