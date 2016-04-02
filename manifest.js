module.exports = [
  {
    re: /www\..+\.com/,
    fn: require('./scrapes/google')
  }
];
