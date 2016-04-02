module.exports = [
  {
    re: /^http:\/\/www\.lorenhoward\.com/,
    fn: require('./extractors/google')
  },

  {
    re: /^https?:\/\/www\.audiosf\.com/,
    fn: require('./extractors/audiosf')
  }
];
