var scripts = require('./manifest');

module.exports = {
  match: function(path) {
    path = String(path);

    for (var i = 0; i < scripts.length; i++) {
      if (path.match(scripts[i].re)) {
        return scripts[i].fn;
      }
    }
  }
};
