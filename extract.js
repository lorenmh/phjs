var extractors = require('./manifest');

module.exports = function(path, next) {
  path = String(path);

  for (var i = 0; i < extractors.length; i++) {

    if (path.match(extractors[i].re)) {
      return extractors[i].fn(path, next);
    }

  }

  next('couldn\'t find extractor for path ' + path);
};
