(function() {
  var PILL_MARGIN = 2;

  function yoinkElements(parent, selector) {
    var i = 0,
        widthMap = [],
        items = parent.querySelectorAll(selector),
        len = items.length;
    for (i; i < len; i++) {
      widthMap[i] = {
        width: items[i].offsetWidth,
        node: parent.removeChild(items[i])
      };
    }
    // Sort by width
    widthMap.sort(function(first, second) {
      return second.width - first.width;
    });
    return widthMap;
  }

  function getLine(elements, maxWidth) {
    // The last element in elements is the smallest one
    var i = Math.floor((Math.random() * elements.length)),
        line = elements.splice(i,1),
        pill,
        currentWidth;
    while ((currentWidth = lineWidth(line, PILL_MARGIN)) < maxWidth) {
      pill = bestPill(elements, maxWidth - currentWidth)
      if (pill === null) { break; }
      line.push(pill);
    }
    return line;
  }

  function lineWidth(line, padding) {
    var sum = 0,
        padding = padding || 0;
    for (elem in line) {
      sum += line[elem].width;
    }
    sum += (line.length - 1) * padding;
    return sum;
  }

  // Find the (single) pill that will fit best in the allotted width.
  // Elements is already sorted in descending width.
  function bestPill(elements, maxWidth) {
    var i = 0,
        len = elements.length;
    for (i; i < len; i++) {
      if (elements[i].width <= maxWidth) {
        return elements.splice(i, 1)[0];
      }
    }
    return null;
  }

  function replace(parent, lines) {
    var i = 0, j = 0, line, wrapper;
    for (i; i < lines.length; i++) {
      line = lines[i];
      wrapper = document.createElement('div');
      for (j=0; j < line.length; j++) {
        wrapper.appendChild(line[j].node);
        parent.appendChild(wrapper);
      }
    }
  }

  function pack(parent, selector) {
    var items = yoinkElements(parent, selector),
        lines = [],
        maxWidth = parent.offsetWidth;
    while (items.length > 0) {
      lines.push(getLine(items, maxWidth));
    }
    replace(parent, lines);
  }
  window.pack = pack;
})();
