// transform /path/to/file.ext to file.ext
export function basename(path = '') {
  return path ? path.replace(/^.*\//, '') : '';
}

// transform /path/to/file.ext to .ext
export function extname(path = '') {
  const [_, ext] = /.+(\.[^.]+)$/.exec(path) || [];
  return ext || '';
}

// test whether a variable is an object
export function isObject(obj) {
  return Object.prototype.toString.call(obj) === '[object Object]';
}

// returns whether a dark theme is enabled
export function isDarkTheme() {
  return document.documentElement.classList.contains('theme-arc-green');
}

// removes duplicate elements in an array
export function uniq(arr) {
  return Array.from(new Set(arr));
}

// strip <tags> from a string
export function stripTags(text) {
  return text.replace(/<[^>]*>?/gm, '');
}

// searches the inclusive range [minValue, maxValue].
// credits: https://matthiasott.com/notes/write-your-media-queries-in-pixels-not-ems
export function mqBinarySearch(feature, minValue, maxValue, step, unit) {
  if (maxValue - minValue < step) {
    return minValue;
  }
  const mid = Math.ceil((minValue + maxValue) / 2 / step) * step;
  if (matchMedia(`screen and (min-${feature}:${mid}${unit})`).matches) {
    return mqBinarySearch(feature, mid, maxValue, step, unit); // feature is >= mid
  }
  return mqBinarySearch(feature, minValue, mid - step, step, unit); // feature is < mid
}
